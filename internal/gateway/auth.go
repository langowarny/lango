package gateway

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/langowarny/lango/internal/config"
	"github.com/langowarny/lango/internal/logging"
	"github.com/langowarny/lango/internal/session"
	"golang.org/x/oauth2"
)

var authLogger = logging.SubsystemSugar("gateway-auth")

// AuthManager manages authentication providers
type AuthManager struct {
	providers map[string]*OIDCProvider
	store     session.Store
}

// OIDCProvider handles OIDC authentication for a specific provider
type OIDCProvider struct {
	Name         string
	Config       config.OIDCProviderConfig
	OAuthConfig  *oauth2.Config
	OIDCProvider *oidc.Provider
	Verifier     *oidc.IDTokenVerifier
}

// NewAuthManager creates a new AuthManager
func NewAuthManager(cfg config.AuthConfig, store session.Store) (*AuthManager, error) {
	am := &AuthManager{
		providers: make(map[string]*OIDCProvider),
		store:     store,
	}

	for name, providerCfg := range cfg.Providers {
		provider, err := NewOIDCProvider(name, providerCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create provider %s: %w", name, err)
		}
		am.providers[name] = provider
	}

	return am, nil
}

// NewOIDCProvider creates a new OIDCProvider
func NewOIDCProvider(name string, cfg config.OIDCProviderConfig) (*OIDCProvider, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider %q: %w", cfg.IssuerURL, err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       cfg.Scopes,
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return &OIDCProvider{
		Name:         name,
		Config:       cfg,
		OAuthConfig:  oauthConfig,
		OIDCProvider: provider,
		Verifier:     verifier,
	}, nil
}

// RegisterRoutes registers auth routes on the router
func (am *AuthManager) RegisterRoutes(r chi.Router) {
	r.Get("/auth/login/{provider}", am.handleLogin)
	r.Get("/auth/callback/{provider}", am.handleCallback)
}

func (am *AuthManager) handleLogin(w http.ResponseWriter, r *http.Request) {
	provName := chi.URLParam(r, "provider")
	provider, ok := am.providers[provName]
	if !ok {
		http.Error(w, "provider not found", http.StatusNotFound)
		return
	}

	state, err := generateRandomString(32)
	if err != nil {
		http.Error(w, "failed to generate state", http.StatusInternalServerError)
		return
	}

	// TODO: Store state in cookie to verify callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(10 * time.Minute),
	})

	http.Redirect(w, r, provider.OAuthConfig.AuthCodeURL(state), http.StatusFound)
}

func (am *AuthManager) handleCallback(w http.ResponseWriter, r *http.Request) {
	provName := chi.URLParam(r, "provider")
	provider, ok := am.providers[provName]
	if !ok {
		http.Error(w, "provider not found", http.StatusNotFound)
		return
	}

	// Verify state
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "state cookie missing", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != cookie.Value {
		http.Error(w, "state mismatch", http.StatusBadRequest)
		return
	}

	// Exchange code
	ctx := r.Context()
	oauth2Token, err := provider.OAuthConfig.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		authLogger.Errorw("token exchange error", "provider", provName, "error", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	// Verify ID Token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token field in oauth2 token", http.StatusInternalServerError)
		return
	}

	idToken, err := provider.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		authLogger.Errorw("token verification error", "provider", provName, "error", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	// Extract claims
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Sub           string `json:"sub"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "failed to parse claims", http.StatusInternalServerError)
		return
	}

	// Create Session
	sessionKey := fmt.Sprintf("sess_%s", uuid.New().String())
	sess := &session.Session{
		Key:       sessionKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]string{
			"email":    claims.Email,
			"sub":      claims.Sub,
			"provider": provName,
		},
	}

	if err := am.store.Create(sess); err != nil {
		authLogger.Errorw("session creation error", "provider", provName, "error", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	// Write Session Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "lango_session",
		Value:    sessionKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour), // Configurable TTL?
	})

	fmt.Fprintf(w, "Login successful! Welcome %s", claims.Email)
}

func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
