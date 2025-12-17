package buckets

import (
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		size     float64
		expected string
	}{
		{
			name:     "zero bytes",
			size:     0,
			expected: "0.0B",
		},
		{
			name:     "less than 1 KiB",
			size:     512,
			expected: "512.0B",
		},
		{
			name:     "exactly 1 KiB",
			size:     1024,
			expected: "1.0 KiB",
		},
		{
			name:     "exactly 2 KiB",
			size:     2048,
			expected: "2.0 KiB",
		},
		{
			name:     "exactly 1 MiB",
			size:     1048576, // 1024 * 1024
			expected: "1.0 MiB",
		},
		{
			name:     "exactly 1 GiB",
			size:     1073741824, // 1024 * 1024 * 1024
			expected: "1.0 GiB",
		},
		{
			name:     "between KiB and MiB",
			size:     1536, // 1.5 KiB
			expected: "1.5 KiB",
		},
		{
			name:     "between MiB and GiB",
			size:     1572864, // 1.5 MiB
			expected: "1.5 MiB",
		},
		{
			name:     "just below 1 KiB",
			size:     1023,
			expected: "1023.0B",
		},
		{
			name:     "just above 1 KiB",
			size:     1025,
			expected: "1.0 KiB",
		},
		{
			name:     "exactly 1 TiB",
			size:     1099511627776, // 1024^4
			expected: "1.0 TiB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSize(tt.size)
			if got != tt.expected {
				t.Errorf("FormatSize(%f) = %q, want %q", tt.size, got, tt.expected)
			}
		})
	}
}

