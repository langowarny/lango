package security

import (
	"context"
	"testing"
)

// mockConnectionChecker for testing
type mockConnectionChecker struct {
	connected bool
}

func (m *mockConnectionChecker) IsConnected() bool {
	return m.connected
}

// mockCryptoProvider for testing
type mockCryptoProvider struct {
	signResult    []byte
	encryptResult []byte
	decryptResult []byte
	signErr       error
	encryptErr    error
	decryptErr    error
	called        bool
}

func (m *mockCryptoProvider) Sign(ctx context.Context, keyID string, payload []byte) ([]byte, error) {
	m.called = true
	return m.signResult, m.signErr
}

func (m *mockCryptoProvider) Encrypt(ctx context.Context, keyID string, plaintext []byte) ([]byte, error) {
	m.called = true
	return m.encryptResult, m.encryptErr
}

func (m *mockCryptoProvider) Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error) {
	m.called = true
	return m.decryptResult, m.decryptErr
}

func TestCompositeProvider_UsesPrimaryWhenConnected(t *testing.T) {
	primary := &mockCryptoProvider{encryptResult: []byte("primary-encrypted")}
	fallback := &mockCryptoProvider{encryptResult: []byte("fallback-encrypted")}
	checker := &mockConnectionChecker{connected: true}

	composite := NewCompositeCryptoProvider(primary, fallback, checker)

	result, err := composite.Encrypt(context.Background(), "key1", []byte("data"))
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != "primary-encrypted" {
		t.Errorf("expected primary result, got %s", result)
	}

	if !primary.called {
		t.Error("primary should have been called")
	}

	if fallback.called {
		t.Error("fallback should not have been called")
	}

	if composite.UsedLocal() {
		t.Error("should not have used local")
	}
}

func TestCompositeProvider_UsesFallbackWhenDisconnected(t *testing.T) {
	primary := &mockCryptoProvider{encryptResult: []byte("primary-encrypted")}
	fallback := &mockCryptoProvider{encryptResult: []byte("fallback-encrypted")}
	checker := &mockConnectionChecker{connected: false}

	composite := NewCompositeCryptoProvider(primary, fallback, checker)

	result, err := composite.Encrypt(context.Background(), "key1", []byte("data"))
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != "fallback-encrypted" {
		t.Errorf("expected fallback result, got %s", result)
	}

	if primary.called {
		t.Error("primary should not have been called")
	}

	if !fallback.called {
		t.Error("fallback should have been called")
	}

	if !composite.UsedLocal() {
		t.Error("should have used local")
	}
}

func TestCompositeProvider_ErrorsWhenNoProvider(t *testing.T) {
	checker := &mockConnectionChecker{connected: false}
	composite := NewCompositeCryptoProvider(nil, nil, checker)

	_, err := composite.Encrypt(context.Background(), "key1", []byte("data"))
	if err == nil {
		t.Error("expected error when no provider available")
	}
}

func TestCompositeProvider_Sign(t *testing.T) {
	primary := &mockCryptoProvider{signResult: []byte("primary-sig")}
	fallback := &mockCryptoProvider{signResult: []byte("fallback-sig")}
	checker := &mockConnectionChecker{connected: true}

	composite := NewCompositeCryptoProvider(primary, fallback, checker)

	result, err := composite.Sign(context.Background(), "key1", []byte("data"))
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != "primary-sig" {
		t.Errorf("expected primary signature, got %s", result)
	}
}

func TestCompositeProvider_Decrypt(t *testing.T) {
	fallback := &mockCryptoProvider{decryptResult: []byte("decrypted-data")}
	checker := &mockConnectionChecker{connected: false}

	composite := NewCompositeCryptoProvider(nil, fallback, checker)

	result, err := composite.Decrypt(context.Background(), "key1", []byte("encrypted"))
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != "decrypted-data" {
		t.Errorf("expected decrypted data, got %s", result)
	}
}
