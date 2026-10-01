package fsstore

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}

func TestNames(t *testing.T) {
	t.Run("missing dir has no names", func(t *testing.T) {
		names, err := NewDir(filepath.Join(t.TempDir(), "nope")).Names()
		require.NoError(t, err)
		assert.Empty(t, names)
	})

	t.Run("lists only spec files, sorted", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "zeta.openapi.yaml", "z")
		writeFile(t, dir, "alpha.openapi.yaml", "a")
		writeFile(t, dir, ".openapi.yaml", "no name")
		writeFile(t, dir, "notes.txt", "hi")
		require.NoError(t, os.Mkdir(filepath.Join(dir, "dir.openapi.yaml"), 0o700))

		names, err := NewDir(dir).Names()
		require.NoError(t, err)
		assert.Equal(t, []string{"alpha", "zeta"}, names)
	})

	t.Run("an unreadable dir is an error", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("root reads any dir")
		}
		dir := t.TempDir()
		require.NoError(t, os.Chmod(dir, 0))
		t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

		_, err := NewDir(dir).Names()
		assert.Error(t, err)
	})
}

func TestRead(t *testing.T) {
	t.Run("reads the spec of a name", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "foo.openapi.yaml", "spec")

		pkg, err := NewDir(dir).Read("foo")
		require.NoError(t, err)
		assert.Equal(t, "foo", pkg.Name)
		assert.Equal(t, "spec", string(pkg.Spec))
	})

	t.Run("a missing name is an error", func(t *testing.T) {
		_, err := NewDir(t.TempDir()).Read("foo")
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("refuses names outside the dir", func(t *testing.T) {
		parent := t.TempDir()
		dir := filepath.Join(parent, "beta")
		require.NoError(t, os.Mkdir(dir, 0o700))
		writeFile(t, parent, "secret.openapi.yaml", "secret")

		for _, name := range []string{"", ".", "..", "../secret", "/etc/passwd", "sub/foo", `sub\foo`} {
			_, err := NewDir(dir).Read(name)
			assert.Error(t, err, name)
		}
	})
}

func TestString(t *testing.T) {
	assert.Equal(t, "/some/dir", NewDir("/some/dir").String())
}

// The beta specs ship embedded in the binary: the same store reads them from an
// fs.FS, and files that are not specs (a README) are not modules.
func TestEmbedded(t *testing.T) {
	store := New(fstest.MapFS{
		"README.md":              {Data: []byte("beta specs")},
		"foo.openapi.yaml":       {Data: []byte("spec")},
		"sub/bar.openapi.yaml":   {Data: []byte("nested")},
		"dir.openapi.yaml/x.txt": {Data: []byte("a dir named like a spec")},
	}, "this mgc release")

	names, err := store.Names()
	require.NoError(t, err)
	assert.Equal(t, []string{"foo"}, names)

	pkg, err := store.Read("foo")
	require.NoError(t, err)
	assert.Equal(t, "spec", string(pkg.Spec))

	for _, name := range []string{"sub/bar", "../foo", ""} {
		_, err := store.Read(name)
		assert.Error(t, err, name)
	}
	assert.Equal(t, "this mgc release", store.String())
}
