package skill

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

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

// ImportFromRepo discovers and imports all skills from a GitHub repository.
func (im *Importer) ImportFromRepo(ctx context.Context, ref *GitHubRef, store SkillStore) (*ImportResult, error) {
	skillNames, err := im.DiscoverSkills(ctx, ref)
	if err != nil {
		return nil, err
	}

	sourceURL := fmt.Sprintf("https://github.com/%s/%s", ref.Owner, ref.Repo)
	result := &ImportResult{}

	for _, name := range skillNames {
		raw, err := im.FetchSkillMD(ctx, ref, name)
		if err != nil {
			im.logger.Warnw("skip skill: fetch failed", "skill", name, "error", err)
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", name, err))
			continue
		}

		entry, err := ParseSkillMD(raw)
		if err != nil {
			im.logger.Warnw("skip skill: parse failed", "skill", name, "error", err)
			result.Errors = append(result.Errors, fmt.Sprintf("%s: parse: %v", name, err))
			continue
		}

		// Override source to track import origin.
		entry.Source = sourceURL

		// Check if skill already exists.
		if existing, _ := store.Get(ctx, entry.Name); existing != nil {
			result.Skipped = append(result.Skipped, entry.Name)
			continue
		}

		if err := store.Save(ctx, *entry); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: save: %v", name, err))
			continue
		}
		result.Imported = append(result.Imported, entry.Name)
	}

	return result, nil
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
