package security

import (
	"context"
	"fmt"
	"sync"
)

// ConnectionChecker provides connection status checking.
type ConnectionChecker interface {
	IsConnected() bool
}

// CompositeCryptoProvider implements CryptoProvider with fallback logic.
// It tries the primary provider first (typically companion), then falls back to local.
type CompositeCryptoProvider struct {
	mu        sync.RWMutex
	primary   CryptoProvider
	fallback  CryptoProvider
	checker   ConnectionChecker
	usedLocal bool
}

// NewCompositeCryptoProvider creates a new CompositeCryptoProvider.
func NewCompositeCryptoProvider(primary CryptoProvider, fallback CryptoProvider, checker ConnectionChecker) *CompositeCryptoProvider {
	return &CompositeCryptoProvider{
		primary:  primary,
		fallback: fallback,
		checker:  checker,
	}
}

// UsedLocal returns true if the last operation used the local fallback.
func (c *CompositeCryptoProvider) UsedLocal() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.usedLocal
}

// Sign implements CryptoProvider.
func (c *CompositeCryptoProvider) Sign(ctx context.Context, keyID string, payload []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Try primary (companion) if connected
	if c.checker != nil && c.checker.IsConnected() {
		sig, err := c.primary.Sign(ctx, keyID, payload)
		if err == nil {
			c.usedLocal = false
			return sig, nil
		}
		// Fall through to fallback on error
	}

	// Use fallback (local)
	if c.fallback == nil {
		return nil, fmt.Errorf("no crypto provider available")
	}

	c.usedLocal = true
	return c.fallback.Sign(ctx, keyID, payload)
}

// Encrypt implements CryptoProvider.
func (c *CompositeCryptoProvider) Encrypt(ctx context.Context, keyID string, plaintext []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Try primary (companion) if connected
	if c.checker != nil && c.checker.IsConnected() {
		ciphertext, err := c.primary.Encrypt(ctx, keyID, plaintext)
		if err == nil {
			c.usedLocal = false
			return ciphertext, nil
		}
		// Fall through to fallback on error
	}

	// Use fallback (local)
	if c.fallback == nil {
		return nil, fmt.Errorf("no crypto provider available")
	}

	c.usedLocal = true
	return c.fallback.Encrypt(ctx, keyID, plaintext)
}

// Decrypt implements CryptoProvider.
func (c *CompositeCryptoProvider) Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Try primary (companion) if connected
	if c.checker != nil && c.checker.IsConnected() {
		plaintext, err := c.primary.Decrypt(ctx, keyID, ciphertext)
		if err == nil {
			c.usedLocal = false
			return plaintext, nil
		}
		// Fall through to fallback on error
	}

	// Use fallback (local)
	if c.fallback == nil {
		return nil, fmt.Errorf("no crypto provider available")
	}

	c.usedLocal = true
	return c.fallback.Decrypt(ctx, keyID, ciphertext)
}
