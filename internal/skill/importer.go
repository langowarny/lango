package skill

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// resourceDirs lists the conventional subdirectories that may contain resource files.
var resourceDirs = []string{"scripts", "references", "assets"}

// ImportConfig holds configuration for bulk import operations.
type ImportConfig struct {
	MaxSkills   int
	Concurrency int
	Timeout     time.Duration
}

// Importer fetches SKILL.md files from GitHub repositories or arbitrary URLs.
type Importer struct {
	client *http.Client
	logger *zap.SugaredLogger
}

// GitHubRef represents a parsed GitHub repository reference.
type GitHubRef struct {
	Owner  string
	Repo   string
	Branch string
	Path   string
}

// ImportResult summarises the outcome of a bulk import operation.
type ImportResult struct {
	Imported []string
	Skipped  []string
	Errors   []string
}

// NewImporter creates a new skill importer.
func NewImporter(logger *zap.SugaredLogger) *Importer {
	return &Importer{
		client: &http.Client{Timeout: 30 * time.Second},
		logger: logger,
	}
}

// NewImporterWithClient creates an importer with a custom HTTP client (for testing).
func NewImporterWithClient(client *http.Client, logger *zap.SugaredLogger) *Importer {
	return &Importer{
		client: client,
		logger: logger,
	}
}

// IsGitHubURL returns true if the URL points to github.com.
func IsGitHubURL(rawURL string) bool {
	return strings.Contains(rawURL, "github.com")
}

// ParseGitHubURL parses a GitHub URL into owner, repo, branch, and path components.
// Supported formats:
//   - https://github.com/owner/repo
//   - https://github.com/owner/repo/tree/branch
//   - https://github.com/owner/repo/tree/branch/path/to/dir
func ParseGitHubURL(rawURL string) (*GitHubRef, error) {
	rawURL = strings.TrimSuffix(rawURL, "/")

	// Remove protocol prefix.
	u := rawURL
	for _, prefix := range []string{"https://", "http://"} {
		u = strings.TrimPrefix(u, prefix)
	}

	// Remove github.com prefix.
	u = strings.TrimPrefix(u, "github.com/")

	parts := strings.SplitN(u, "/", 4) // owner/repo[/tree/branch/path]
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid GitHub URL: need at least owner/repo")
	}

	ref := &GitHubRef{
		Owner:  parts[0],
		Repo:   parts[1],
		Branch: "main",
	}

	if len(parts) >= 4 {
		// parts[2] is "tree" or "blob"
		rest := parts[3] // branch/path or just branch
		slashIdx := strings.Index(rest, "/")
		if slashIdx < 0 {
			ref.Branch = rest
		} else {
			ref.Branch = rest[:slashIdx]
			ref.Path = rest[slashIdx+1:]
		}
	}

	return ref, nil
}

// gitHubContentsEntry is a single entry from the GitHub Contents API response.
type gitHubContentsEntry struct {
	Name string `json:"name"`
	Type string `json:"type"` // "file" or "dir"
	Path string `json:"path"`
}

// DiscoverSkills lists subdirectories at the given path in a GitHub repo.
// Each subdirectory is assumed to contain a SKILL.md file.
func (im *Importer) DiscoverSkills(ctx context.Context, ref *GitHubRef) ([]string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
		ref.Owner, ref.Repo, ref.Path, ref.Branch)

	body, err := im.doGet(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("discover skills: %w", err)
	}

	var entries []gitHubContentsEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("parse GitHub contents: %w", err)
	}

	var skills []string
	for _, e := range entries {
		if e.Type == "dir" {
			skills = append(skills, e.Name)
		}
	}
	return skills, nil
}

