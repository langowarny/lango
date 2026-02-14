package agent

import (
	"sync"
	"testing"
)

func TestSecretScanner_Scan(t *testing.T) {
	tests := []struct {
		name       string
		giveNames  []string
		giveValues [][]byte
		giveText   string
		want       string
	}{
		{
			name:       "registered secret is masked",
			giveNames:  []string{"API_KEY"},
			giveValues: [][]byte{[]byte("super-secret-key-1234")},
			giveText:   "The key is super-secret-key-1234 here",
			want:       "The key is [SECRET:API_KEY] here",
		},
		{
			name:       "text without secrets is unchanged",
			giveNames:  []string{"DB_PASS"},
			giveValues: [][]byte{[]byte("hunter2hunter2")},
			giveText:   "nothing to see here",
			want:       "nothing to see here",
		},
		{
			name:       "multiple secrets in same text",
			giveNames:  []string{"TOKEN", "PASSWORD"},
			giveValues: [][]byte{[]byte("tok-abc123"), []byte("p@ssw0rd!")},
			giveText:   "token=tok-abc123 password=p@ssw0rd!",
			want:       "token=[SECRET:TOKEN] password=[SECRET:PASSWORD]",
		},
		{
			name:       "short value not masked",
			giveNames:  []string{"PIN"},
			giveValues: [][]byte{[]byte("123")},
			giveText:   "my pin is 123 ok",
			want:       "my pin is 123 ok",
		},
		{
			name:       "exactly 4 chars is masked",
			giveNames:  []string{"CODE"},
			giveValues: [][]byte{[]byte("abcd")},
			giveText:   "code=abcd end",
			want:       "code=[SECRET:CODE] end",
		},
		{
			name:     "empty scanner leaves text unchanged",
			giveText: "no secrets registered at all",
			want:     "no secrets registered at all",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewSecretScanner()
			for i := range tt.giveNames {
				scanner.Register(tt.giveNames[i], tt.giveValues[i])
			}

			got := scanner.Scan(tt.giveText)
			if got != tt.want {
				t.Errorf("Scan() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSecretScanner_Clear(t *testing.T) {
	scanner := NewSecretScanner()
	scanner.Register("SECRET", []byte("mysecretvalue"))

	if !scanner.HasSecrets() {
		t.Fatal("HasSecrets() = false after Register, want true")
	}

	got := scanner.Scan("mysecretvalue")
	if got != "[SECRET:SECRET]" {
		t.Errorf("Scan() before Clear = %q, want %q", got, "[SECRET:SECRET]")
	}

	scanner.Clear()

	if scanner.HasSecrets() {
		t.Fatal("HasSecrets() = true after Clear, want false")
	}

	got = scanner.Scan("mysecretvalue")
	if got != "mysecretvalue" {
		t.Errorf("Scan() after Clear = %q, want %q", got, "mysecretvalue")
	}
}

func TestSecretScanner_HasSecrets(t *testing.T) {
	scanner := NewSecretScanner()

	if scanner.HasSecrets() {
		t.Fatal("HasSecrets() = true on new scanner, want false")
	}

	scanner.Register("KEY", []byte("longvalue"))

	if !scanner.HasSecrets() {
		t.Fatal("HasSecrets() = false after Register, want true")
	}

	// Short value should be ignored, count should remain at 1.
	scanner.Register("SHORT", []byte("ab"))

	scanner.Clear()

	if scanner.HasSecrets() {
		t.Fatal("HasSecrets() = true after Clear, want false")
	}
}

func TestSecretScanner_ConcurrentAccess(t *testing.T) {
	scanner := NewSecretScanner()
	scanner.Register("BASE", []byte("base-secret-value"))

	var wg sync.WaitGroup
	const goroutines = 50

	// Half the goroutines register new secrets.
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			scanner.Register("KEY_"+string(rune('A'+id%26)),
				[]byte("secret-value-for-goroutine"))
		}(i)
	}

	// Other half scans concurrently.
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = scanner.Scan("text containing base-secret-value here")
		}()
	}

	wg.Wait()

	// If we reach here without a race detector failure, concurrency is safe.
	if !scanner.HasSecrets() {
		t.Fatal("HasSecrets() = false, want true after concurrent registers")
	}
}
