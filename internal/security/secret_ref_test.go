package security

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestRefStore_Store(t *testing.T) {
	tests := []struct {
		give      string
		giveValue []byte
		want      string
	}{
		{
			give:      "api-key",
			giveValue: []byte("sk-12345"),
			want:      "{{secret:api-key}}",
		},
		{
			give:      "db-password",
			giveValue: []byte("p@ssw0rd"),
			want:      "{{secret:db-password}}",
		},
		{
			give:      "empty",
			giveValue: []byte(""),
			want:      "{{secret:empty}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			rs := NewRefStore()
			got := rs.Store(tt.give, tt.giveValue)
			if got != tt.want {
				t.Errorf("Store(%q) = %q, want %q", tt.give, got, tt.want)
			}
		})
	}
}

func TestRefStore_StoreDecrypted(t *testing.T) {
	tests := []struct {
		give      string
		giveValue []byte
		want      string
	}{
		{
			give:      "abc-123",
			giveValue: []byte("decrypted-data"),
			want:      "{{decrypt:abc-123}}",
		},
		{
			give:      "session-42",
			giveValue: []byte("session-payload"),
			want:      "{{decrypt:session-42}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			rs := NewRefStore()
			got := rs.StoreDecrypted(tt.give, tt.giveValue)
			if got != tt.want {
				t.Errorf("StoreDecrypted(%q) = %q, want %q",
					tt.give, got, tt.want)
			}
		})
	}
}

func TestRefStore_Resolve(t *testing.T) {
	rs := NewRefStore()
	rs.Store("my-secret", []byte("secret-value"))
	rs.StoreDecrypted("dec-1", []byte("decrypted-value"))

	tests := []struct {
		give     string
		wantVal  []byte
		wantOK   bool
	}{
		{
			give:    "{{secret:my-secret}}",
			wantVal: []byte("secret-value"),
			wantOK:  true,
		},
		{
			give:    "{{decrypt:dec-1}}",
			wantVal: []byte("decrypted-value"),
			wantOK:  true,
		},
		{
			give:    "{{secret:unknown}}",
			wantVal: nil,
			wantOK:  false,
		},
		{
			give:    "not-a-token",
			wantVal: nil,
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			gotVal, gotOK := rs.Resolve(tt.give)
			if gotOK != tt.wantOK {
				t.Errorf("Resolve(%q) ok = %v, want %v",
					tt.give, gotOK, tt.wantOK)
			}
			if !bytes.Equal(gotVal, tt.wantVal) {
				t.Errorf("Resolve(%q) val = %q, want %q",
					tt.give, gotVal, tt.wantVal)
			}
		})
	}
}

func TestRefStore_Resolve_DoesNotMutateInternal(t *testing.T) {
	rs := NewRefStore()
	rs.Store("key", []byte("original"))

	val, ok := rs.Resolve("{{secret:key}}")
	if !ok {
		t.Fatal("expected to resolve stored secret")
	}

	// Mutate the returned slice
	val[0] = 'X'

	// Verify internal state is unchanged
	val2, _ := rs.Resolve("{{secret:key}}")
	if string(val2) != "original" {
		t.Errorf("internal state mutated: got %q, want %q", val2, "original")
	}
}

func TestRefStore_ResolveAll(t *testing.T) {
	rs := NewRefStore()
	rs.Store("user", []byte("admin"))
	rs.Store("pass", []byte("s3cret"))
	rs.StoreDecrypted("data-1", []byte("payload"))

	tests := []struct {
		give string
		want string
	}{
		{
			give: "login with {{secret:user}} and {{secret:pass}}",
			want: "login with admin and s3cret",
		},
		{
			give: "result: {{decrypt:data-1}}",
			want: "result: payload",
		},
		{
			give: "no tokens here",
			want: "no tokens here",
		},
		{
			give: "unknown {{secret:missing}} stays",
			want: "unknown {{secret:missing}} stays",
		},
		{
			give: "",
			want: "",
		},
		{
			give: "mixed {{secret:user}} and {{decrypt:data-1}} tokens",
			want: "mixed admin and payload tokens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			got := rs.ResolveAll(tt.give)
			if got != tt.want {
				t.Errorf("ResolveAll(%q) = %q, want %q",
					tt.give, got, tt.want)
			}
		})
	}
}

func TestRefStore_Values(t *testing.T) {
	rs := NewRefStore()
	rs.Store("a", []byte("val-a"))
	rs.Store("b", []byte("val-b"))
	rs.StoreDecrypted("c", []byte("val-c"))

	values := rs.Values()
	if len(values) != 3 {
		t.Fatalf("Values() returned %d items, want 3", len(values))
	}

	// Collect all values into a set for order-independent comparison
	seen := make(map[string]bool, len(values))
	for _, v := range values {
		seen[string(v)] = true
	}

	for _, want := range []string{"val-a", "val-b", "val-c"} {
		if !seen[want] {
			t.Errorf("Values() missing %q", want)
		}
	}
}

