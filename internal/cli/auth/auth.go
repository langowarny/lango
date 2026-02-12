package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/langowarny/lango/internal/config"
	"github.com/langowarny/lango/internal/logging"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [provider]",
		Short: "Login to an OAuth provider (google, github)",
		Args:  cobra.ExactArgs(1),
		RunE:  runLogin,
	}
	// Do not add config flag here, inherit from root command
	return cmd
}

func runLogin(cmd *cobra.Command, args []string) error {
	providerName := args[0]
	ctx := context.Background()
	log := logging.Sugar()

	// Load config
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	// If empty string, load default
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get provider config
	pCfg, ok := cfg.Providers[providerName]
	if !ok {
		// Fallback for "google" if configured as "gemini" type?
		// Or assume the key in the map IS the provider name.
		// Let's check if the user configured "google" explicitly.
		// If not, maybe check "gemini" and see if it has OAuth config?
		// For now, strict: key must match.
		return fmt.Errorf("provider '%s' not found in configuration", providerName)
	}

	if pCfg.ClientID == "" || pCfg.ClientSecret == "" {
		return fmt.Errorf("provider '%s' missing clientId or clientSecret", providerName)
	}

	var endpoint oauth2.Endpoint
	var scopes []string

	switch providerName {
	case "google":
		endpoint = google.Endpoint
		scopes = []string{"https://www.googleapis.com/auth/cloud-platform"} // Default for Vertex AI
		if len(pCfg.Scopes) > 0 {
			scopes = pCfg.Scopes
		}
	case "github":
		endpoint = github.Endpoint
		scopes = []string{"read:user", "user:email"}
		if len(pCfg.Scopes) > 0 {
			scopes = pCfg.Scopes
		}
	default:
		return fmt.Errorf("unsupported provider: %s", providerName)
	}

	// Setup local listener
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to start local listener: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURL := fmt.Sprintf("http://localhost:%d/callback", port)

	conf := &oauth2.Config{
		ClientID:     pCfg.ClientID,
		ClientSecret: pCfg.ClientSecret,
		Scopes:       scopes,
		Endpoint:     endpoint,
		RedirectURL:  redirectURL,
	}

	state, err := generateState()
	if err != nil {
		return fmt.Errorf("generate OAuth state: %w", err)
	}
	authURL := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	fmt.Printf("Opening browser for login: %s\n", authURL)
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Failed to open browser: %v\nPlease visit the URL manually.\n", err)
	}

	codeChan := make(chan string)
	errChan := make(chan error)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/callback" {
				http.NotFound(w, r)
				return
			}
			queryState := r.URL.Query().Get("state")
			if queryState != state {
				http.Error(w, "State mismatch", http.StatusBadRequest)
				errChan <- fmt.Errorf("state mismatch")
				return
			}
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "Code missing", http.StatusBadRequest)
				errChan <- fmt.Errorf("code missing")
				return
			}
			fmt.Fprintf(w, "Login successful! You can close this window.")
			codeChan <- code
		}),
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case code := <-codeChan:
		// Exchange code
		token, err := conf.Exchange(ctx, code)
		if err != nil {
			return fmt.Errorf("failed to exchange token: %w", err)
		}

		if err := saveToken(providerName, token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
		log.Infow("Login successful", "provider", providerName)
		fmt.Printf("Successfully logged in to %s\n", providerName)

	case err := <-errChan:
		return fmt.Errorf("login failed: %w", err)
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("login timed out")
	}

	return nil
}

func saveToken(provider string, token *oauth2.Token) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	tokenDir := filepath.Join(home, ".lango", "tokens")
	if err := os.MkdirAll(tokenDir, 0700); err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(tokenDir, provider+".json"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