// gitHubFileResponse is the response from the GitHub Contents API for a single file.
type gitHubFileResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// FetchSkillMD fetches a SKILL.md file from a GitHub repo at {path}/{skillName}/SKILL.md.
func (im *Importer) FetchSkillMD(ctx context.Context, ref *GitHubRef, skillName string) ([]byte, error) {
	filePath := skillName + "/SKILL.md"
	if ref.Path != "" {
		filePath = ref.Path + "/" + filePath
	}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
		ref.Owner, ref.Repo, filePath, ref.Branch)

	body, err := im.doGet(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("fetch SKILL.md for %q: %w", skillName, err)
	}

	var file gitHubFileResponse
	if err := json.Unmarshal(body, &file); err != nil {
		return nil, fmt.Errorf("parse GitHub file response: %w", err)
	}

	if file.Encoding != "base64" {
		return nil, fmt.Errorf("unexpected encoding %q (expected base64)", file.Encoding)
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(file.Content, "\n", ""))
	if err != nil {
		return nil, fmt.Errorf("decode base64 content: %w", err)
	}

	return decoded, nil
}

// FetchFromURL fetches raw content from an arbitrary URL via HTTP GET.
func (im *Importer) FetchFromURL(ctx context.Context, rawURL string) ([]byte, error) {
	return im.doGet(ctx, rawURL)
}

// hasGit returns true if the git binary is available on PATH.
func hasGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// cloneRepo clones a GitHub repo to a temp directory and returns the path.
// Shallow clone (depth=1) for speed.
func (im *Importer) cloneRepo(ctx context.Context, ref *GitHubRef) (string, error) {
	tmpDir, err := os.MkdirTemp("", "lango-skill-import-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", ref.Owner, ref.Repo)
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", "--branch", ref.Branch, repoURL, tmpDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("git clone: %s: %w", string(out), err)
	}

	return tmpDir, nil
}

// copyResourceDirs copies conventional resource directories from a cloned skill dir to the store.
func copyResourceDirs(ctx context.Context, srcDir, skillName string, store SkillStore) {
	for _, dir := range resourceDirs {
		resDir := filepath.Join(srcDir, dir)
		entries, err := os.ReadDir(resDir)
		if err != nil {
			continue // directory doesn't exist — skip
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			data, err := os.ReadFile(filepath.Join(resDir, e.Name()))
			if err != nil {
				continue
			}
			_ = store.SaveResource(ctx, skillName, filepath.Join(dir, e.Name()), data)
		}
	}
}

// importViaGit clones the repo and imports skills from the local filesystem.
func (im *Importer) importViaGit(ctx context.Context, ref *GitHubRef, store SkillStore, cfg ImportConfig) (*ImportResult, error) {
	cloneDir, err := im.cloneRepo(ctx, ref)
	if err != nil {
		im.logger.Warnw("git clone failed, falling back to HTTP", "error", err)
		return im.importViaHTTP(ctx, ref, store, cfg)
	}
	defer os.RemoveAll(cloneDir)

	baseDir := cloneDir
	if ref.Path != "" {
		baseDir = filepath.Join(cloneDir, ref.Path)
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("read cloned dir: %w", err)
	}

	sourceURL := fmt.Sprintf("https://github.com/%s/%s", ref.Owner, ref.Repo)
	result := &ImportResult{}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		skillDir := filepath.Join(baseDir, e.Name())
		mdPath := filepath.Join(skillDir, "SKILL.md")
		raw, err := os.ReadFile(mdPath)
		if err != nil {
			continue // not a skill directory
		}

		if cfg.MaxSkills > 0 && len(result.Imported)+len(result.Skipped) >= cfg.MaxSkills {
			break
		}

		entry, parseErr := ParseSkillMD(raw)
		if parseErr != nil {
			im.logger.Warnw("skip skill: parse failed", "skill", e.Name(), "error", parseErr)
			result.Errors = append(result.Errors, fmt.Sprintf("%s: parse: %v", e.Name(), parseErr))
			continue
		}

		entry.Source = sourceURL

		if existing, _ := store.Get(ctx, entry.Name); existing != nil {
			result.Skipped = append(result.Skipped, entry.Name)
			continue
		}

		if saveErr := store.Save(ctx, *entry); saveErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: save: %v", e.Name(), saveErr))
			continue
		}

		// Copy resource directories (scripts/, references/, assets/).
		copyResourceDirs(ctx, skillDir, entry.Name, store)

		result.Imported = append(result.Imported, entry.Name)
	}

	return result, nil
}

