package common

import (
	"testing"
)

func TestFixFilenameEncoding(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantFixed     string
		wantConverted bool
		wantErr       bool
	}{
		{
			name:          "already valid UTF-8 is returned unchanged",
			input:         "café-relatório.txt",
			wantFixed:     "café-relatório.txt",
			wantConverted: false,
		},
		{
			name:          "windows-1252 'ô' byte (0xF4) is converted to UTF-8",
			input:         "rapport-\xf4.txt",
			wantFixed:     "rapport-ô.txt",
			wantConverted: true,
		},
		{
			name:          "windows-1252 'ê' byte (0xEA) is converted to UTF-8",
			input:         "fen\xeatre.txt",
			wantFixed:     "fenêtre.txt",
			wantConverted: true,
		},
		{
			name:          "ascii name with a single invalid high byte in the middle",
			input:         "report\xe9final.txt",
			wantFixed:     "reportéfinal.txt",
			wantConverted: true,
		},
		{
			name:          "byte undefined in windows-1252 falls back to windows-1251 (cyrillic)",
			input:         "test-\x81-name.txt",
			wantFixed:     "test-Ѓ-name.txt",
			wantConverted: true,
		},
		{
			name:          "byte undefined in windows-1252 falls back to windows-1250 (central european)",
			input:         "test-\x8d-name.txt",
			wantFixed:     "test-Ť-name.txt",
			wantConverted: true,
		},
		{
			name:          "byte plausible in both windows-1252 and windows-1251 prefers windows-1252 (priority order)",
			input:         "cost-\x80-report.txt",
			wantFixed:     "cost-€-report.txt",
			wantConverted: true,
		},
		{
			name:    "byte with no plausible mapping in any candidate encoding is rejected",
			input:   "broken-\x81\x01-name.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixed, converted, err := FixFilenameEncoding(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("FixFilenameEncoding(%q) expected an error, got none (fixed=%q)", tt.input, fixed)
				}
				return
			}

			if err != nil {
				t.Fatalf("FixFilenameEncoding(%q) unexpected error: %v", tt.input, err)
			}
			if fixed != tt.wantFixed {
				t.Errorf("FixFilenameEncoding(%q) fixed = %q, want %q", tt.input, fixed, tt.wantFixed)
			}
			if converted != tt.wantConverted {
				t.Errorf("FixFilenameEncoding(%q) converted = %v, want %v", tt.input, converted, tt.wantConverted)
			}
		})
	}
}
