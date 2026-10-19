package sdk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview/alphasync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreviewCommands(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("MGC_CONFIG_DIR", dir)

	child := func(name string) core.Descriptor {
		t.Helper()
		desc, err := NewSdk().Group().GetChildByName(name)
		require.NoError(t, err, name)
		return desc
	}
	assert.False(t, child("preview").IsInternal())
	assert.True(t, child("alpha").IsInternal(), "hidden while nothing is installed")
	_, err := NewSdk().Group().GetChildByName("beta")
	assert.Error(t, err, "nothing is in beta")
	assert.NotEqual(t, preview.GroupID, child("alpha").GroupID(), "the help has no Preview section")

	installed := filepath.Join(dir, ".preview", "alpha", "prod")
	require.NoError(t, os.MkdirAll(installed, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(installed, "index.json"), []byte(`{"foo":{"module":"foo","manifest":"m"}}`), 0o600))
	alpha := child("alpha")
	assert.False(t, alpha.IsInternal(), "an installed module shows it")
	assert.Equal(t, preview.GroupID, alpha.GroupID(), "in the Preview section of the help")
	_, err = alpha.(core.Grouper).GetChildByName("foo")
	assert.NoError(t, err, "a module that fails its checks stays, to tell why")

	t.Setenv("MGC_ENV", "pre-prod")
	assert.True(t, child("alpha").IsInternal(), "a module of prod never shows in pre-prod")
}

// The notice has 3 s: retrying a source that fails would outlast them.
func TestCheckAlphaUpdatesTriesOnce(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("MGC_CONFIG_DIR", dir)
	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)
	t.Setenv("MGC_ALPHA_URL", server.URL)

	installed := filepath.Join(dir, ".preview", "alpha", "prod")
	require.NoError(t, os.MkdirAll(installed, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(installed, "index.json"), []byte(`{"foo":{"module":"foo","manifest":"m"}}`), 0o600))

	NewSdk().CheckAlphaUpdates(context.Background())
	assert.EqualValues(t, 1, hits.Load())
}

// With an API key, 'mgc preview install' finds the tenant through this command:
// renaming it in the network spec would break that.
func TestTheSecurityGroupsListExists(t *testing.T) {
	t.Setenv("MGC_CONFIG_DIR", t.TempDir())

	var node core.Descriptor = NewSdk().Group()
	for _, name := range alphasync.SecurityGroupsList {
		group, ok := node.(core.Grouper)
		require.True(t, ok, "%q has no subcommands", node.Name())
		child, err := group.GetChildByName(name)
		require.NoError(t, err, name)
		node = child
	}
	_, ok := node.(core.Executor)
	assert.True(t, ok)
}
