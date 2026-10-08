package dispatch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTake(t *testing.T) {
	valid, err := Encode(testEvent(), "")
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		content []byte
		wantErr bool
	}{
		{"valid payload", valid, false},
		{"invalid payload", []byte("not json"), true},
		{"missing file", nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "mgc-telemetry-test.json")
			if tc.content != nil {
				if err := os.WriteFile(path, tc.content, 0o600); err != nil {
					t.Fatal(err)
				}
			}

			payload, err := Take(path)

			if (err != nil) != tc.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && payload.Event.Action != testEvent().Action {
				t.Errorf("action = %q, want %q", payload.Event.Action, testEvent().Action)
			}
			if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
				t.Errorf("Take must always delete the file, stat err = %v", statErr)
			}
		})
	}
}
