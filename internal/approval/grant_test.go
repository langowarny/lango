package approval

import (
	"sync"
	"testing"
	"time"
)

func TestGrantStore_GrantAndIsGranted(t *testing.T) {
	gs := NewGrantStore()

	if gs.IsGranted("session-1", "exec") {
		t.Error("expected no grant before granting")
	}

	gs.Grant("session-1", "exec")

	if !gs.IsGranted("session-1", "exec") {
		t.Error("expected grant after granting")
	}
}

func TestGrantStore_GrantIsolation(t *testing.T) {
	gs := NewGrantStore()
	gs.Grant("session-1", "exec")

	tests := []struct {
		give       string
		giveTool   string
		wantResult bool
	}{
		{give: "session-1", giveTool: "exec", wantResult: true},
		{give: "session-1", giveTool: "fs_delete", wantResult: false},
		{give: "session-2", giveTool: "exec", wantResult: false},
	}

	for _, tt := range tests {
		t.Run(tt.give+":"+tt.giveTool, func(t *testing.T) {
			if got := gs.IsGranted(tt.give, tt.giveTool); got != tt.wantResult {
				t.Errorf("IsGranted(%q, %q) = %v, want %v", tt.give, tt.giveTool, got, tt.wantResult)
			}
		})
	}
}

func TestGrantStore_Revoke(t *testing.T) {
	gs := NewGrantStore()
	gs.Grant("session-1", "exec")
	gs.Grant("session-1", "fs_write")

	gs.Revoke("session-1", "exec")

	if gs.IsGranted("session-1", "exec") {
		t.Error("expected exec grant to be revoked")
	}
	if !gs.IsGranted("session-1", "fs_write") {
		t.Error("expected fs_write grant to remain")
	}
}

func TestGrantStore_RevokeSession(t *testing.T) {
	gs := NewGrantStore()
	gs.Grant("session-1", "exec")
	gs.Grant("session-1", "fs_write")
	gs.Grant("session-2", "exec")

	gs.RevokeSession("session-1")

	if gs.IsGranted("session-1", "exec") {
		t.Error("expected session-1 exec grant to be revoked")
	}
	if gs.IsGranted("session-1", "fs_write") {
		t.Error("expected session-1 fs_write grant to be revoked")
	}
	if !gs.IsGranted("session-2", "exec") {
		t.Error("expected session-2 exec grant to remain")
	}
}

func TestGrantStore_ConcurrentAccess(t *testing.T) {
	gs := NewGrantStore()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gs.Grant("session-1", "exec")
			gs.IsGranted("session-1", "exec")
			gs.Grant("session-2", "fs_write")
			gs.Revoke("session-2", "fs_write")
		}()
	}

	wg.Wait()

	if !gs.IsGranted("session-1", "exec") {
		t.Error("expected session-1 exec grant after concurrent access")
	}
}

func TestGrantStore_RevokeNonExistent(t *testing.T) {
	gs := NewGrantStore()

	// Should not panic
	gs.Revoke("nonexistent", "tool")
	gs.RevokeSession("nonexistent")
}

func TestGrantStore_TTLExpired(t *testing.T) {
	now := time.Now()
	gs := NewGrantStore()
	gs.nowFn = func() time.Time { return now }
	gs.SetTTL(10 * time.Minute)

	gs.Grant("session-1", "echo")

	// Still valid within TTL.
	gs.nowFn = func() time.Time { return now.Add(9 * time.Minute) }
	if !gs.IsGranted("session-1", "echo") {
		t.Error("expected grant to be valid within TTL")
	}

	// Expired after TTL.
	gs.nowFn = func() time.Time { return now.Add(11 * time.Minute) }
	if gs.IsGranted("session-1", "echo") {
		t.Error("expected grant to be expired after TTL")
	}
}

func TestGrantStore_TTLZeroMeansNoExpiry(t *testing.T) {
	now := time.Now()
	gs := NewGrantStore()
	gs.nowFn = func() time.Time { return now }
	// TTL = 0 (default).

	gs.Grant("session-1", "echo")

	// 100 hours later, still valid.
	gs.nowFn = func() time.Time { return now.Add(100 * time.Hour) }
	if !gs.IsGranted("session-1", "echo") {
		t.Error("expected grant to be valid indefinitely when TTL = 0")
	}
}

func TestGrantStore_CleanExpired(t *testing.T) {
	now := time.Now()
	gs := NewGrantStore()
	gs.nowFn = func() time.Time { return now }
	gs.SetTTL(5 * time.Minute)

	gs.Grant("session-1", "echo")
	gs.Grant("session-1", "exec")
	gs.Grant("session-2", "echo")

	// Advance time past TTL for the first two, but grant session-2:echo later.
	gs.nowFn = func() time.Time { return now.Add(3 * time.Minute) }
	gs.Grant("session-2", "echo") // refresh

	gs.nowFn = func() time.Time { return now.Add(6 * time.Minute) }
	removed := gs.CleanExpired()
	if removed != 2 {
		t.Errorf("expected 2 expired grants removed, got %d", removed)
	}

	if gs.IsGranted("session-1", "echo") {
		t.Error("session-1:echo should be cleaned")
	}
	if gs.IsGranted("session-1", "exec") {
		t.Error("session-1:exec should be cleaned")
	}
	if !gs.IsGranted("session-2", "echo") {
		t.Error("session-2:echo should still be valid (refreshed)")
	}
}

func TestGrantStore_CleanExpiredNoOpWhenTTLZero(t *testing.T) {
	gs := NewGrantStore()
	gs.Grant("session-1", "echo")

	removed := gs.CleanExpired()
	if removed != 0 {
		t.Errorf("expected 0 removed with TTL=0, got %d", removed)
	}
	if !gs.IsGranted("session-1", "echo") {
		t.Error("grant should remain when TTL=0")
	}
}
