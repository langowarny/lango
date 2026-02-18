package security

import (
	"context"
	"testing"
	"time"
)

func TestRPCProvider_Sign(t *testing.T) {
	provider := NewRPCProvider()

	// Mock Sender that replies immediately
	provider.SetSender(func(event string, payload interface{}) error {
		if event != "sign.request" {
			t.Errorf("expected sign.request, got %s", event)
			return nil
		}
		req := payload.(SignRequest)

		// Simulate response
		resp := SignResponse{
			ID:        req.ID,
			Signature: []byte("signature_bytes"),
		}

		// Handle response in a goroutine to avoid blocking if the channel buffer was 0 (it's 1, but good practice)
		go func() {
			if err := provider.HandleSignResponse(resp); err != nil {
				t.Errorf("handle response failed: %v", err)
			}
		}()
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sig, err := provider.Sign(ctx, "key1", []byte("data"))
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if string(sig) != "signature_bytes" {
		t.Errorf("expected 'signature_bytes', got %s", string(sig))
	}
}

func TestRPCProvider_Encrypt(t *testing.T) {
	provider := NewRPCProvider()

	provider.SetSender(func(event string, payload interface{}) error {
		if event != "encrypt.request" {
			t.Errorf("expected encrypt.request, got %s", event)
			return nil
		}
		req := payload.(EncryptRequest)

		resp := EncryptResponse{
			ID:         req.ID,
			Ciphertext: []byte("encrypted_bytes"),
		}

		go func() {
			provider.HandleEncryptResponse(resp)
		}()
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cipher, err := provider.Encrypt(ctx, "key1", []byte("plaintext"))
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if string(cipher) != "encrypted_bytes" {
		t.Errorf("expected 'encrypted_bytes', got %s", string(cipher))
	}
}

func TestRPCProvider_Decrypt(t *testing.T) {
	provider := NewRPCProvider()

	provider.SetSender(func(event string, payload interface{}) error {
		if event != "decrypt.request" {
			t.Errorf("expected decrypt.request, got %s", event)
			return nil
		}
		req := payload.(DecryptRequest)

		resp := DecryptResponse{
			ID:        req.ID,
			Plaintext: []byte("decrypted_bytes"),
		}

		go func() {
			provider.HandleDecryptResponse(resp)
		}()
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	plain, err := provider.Decrypt(ctx, "key1", []byte("ciphertext"))
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(plain) != "decrypted_bytes" {
		t.Errorf("expected 'decrypted_bytes', got %s", string(plain))
	}
}