// importViaHTTP imports skills using the GitHub Contents API (original HTTP approach).
func (im *Importer) importViaHTTP(ctx context.Context, ref *GitHubRef, store SkillStore, cfg ImportConfig) (*ImportResult, error) {
	skillNames, err := im.DiscoverSkills(ctx, ref)
	if err != nil {
		return nil, err
	}

	if cfg.MaxSkills > 0 && len(skillNames) > cfg.MaxSkills {
		im.logger.Warnw("skill count exceeds max, truncating",
			"discovered", len(skillNames), "max", cfg.MaxSkills)
		skillNames = skillNames[:cfg.MaxSkills]
	}

	sourceURL := fmt.Sprintf("https://github.com/%s/%s", ref.Owner, ref.Repo)
	result := &ImportResult{}
	var mu sync.Mutex

	concurrency := cfg.Concurrency
	if concurrency <= 0 {
		concurrency = 5
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, name := range skillNames {
		if ctx.Err() != nil {
			mu.Lock()
			result.Errors = append(result.Errors, fmt.Sprintf("cancelled: %v", ctx.Err()))
			mu.Unlock()
			break
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(skillName string) {
			defer wg.Done()
			defer func() { <-sem }()

			raw, fetchErr := im.FetchSkillMD(ctx, ref, skillName)
			if fetchErr != nil {
				im.logger.Warnw("skip skill: fetch failed", "skill", skillName, "error", fetchErr)
				mu.Lock()
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", skillName, fetchErr))
				mu.Unlock()
				return
			}

			entry, parseErr := ParseSkillMD(raw)
			if parseErr != nil {
				im.logger.Warnw("skip skill: parse failed", "skill", skillName, "error", parseErr)
				mu.Lock()
				result.Errors = append(result.Errors, fmt.Sprintf("%s: parse: %v", skillName, parseErr))
				mu.Unlock()
				return
			}

			entry.Source = sourceURL

			if existing, _ := store.Get(ctx, entry.Name); existing != nil {
				mu.Lock()
				result.Skipped = append(result.Skipped, entry.Name)
				mu.Unlock()
				return
			}

			if saveErr := store.Save(ctx, *entry); saveErr != nil {
				mu.Lock()
				result.Errors = append(result.Errors, fmt.Sprintf("%s: save: %v", skillName, saveErr))
				mu.Unlock()
				return
			}

			// Fetch and save resource files via HTTP.
			im.fetchAndSaveResources(ctx, ref, skillName, entry.Name, store)

			mu.Lock()
			result.Imported = append(result.Imported, entry.Name)
			mu.Unlock()
		}(name)
	}

	wg.Wait()
	return result, nil
}

// fetchAndSaveResources fetches resource files for a skill from GitHub via the Contents API.
func (im *Importer) fetchAndSaveResources(ctx context.Context, ref *GitHubRef, dirName, skillName string, store SkillStore) {
	for _, dir := range resourceDirs {
		resPath := dirName + "/" + dir
		if ref.Path != "" {
			resPath = ref.Path + "/" + resPath
		}

		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
			ref.Owner, ref.Repo, resPath, ref.Branch)

		body, err := im.doGet(ctx, apiURL)
		if err != nil {
			continue // directory doesn't exist — skip
		}

		var entries []gitHubContentsEntry
		if err := json.Unmarshal(body, &entries); err != nil {
			continue
		}

		for _, e := range entries {
			if e.Type != "file" {
				continue
			}
			data, err := im.fetchGitHubFileContent(ctx, ref, e.Path)
			if err != nil {
				continue
			}
			_ = store.SaveResource(ctx, skillName, filepath.Join(dir, e.Name), data)
		}
	}
}

// fetchGitHubFileContent fetches a single file's content from GitHub via the Contents API.
func (im *Importer) fetchGitHubFileContent(ctx context.Context, ref *GitHubRef, filePath string) ([]byte, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
		ref.Owner, ref.Repo, filePath, ref.Branch)

	body, err := im.doGet(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	var file gitHubFileResponse
	if err := json.Unmarshal(body, &file); err != nil {
		return nil, fmt.Errorf("parse file response: %w", err)
	}

	if file.Encoding != "base64" {
		return nil, fmt.Errorf("unexpected encoding %q", file.Encoding)
	}

	return base64.StdEncoding.DecodeString(strings.ReplaceAll(file.Content, "\n", ""))
}

// ImportFromRepo discovers and imports all skills from a GitHub repository.
// Prefers git clone when available, falls back to GitHub HTTP API.
func (im *Importer) ImportFromRepo(ctx context.Context, ref *GitHubRef, store SkillStore, cfg ImportConfig) (*ImportResult, error) {
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	if hasGit() {
		return im.importViaGit(ctx, ref, store, cfg)
	}

	return im.importViaHTTP(ctx, ref, store, cfg)
}

// ImportSingleWithResources imports a single skill from a GitHub repo, including resource files.
func (im *Importer) ImportSingleWithResources(ctx context.Context, ref *GitHubRef, skillName string, store SkillStore) (*SkillEntry, error) {
	if hasGit() {
		return im.importSingleViaGit(ctx, ref, skillName, store)
	}
	return im.importSingleViaHTTP(ctx, ref, skillName, store)
}

func (im *Importer) importSingleViaGit(ctx context.Context, ref *GitHubRef, skillName string, store SkillStore) (*SkillEntry, error) {
	cloneDir, err := im.cloneRepo(ctx, ref)
	if err != nil {
		im.logger.Warnw("git clone failed, falling back to HTTP", "error", err)
		return im.importSingleViaHTTP(ctx, ref, skillName, store)
	}
	defer os.RemoveAll(cloneDir)

	skillDir := filepath.Join(cloneDir, skillName)
	if ref.Path != "" {
		skillDir = filepath.Join(cloneDir, ref.Path, skillName)
	}

	raw, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if err != nil {
		return nil, fmt.Errorf("read SKILL.md: %w", err)
	}

	sourceURL := fmt.Sprintf("https://github.com/%s/%s", ref.Owner, ref.Repo)
	entry, err := ParseSkillMD(raw)
	if err != nil {
		return nil, fmt.Errorf("parse SKILL.md: %w", err)
	}
	entry.Source = sourceURL

	if err := store.Save(ctx, *entry); err != nil {
		return nil, fmt.Errorf("save skill %q: %w", entry.Name, err)
	}

	copyResourceDirs(ctx, skillDir, entry.Name, store)
	return entry, nil
}

func (im *Importer) importSingleViaHTTP(ctx context.Context, ref *GitHubRef, skillName string, store SkillStore) (*SkillEntry, error) {
	raw, err := im.FetchSkillMD(ctx, ref, skillName)
	if err != nil {
		return nil, fmt.Errorf("fetch SKILL.md for %q: %w", skillName, err)
	}

	sourceURL := fmt.Sprintf("https://github.com/%s/%s", ref.Owner, ref.Repo)
	entry, err := ParseSkillMD(raw)
	if err != nil {
		return nil, fmt.Errorf("parse SKILL.md: %w", err)
	}
	entry.Source = sourceURL

	if err := store.Save(ctx, *entry); err != nil {
		return nil, fmt.Errorf("save skill %q: %w", entry.Name, err)
	}

	im.fetchAndSaveResources(ctx, ref, skillName, entry.Name, store)
	return entry, nil
}

// ImportSingle imports a single skill from raw SKILL.md content.
func (im *Importer) ImportSingle(ctx context.Context, raw []byte, sourceURL string, store SkillStore) (*SkillEntry, error) {
	entry, err := ParseSkillMD(raw)
	if err != nil {
		return nil, fmt.Errorf("parse SKILL.md: %w", err)
	}
	entry.Source = sourceURL

	if err := store.Save(ctx, *entry); err != nil {
		return nil, fmt.Errorf("save skill %q: %w", entry.Name, err)
	}
	return entry, nil
}

func (im *Importer) doGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "lango-skill-importer")

	resp, err := im.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return body, nil
}
