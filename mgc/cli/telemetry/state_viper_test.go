package telemetry

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestViperStateStoreLoad(t *testing.T) {
	testCases := []struct {
		name    string
		content string // "" significa que o arquivo não existe
		wantErr bool
	}{
		{"missing file loads the zero state and creates nothing", "", false},
		{"corrupted file reports an error and falls back to the zero state", "telemetry: [not: valid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, StateFileName)
			if tc.content != "" {
				if err := os.WriteFile(path, []byte(tc.content), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			state, err := NewViperStateStore(dir).Load()

			if (err != nil) != tc.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if state != (State{}) {
				t.Errorf("want the zero state, got %+v", state)
			}
			if _, statErr := os.Stat(path); tc.content == "" && !os.IsNotExist(statErr) {
				t.Errorf("Load must not create the state file, stat err = %v", statErr)
			}
		})
	}
}

func TestViperStateStoreSave(t *testing.T) {
	setAt := time.Date(2026, 9, 22, 14, 30, 45, 123000000, time.UTC)
	full := State{
		Disabled:           true,
		NoticeShown:        true,
		CredentialsSetAt:   &setAt,
		FirstValueRecorded: true,
		InstallationID:     "b3f1c2de-0000-4000-8000-000000000000",
	}

	testCases := []struct {
		name  string
		saves []State
	}{
		{"round trip with every key", []State{full}},
		{"repeated saves leave no temporary files", []State{{NoticeShown: true}, {NoticeShown: true}, {NoticeShown: true}}},
		{"last save wins", []State{full, {NoticeShown: true}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			store := NewViperStateStore(dir)

			for _, s := range tc.saves {
				if err := store.Save(s); err != nil {
					t.Fatal(err)
				}
			}

			got, err := store.Load()
			if err != nil {
				t.Fatal(err)
			}
			if want := tc.saves[len(tc.saves)-1]; !sameState(got, want) {
				t.Errorf("got %+v, want %+v", got, want)
			}

			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 || entries[0].Name() != StateFileName {
				t.Errorf("only %s must remain in %s, got %d entries", StateFileName, dir, len(entries))
			}
		})
	}
}

func TestViperStateStoreUnwritableDirectory(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}

	dir := t.TempDir()
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	if err := NewViperStateStore(dir).Save(State{Disabled: true}); err == nil {
		t.Error("expected an error when the directory is not writable")
	}
}

func TestViperStateStoreConcurrentSavesNeverCorrupt(t *testing.T) {
	store := NewViperStateStore(t.TempDir())
	want := State{NoticeShown: true, InstallationID: "3f1c2b9e-0000-4000-8000-000000000000"}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.Save(want)
		}()
	}
	wg.Wait()

	got, err := store.Load()
	if err != nil {
		t.Fatalf("state file must stay readable after concurrent saves: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func sameState(a, b State) bool {
	if (a.CredentialsSetAt == nil) != (b.CredentialsSetAt == nil) {
		return false
	}
	if a.CredentialsSetAt != nil && !a.CredentialsSetAt.Equal(*b.CredentialsSetAt) {
		return false
	}
	a.CredentialsSetAt, b.CredentialsSetAt = nil, nil
	return a == b
}
