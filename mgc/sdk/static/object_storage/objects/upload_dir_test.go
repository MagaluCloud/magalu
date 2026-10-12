package objects

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

func TestWalkDirFilters(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.log", "sub/c.txt", "sub/d.log"} {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name     string
		filters  []common.FilterParams
		shallow  bool
		expected []string
	}{
		{
			name:     "no filters",
			expected: []string{"a.txt", "b.log", "sub/c.txt", "sub/d.log"},
		},
		{
			name:     "exclude",
			filters:  []common.FilterParams{{Exclude: "*.log"}},
			expected: []string{"a.txt", "sub/c.txt"},
		},
		{
			name:     "exclude then include",
			filters:  []common.FilterParams{{Exclude: "*"}, {Include: "*.txt"}},
			expected: []string{"a.txt", "sub/c.txt"},
		},
		{
			name:     "exclude with shallow",
			filters:  []common.FilterParams{{Exclude: "*.txt"}},
			shallow:  true,
			expected: []string{"b.log"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancelCause(context.Background())
			defer cancel(nil)

			files, err := walkDir(ctx, root, tt.shallow, common.NewFilterRule(tt.filters, cancel))
			if err != nil {
				t.Fatal(err)
			}

			var got []string
			for _, f := range files {
				rel, _ := filepath.Rel(root, f)
				got = append(got, filepath.ToSlash(rel))
			}
			slices.Sort(got)

			if !slices.Equal(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
