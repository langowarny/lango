package browser

import (
	"errors"
	"fmt"
	"testing"
)

func TestSafeRodCall_RecoversPanic(t *testing.T) {
	err := safeRodCall(func() error {
		panic("CDP connection lost")
	})
	if err == nil {
		t.Fatal("expected error from panic, got nil")
	}
	if !errors.Is(err, ErrBrowserPanic) {
		t.Errorf("expected ErrBrowserPanic, got %v", err)
	}
	if want := "CDP connection lost"; !containsStr(err.Error(), want) {
		t.Errorf("error should contain %q, got %q", want, err.Error())
	}
}

func TestSafeRodCall_PassesNormalError(t *testing.T) {
	sentinel := fmt.Errorf("normal error")
	err := safeRodCall(func() error {
		return sentinel
	})
	if err != sentinel {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestSafeRodCall_ReturnsNilOnSuccess(t *testing.T) {
	err := safeRodCall(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestSafeRodCallValue_RecoversPanic(t *testing.T) {
	val, err := safeRodCallValue(func() (string, error) {
		panic("websocket closed")
	})
	if err == nil {
		t.Fatal("expected error from panic, got nil")
	}
	if !errors.Is(err, ErrBrowserPanic) {
		t.Errorf("expected ErrBrowserPanic, got %v", err)
	}
	if val != "" {
		t.Errorf("expected zero value on panic, got %q", val)
	}
}

func TestSafeRodCallValue_PassesNormalError(t *testing.T) {
	sentinel := fmt.Errorf("element not found")
	val, err := safeRodCallValue(func() (int, error) {
		return 0, sentinel
	})
	if err != sentinel {
		t.Errorf("expected sentinel error, got %v", err)
	}
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}
}

func TestSafeRodCallValue_ReturnsValueOnSuccess(t *testing.T) {
	val, err := safeRodCallValue(func() (string, error) {
		return "hello", nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if val != "hello" {
		t.Errorf("expected %q, got %q", "hello", val)
	}
}

func TestErrBrowserPanic_Unwrap(t *testing.T) {
	wrapped := fmt.Errorf("%w: runtime crash", ErrBrowserPanic)
	if !errors.Is(wrapped, ErrBrowserPanic) {
		t.Error("wrapped error should match ErrBrowserPanic via errors.Is")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && searchStr(s, sub)
}

func searchStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
