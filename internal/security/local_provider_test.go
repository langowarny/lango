package security

import (
	"context"
	"testing"
)

func TestLocalCryptoProvider_Initialize(t *testing.T) {
	p := NewLocalCryptoProvider()

	// Test short passphrase
	err := p.Initialize("short")
	if err == nil {
		t.Error("expected error for short passphrase")
	}

	// Test valid passphrase
	err = p.Initialize("secure-passphrase-123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !p.IsInitialized() {
		t.Error("expected provider to be initialized")
	}

	if len(p.Salt()) != SaltSize {
		t.Errorf("expected salt size %d, got %d", SaltSize, len(p.Salt()))
	}
}

func TestLocalCryptoProvider_EncryptDecrypt(t *testing.T) {
	p := NewLocalCryptoProvider()
	if err := p.Initialize("test-passphrase-123"); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	plaintext := []byte("secret message to encrypt")

	// Encrypt
	ciphertext, err := p.Encrypt(ctx, "local", plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if len(ciphertext) <= len(plaintext) {
		t.Error("ciphertext should be longer than plaintext")
	}

	// Decrypt
	decrypted, err := p.Decrypt(ctx, "local", ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted text mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestLocalCryptoProvider_DecryptWithWrongKey(t *testing.T) {
	p1 := NewLocalCryptoProvider()
	if err := p1.Initialize("passphrase-one-123"); err != nil {
		t.Fatal(err)
	}

	p2 := NewLocalCryptoProvider()
	if err := p2.Initialize("passphrase-two-456"); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	plaintext := []byte("secret message")

	// Encrypt with p1
	ciphertext, err := p1.Encrypt(ctx, "local", plaintext)
	if err != nil {
		t.Fatal(err)
	}

	// Try to decrypt with p2 - should fail
	_, err = p2.Decrypt(ctx, "local", ciphertext)
	if err == nil {
		t.Error("expected decryption to fail with wrong key")
	}
}

func TestLocalCryptoProvider_Sign(t *testing.T) {
	p := NewLocalCryptoProvider()
	if err := p.Initialize("test-passphrase-123"); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	payload := []byte("data to sign")

	sig1, err := p.Sign(ctx, "local", payload)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}

	// Same payload should produce same signature
	sig2, err := p.Sign(ctx, "local", payload)
	if err != nil {
		t.Fatal(err)
	}

	if string(sig1) != string(sig2) {
		t.Error("signatures should match for same payload")
	}

	// Different payload should produce different signature
	sig3, err := p.Sign(ctx, "local", []byte("different data"))
	if err != nil {
		t.Fatal(err)
	}

	if string(sig1) == string(sig3) {
		t.Error("signatures should differ for different payloads")
	}
}

func TestLocalCryptoProvider_NotInitialized(t *testing.T) {
	p := NewLocalCryptoProvider()
	ctx := context.Background()

	_, err := p.Encrypt(ctx, "local", []byte("test"))
	if err == nil {
		t.Error("expected error for uninitialized provider")
	}

	_, err = p.Decrypt(ctx, "local", []byte("test"))
	if err == nil {
		t.Error("expected error for uninitialized provider")
	}

	_, err = p.Sign(ctx, "local", []byte("test"))
	if err == nil {
		t.Error("expected error for uninitialized provider")
	}
}

func TestLocalCryptoProvider_InitializeWithSalt(t *testing.T) {
	p1 := NewLocalCryptoProvider()
	passphrase := "test-passphrase-123"
	if err := p1.Initialize(passphrase); err != nil {
		t.Fatal(err)
	}

	salt := p1.Salt()
	ctx := context.Background()
	plaintext := []byte("secret message")

	ciphertext, err := p1.Encrypt(ctx, "local", plaintext)
	if err != nil {
		t.Fatal(err)
	}

	// Create new provider with same passphrase and salt
	p2 := NewLocalCryptoProvider()
	if err := p2.InitializeWithSalt(passphrase, salt); err != nil {
		t.Fatal(err)
	}

	// Should be able to decrypt
	decrypted, err := p2.Decrypt(ctx, "local", ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Error("decrypted text mismatch")
	}
}