func TestRefStore_Values_DoesNotMutateInternal(t *testing.T) {
	rs := NewRefStore()
	rs.Store("key", []byte("original"))

	values := rs.Values()
	if len(values) != 1 {
		t.Fatalf("Values() returned %d items, want 1", len(values))
	}

	// Mutate the returned slice
	values[0][0] = 'X'

	// Verify internal state is unchanged
	val, _ := rs.Resolve("{{secret:key}}")
	if string(val) != "original" {
		t.Errorf("internal state mutated: got %q, want %q", val, "original")
	}
}

func TestRefStore_Names(t *testing.T) {
	rs := NewRefStore()
	rs.Store("api-key", []byte("sk-12345"))
	rs.Store("db-pass", []byte("p@ss"))
	rs.StoreDecrypted("session", []byte("sess-data"))

	names := rs.Names()

	tests := []struct {
		givePlaintext string
		wantName      string
	}{
		{
			givePlaintext: "sk-12345",
			wantName:      "api-key",
		},
		{
			givePlaintext: "p@ss",
			wantName:      "db-pass",
		},
		{
			givePlaintext: "sess-data",
			wantName:      "session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.wantName, func(t *testing.T) {
			got, ok := names[tt.givePlaintext]
			if !ok {
				t.Errorf("Names() missing key %q", tt.givePlaintext)
				return
			}
			if got != tt.wantName {
				t.Errorf("Names()[%q] = %q, want %q",
					tt.givePlaintext, got, tt.wantName)
			}
		})
	}
}

func TestRefStore_Clear(t *testing.T) {
	rs := NewRefStore()
	rs.Store("secret-1", []byte("value-1"))
	rs.StoreDecrypted("dec-1", []byte("value-2"))

	// Verify data exists before clearing
	if _, ok := rs.Resolve("{{secret:secret-1}}"); !ok {
		t.Fatal("expected secret to exist before Clear")
	}

	rs.Clear()

	// Verify all data is gone
	if _, ok := rs.Resolve("{{secret:secret-1}}"); ok {
		t.Error("expected secret to be cleared")
	}
	if _, ok := rs.Resolve("{{decrypt:dec-1}}"); ok {
		t.Error("expected decrypt ref to be cleared")
	}
	if len(rs.Values()) != 0 {
		t.Error("expected Values() to be empty after Clear")
	}
	if len(rs.Names()) != 0 {
		t.Error("expected Names() to be empty after Clear")
	}
}

func TestRefStore_ConcurrentStoreAndResolve(t *testing.T) {
	rs := NewRefStore()
	const numGoroutines = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	// Concurrent stores
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			name := fmt.Sprintf("key-%d", idx)
			value := []byte(fmt.Sprintf("val-%d", idx))
			rs.Store(name, value)
		}(i)
	}

	// Concurrent resolves (some will find values, some won't)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			token := fmt.Sprintf("{{secret:key-%d}}", idx)
			rs.Resolve(token)
		}(i)
	}

	wg.Wait()

	// After all goroutines finish, verify all values were stored
	for i := 0; i < numGoroutines; i++ {
		token := fmt.Sprintf("{{secret:key-%d}}", i)
		val, ok := rs.Resolve(token)
		if !ok {
			t.Errorf("missing value for %s after concurrent store", token)
			continue
		}
		want := fmt.Sprintf("val-%d", i)
		if string(val) != want {
			t.Errorf("Resolve(%q) = %q, want %q", token, val, want)
		}
	}
}

func TestRefStore_ConcurrentMixedOperations(t *testing.T) {
	rs := NewRefStore()
	rs.Store("pre-existing", []byte("initial"))

	var wg sync.WaitGroup
	wg.Add(4)

	// Store
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			rs.Store(fmt.Sprintf("s-%d", i), []byte(fmt.Sprintf("v-%d", i)))
		}
	}()

	// StoreDecrypted
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			rs.StoreDecrypted(fmt.Sprintf("d-%d", i),
				[]byte(fmt.Sprintf("dv-%d", i)))
		}
	}()

	// ResolveAll
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			rs.ResolveAll("test {{secret:pre-existing}} input")
		}
	}()

	// Values + Names
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			rs.Values()
			rs.Names()
		}
	}()

	wg.Wait()

	// Verify pre-existing value survives concurrent operations
	val, ok := rs.Resolve("{{secret:pre-existing}}")
	if !ok {
		t.Error("pre-existing secret lost during concurrent operations")
	}
	if string(val) != "initial" {
		t.Errorf("pre-existing value changed: got %q, want %q",
			val, "initial")
	}
}
