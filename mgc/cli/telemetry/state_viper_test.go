package telemetry

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestStateRoundTrip(t *testing.T) {
	store := NewViperStateStore(t.TempDir())

	empty, err := store.Load()
	if err != nil {
		t.Fatalf("missing file must not be an error: %v", err)
	}
	if empty != (State{}) {
		t.Fatalf("missing file must load the zero state, got %+v", empty)
	}

	setAt := time.Date(2026, 9, 22, 14, 30, 45, 123000000, time.UTC)
	want := State{
		Disabled:           true,
		NoticeShown:        true,
		CredentialsSetAt:   &setAt,
		FirstValueRecorded: true,
		InstallationID:     "b3f1c2de-0000-4000-8000-000000000000",
	}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Disabled != want.Disabled || got.NoticeShown != want.NoticeShown ||
		got.FirstValueRecorded != want.FirstValueRecorded || got.InstallationID != want.InstallationID ||
		got.CredentialsSetAt == nil || !got.CredentialsSetAt.Equal(setAt) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestStateLoadDoesNotCreateFile(t *testing.T) {
	dir := t.TempDir()
	store := NewViperStateStore(dir)

	if _, err := store.Load(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, StateFileName)); !os.IsNotExist(err) {
		t.Errorf("Load must not create the state file, stat err = %v", err)
	}
}

func TestStateCorruptedFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, StateFileName), []byte("telemetry: [not: valid"), 0o644); err != nil {
		t.Fatal(err)
	}

	state, err := NewViperStateStore(dir).Load()
	if err == nil {
		t.Error("corrupted file should report an error to be logged in debug mode")
	}
	if state != (State{}) {
		t.Errorf("corrupted file must fall back to the zero state, got %+v", state)
	}
}

func TestStateUnwritableDirectory(t *testing.T) {
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

func TestStateSaveLeavesNoTemporaryFiles(t *testing.T) {
	dir := t.TempDir()
	store := NewViperStateStore(dir)

	for i := 0; i < 3; i++ {
		if err := store.Save(State{NoticeShown: true}); err != nil {
			t.Fatal(err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != StateFileName {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("only %s must remain, got %v", StateFileName, names)
	}
}

func TestStateConcurrentSavesNeverCorrupt(t *testing.T) {
	dir := t.TempDir()
	store := NewViperStateStore(dir)
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
