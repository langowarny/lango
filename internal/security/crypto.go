package security

import (
	"context"
)

// CryptoProvider defines the interface for cryptographic operations.
type CryptoProvider interface {
	// Sign generates a signature for the given payload using the specified key ID.
	Sign(ctx context.Context, keyID string, payload []byte) ([]byte, error)
	// Encrypt encrypts the given plaintext using the specified key ID.
	Encrypt(ctx context.Context, keyID string, plaintext []byte) ([]byte, error)
	// Decrypt decrypts the given ciphertext using the specified key ID.
	Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error)
}
