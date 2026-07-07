package objects

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func setupWalkDirTestDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	createFile(t, root, "a.txt")
	makeSymlink(t, root, "broken", "/nonexistent")
	makeSymlink(t, root, "dir-link", "subdir")
	makeSymlink(t, root, "file-link", "a.txt")
	createFile(t, root, "subdir", "b.txt")

	return root
}

func createFile(t *testing.T, root string, parts ...string) {
	t.Helper()
	path := filepath.Join(append([]string{root}, parts...)...)
	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path, []byte("content"), 0o644)
	if err != nil {
		t.Fatal(err)
	}
}

func makeSymlink(t *testing.T, root string, link, target string) {
	t.Helper()
	linkPath := filepath.Join(root, link)
	err := os.Symlink(target, linkPath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWalkDir(t *testing.T) {
	root := setupWalkDirTestDir(t)

	tests := []struct {
		name    string
		shallow bool
		want    []string
	}{
		{
			name:    "recursive walk includes subdirectory files",
			shallow: false,
			want: []string{
				"a.txt",
				"file-link",
				"subdir/b.txt",
			},
		},
		{
			name:    "shallow walk skips subdirectory files",
			shallow: true,
			want: []string{
				"a.txt",
				"file-link",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := walkDir(context.Background(), root, tt.shallow)
			if err != nil {
				t.Fatalf("walkDir() error = %v", err)
			}

			gotRel := make([]string, len(got))
			for i, f := range got {
				rel, _ := filepath.Rel(root, f)
				gotRel[i] = rel
			}
			sort.Strings(gotRel)
			sort.Strings(tt.want)

			if len(gotRel) != len(tt.want) {
				t.Errorf("walkDir() got %d files, want %d\ngot:  %v\nwant: %v",
					len(gotRel), len(tt.want), gotRel, tt.want)
				return
			}
			for i := range gotRel {
				if gotRel[i] != tt.want[i] {
					t.Errorf("walkDir() file[%d] = %q, want %q\ngot:  %v\nwant: %v",
						i, gotRel[i], tt.want[i], gotRel, tt.want)
				}
			}
		})
	}
}

func TestWalkDir_SingleRegularFile(t *testing.T) {
	root := t.TempDir()
	createFile(t, root, "file.txt")

	got, err := walkDir(context.Background(), root, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(got), got)
	}
}

func TestWalkDir_AllBrokenSymlinks(t *testing.T) {
	root := t.TempDir()
	makeSymlink(t, root, "link1", "/nonexistent1")
	makeSymlink(t, root, "link2", "/nonexistent2")

	got, err := walkDir(context.Background(), root, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 files (all symlinks broken), got %d: %v", len(got), got)
	}
}

func TestWalkDir_NonExistentRoot(t *testing.T) {
	_, err := walkDir(context.Background(), "/tmp/mgc-test-nonexistent-12345", false)
	if err == nil {
		t.Error("expected error for non-existent root, got nil")
	}
}

func TestWalkDir_ContextCancelled(t *testing.T) {
	root := setupWalkDirTestDir(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := walkDir(ctx, root, false)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	}
}
