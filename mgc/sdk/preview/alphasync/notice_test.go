package alphasync

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkAt(t *testing.T, cfg Config, now time.Time) (string, error) {
	t.Helper()
	var out bytes.Buffer
	err := checkUpdates(context.Background(), cfg, now, &out)
	return out.String(), err
}

func TestCheckUpdatesAtMostOnceADay(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"), b.spec("iam", "2", "iam spec v2"))

	now := time.Now()
	out, err := checkAt(t, cfg, now)
	require.NoError(t, err)
	assert.Equal(t, "Alpha updates: iam 1 -> 2. Run 'mgc preview upgrade'.\n", out)
	assert.Equal(t, "iam spec v1", specOf(t, cfg.Dir, "iam"), "nothing changes on its own")

	b.hits = nil
	out, err = checkAt(t, cfg, now.Add(23*time.Hour))
	require.NoError(t, err)
	assert.Empty(t, out)
	assert.Empty(t, b.hits)

	out, err = checkAt(t, cfg, now.Add(25*time.Hour))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestCheckUpdatesTellsWhatWasTakenAway(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"), b.spec("dns", "1", "dns spec"), b.spec("connect", "1", "connect spec"))
	_, err := installFor(t, cfg, tenantA, "iam", "dns", "connect@1")
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))
	out, err := checkAt(t, cfg, time.Now())
	require.NoError(t, err)
	assert.Equal(t, "Alpha updates: connect was taken away, dns was taken away. Run 'mgc preview upgrade'.\n", out,
		"a module pinned to a version is told about only when taken away")
}

func TestCheckUpdatesTellsWhenAPinnedVersionIsTakenAway(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"), b.spec("iam", "2", "iam spec v2"))
	_, err := installFor(t, cfg, tenantA, "iam@1")
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, b.spec("iam", "2", "iam spec v2"))
	out, err := checkAt(t, cfg, time.Now())
	require.NoError(t, err)
	assert.Equal(t, "Alpha updates: iam@1 was taken away. Run 'mgc preview upgrade'.\n", out)
}

func TestCheckUpdatesTellsWhenMgcMustBeUpdatedFirst(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"), b.spec("dns", "1", "dns spec v1"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, b.spec("iam", "2", "iam spec v2", minVersion("0.70.0")), b.spec("dns", "2", "dns spec v2"))
	out, err := checkAt(t, cfg, time.Now())
	require.NoError(t, err)
	assert.Equal(t, ""+
		"Alpha updates: dns 1 -> 2. Run 'mgc preview upgrade'.\n"+
		"Alpha updates that need a newer mgc: iam 1 -> 2. Update mgc first.\n", out)
}

func TestCheckUpdatesKeepsQuietWhenItCannotTell(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)

	b.status[manifestPath(t, EnvProd, tenantA)] = http.StatusInternalServerError
	now := time.Now()
	out, err := checkAt(t, cfg, now)
	assert.Error(t, err)
	assert.Empty(t, out)

	b.hits = nil
	_, err = checkAt(t, cfg, now.Add(time.Hour))
	require.NoError(t, err)
	assert.Empty(t, b.hits, "a source failing is not tried again on every command")
}

func TestCheckUpdatesWithNothingInstalled(t *testing.T) {
	b, cfg := newBucket(t)
	out, err := checkAt(t, cfg, time.Now())
	require.NoError(t, err)
	assert.Empty(t, out)
	assert.Empty(t, b.hits)

	_, err = os.Stat(filepath.Join(cfg.Dir, stateFile))
	assert.ErrorIs(t, err, os.ErrNotExist)
}
