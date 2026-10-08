package dispatch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	helperEnv       = "MGC_TELEMETRY_DISPATCH_HELPER"
	helperSignalEnv = "MGC_TELEMETRY_DISPATCH_SIGNAL"
	helperDelay     = 2 * time.Second
	helperWait      = 10 * time.Second
	helperPoll      = 50 * time.Millisecond
	dispatchBudget  = time.Second
)

// TestMain deixa o binário de teste fazer o papel do processo de envio quando é
// executado de novo pelo Dispatch, como o binário da CLI faz com o SenderArg
func TestMain(m *testing.M) {
	if os.Getenv(helperEnv) == "1" && len(os.Args) == 3 && os.Args[1] == SenderArg {
		runHelperSender(os.Args[2])
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func runHelperSender(payloadPath string) {
	payload, err := Take(payloadPath)
	if err != nil {
		return
	}
	time.Sleep(helperDelay)
	_ = os.WriteFile(os.Getenv(helperSignalEnv), []byte(payload.Event.Action), payloadFileMode)
}

// roda o teste real e verifica se o arquivo relamente foi apagado
func TestSubprocessDetached(t *testing.T) {
	if testing.Short() {
		t.Skip("starts a real process")
	}

	tempDir := t.TempDir()
	signal := filepath.Join(t.TempDir(), "sent")
	t.Setenv(helperEnv, "1")
	t.Setenv(helperSignalEnv, signal)

	sub := NewSubprocess(nil, nil)
	sub.tempDir = func() string { return tempDir }

	start := time.Now()
	if err := sub.Dispatch(testEvent()); err != nil {
		t.Fatalf("Dispatch() = %v", err)
	}
	if elapsed := time.Since(start); elapsed > dispatchBudget {
		t.Errorf("Dispatch took %v, it must not wait for the sender", elapsed)
	}
	if _, err := os.Stat(signal); !os.IsNotExist(err) {
		t.Fatalf("the sender finished before Dispatch returned, stat err = %v", err)
	}

	deadline := time.Now().Add(helperWait)
	for time.Now().Before(deadline) {
		if content, err := os.ReadFile(signal); err == nil {
			if string(content) != testEvent().Action {
				t.Errorf("sender got action %q, want %q", content, testEvent().Action)
			}
			entries, _ := os.ReadDir(tempDir)
			if len(entries) != 0 {
				t.Errorf("the sender must delete the payload file, left %d entries", len(entries))
			}
			return
		}
		time.Sleep(helperPoll)
	}
	t.Fatalf("the sender did not finish within %v after Dispatch returned", helperWait)
}
