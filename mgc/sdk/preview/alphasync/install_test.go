package alphasync

import (
	"context"
	"crypto/ed25519"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func installFor(t *testing.T, cfg Config, tenant string, args ...string) (*Result, error) {
	t.Helper()
	var targets []Target
	for _, arg := range args {
		target, err := ParseTarget(arg)
		require.NoError(t, err)
		targets = append(targets, target)
	}
	return Install(context.Background(), cfg, manifestPath(t, EnvProd, tenant), EnvProd, targets)
}

func indexOf(t *testing.T, dir string) map[string]Installed {
	t.Helper()
	index, err := readIndex(dir)
	require.NoError(t, err)
	return index
}

// specOf gives the spec the index vouches for, or nothing when not installed.
func specOf(t *testing.T, dir, module string) string {
	t.Helper()
	entry, ok := indexOf(t, dir)[module]
	if !ok {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(dir, specFile(entry)))
	require.NoError(t, err)
	return string(data)
}

func specPath(t *testing.T, dir, module string) string {
	t.Helper()
	entry, ok := indexOf(t, dir)[module]
	require.True(t, ok, module)
	return filepath.Join(dir, specFile(entry))
}

// assertOnlyIndexedSpecs checks that every spec on disk is one the index names.
func assertOnlyIndexedSpecs(t *testing.T, dir string) {
	t.Helper()
	var named, found []string
	for _, entry := range indexOf(t, dir) {
		named = append(named, specFile(entry))
	}
	for name := range installed(t, dir) {
		if strings.HasSuffix(name, specSuffix) {
			found = append(found, name)
		}
	}
	assert.ElementsMatch(t, named, found, "no spec is left behind")
}

func TestParseTarget(t *testing.T) {
	target, err := ParseTarget("iam@2026-10-08")
	require.NoError(t, err)
	assert.Equal(t, Target{Module: "iam", Version: "2026-10-08"}, target)

	target, err = ParseTarget("iam")
	require.NoError(t, err)
	assert.Equal(t, Target{Module: "iam"}, target)

	for _, arg := range []string{"", "@1", "../iam", "iam@", "iam@1/2"} {
		_, err := ParseTarget(arg)
		assert.Error(t, err, arg)
	}
	assert.Error(t, CheckModule("a@b"), "a module named a@b could never be installed by name")
}

func TestInstallEverythingTheTenantWasGiven(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"), b.spec("iam", "2", "iam spec v2"), b.spec("dns", "1", "dns spec"))

	result, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	assert.Equal(t, []string{"dns@1", "iam@2"}, result.Installed)

	assert.Equal(t, "iam spec v2", specOf(t, cfg.Dir, "iam"))
	assert.Equal(t, "dns spec", specOf(t, cfg.Dir, "dns"))

	entry := indexOf(t, cfg.Dir)["iam"]
	assert.Equal(t, "2", entry.Version)
	assert.Equal(t, manifestPath(t, EnvProd, tenantA), entry.Manifest)
	assert.False(t, entry.Pinned)

	info, err := os.Stat(specPath(t, cfg.Dir, "iam"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	info, err = os.Stat(cfg.Dir)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), info.Mode().Perm())
}

func TestInstallPinsANamedVersion(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"), b.spec("iam", "2", "iam spec v2"))

	result, err := installFor(t, cfg, tenantA, "iam@1")
	require.NoError(t, err)
	assert.Equal(t, []string{"iam@1"}, result.Installed)
	assert.Equal(t, "iam spec v1", specOf(t, cfg.Dir, "iam"))
	assert.True(t, indexOf(t, cfg.Dir)["iam"].Pinned)

	_, err = installFor(t, cfg, tenantA, "iam")
	require.NoError(t, err)
	assert.Equal(t, "iam spec v2", specOf(t, cfg.Dir, "iam"))
	assert.False(t, indexOf(t, cfg.Dir)["iam"].Pinned, "no version means the newest, not pinned")
	assertOnlyIndexedSpecs(t, cfg.Dir)
}

func TestInstallTakesEachModuleOnce(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"), b.spec("iam", "2", "iam spec v2"))

	result, err := installFor(t, cfg, tenantA, "iam", "iam")
	require.NoError(t, err)
	assert.Equal(t, []string{"iam@2"}, result.Installed)

	_, err = installFor(t, cfg, tenantA, "iam@1", "iam")
	assert.EqualError(t, err, `alpha module "iam" is named more than once`)
	assert.Equal(t, "iam spec v2", specOf(t, cfg.Dir, "iam"))
}

func TestInstallAnswersTheSameForWhatIsNotGiven(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))

	_, err := installFor(t, cfg, tenantA, "dns")
	assert.EqualError(t, err, `alpha module "dns" is not available for your tenant`)
	_, err = installFor(t, cfg, tenantB, "iam")
	assert.EqualError(t, err, `alpha module "iam" is not available for your tenant`, "a tenant never invited is told the same")
	_, err = installFor(t, cfg, tenantA, "iam@2")
	assert.EqualError(t, err, `alpha module iam@2 is not available for your tenant`)

	result, err := installFor(t, cfg, tenantB)
	require.NoError(t, err)
	assert.Empty(t, result.Installed)
	_, err = os.Stat(cfg.Dir)
	assert.ErrorIs(t, err, os.ErrNotExist, "a tenant never invited leaves no trace")
}

func TestInstallKeepsWhatIsInstalledOnFailure(t *testing.T) {
	cases := map[string]func(t *testing.T, b *bucket, cfg *Config){
		"checksum mismatch": func(t *testing.T, b *bucket, cfg *Config) {
			b.publish(t, EnvProd, tenantA, b.spec("iam", "2", "iam spec v2", func(s *ManifestSpec) { s.SHA256 = checksum("something else") }))
		},
		"not signed": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("iam", "2", "iam spec v2")
			spec.KeyID, spec.Signature = "", ""
			b.publish(t, EnvProd, tenantA, spec)
		},
		"signed fields changed": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("iam", "2", "iam spec v2", minVersion("0.70.0"))
			spec.MinCLIVersion = ""
			b.publish(t, EnvProd, tenantA, spec)
		},
		"published for another environment": func(t *testing.T, b *bucket, cfg *Config) {
			b.publish(t, EnvProd, tenantA, b.spec("iam", "2", "iam spec v2", forEnv(EnvPreProd)))
		},
		"environment changed after signing": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("iam", "2", "iam spec v2", forEnv(EnvPreProd))
			spec.Env = EnvProd
			b.publish(t, EnvProd, tenantA, spec)
		},
		"module renamed": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("dns", "2", "dns spec")
			spec.Module = "iam"
			b.publish(t, EnvProd, tenantA, spec)
		},
		"signed with an unknown key": func(t *testing.T, b *bucket, cfg *Config) {
			_, other, _ := ed25519.GenerateKey(nil)
			b.publish(t, EnvProd, tenantA, Sign(other, "other", b.spec("iam", "2", "iam spec v2")))
		},
		"signed with another key under a trusted id": func(t *testing.T, b *bucket, cfg *Config) {
			_, other, _ := ed25519.GenerateKey(nil)
			b.publish(t, EnvProd, tenantA, Sign(other, testKeyID, b.spec("iam", "2", "iam spec v2")))
		},
		"spec does not load": func(t *testing.T, b *bucket, cfg *Config) {
			b.publish(t, EnvProd, tenantA, b.spec("iam", "2", "iam spec v2"))
			cfg.Check = func(preview.Package) error { return errors.New("spec has no operations") }
		},
		"spec missing in the bucket": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("iam", "2", "iam spec v2")
			delete(b.files, spec.Path)
			b.publish(t, EnvProd, tenantA, spec)
		},
		"version listed twice": func(t *testing.T, b *bucket, cfg *Config) {
			spec := b.spec("iam", "2", "iam spec v2")
			b.publish(t, EnvProd, tenantA, spec, spec)
		},
		"source refuses": func(t *testing.T, b *bucket, cfg *Config) {
			b.status[manifestPath(t, EnvProd, tenantA)] = http.StatusForbidden
		},
		"source fails": func(t *testing.T, b *bucket, cfg *Config) {
			b.status[manifestPath(t, EnvProd, tenantA)] = http.StatusInternalServerError
		},
		"manifest does not parse": func(t *testing.T, b *bucket, cfg *Config) {
			b.files[manifestPath(t, EnvProd, tenantA)] = "not json"
		},
	}

	for name, breakIt := range cases {
		t.Run(name, func(t *testing.T) {
			b, cfg := newBucket(t)
			b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))
			_, err := installFor(t, cfg, tenantA, "iam")
			require.NoError(t, err)
			before := installed(t, cfg.Dir)

			breakIt(t, b, &cfg)
			_, err = installFor(t, cfg, tenantA, "iam")
			require.Error(t, err)
			assert.NotContains(t, err.Error(), tenantA)
			assert.Equal(t, before, installed(t, cfg.Dir))
		})
	}
}

func TestInstallLeavesOutSpecsForNewerCLIs(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("dns", "1", "dns spec", minVersion("0.70.0")), b.spec("iam", "1", "iam spec", minVersion("0.66.0")))

	result, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	assert.Equal(t, []string{"iam@1"}, result.Installed)
	assert.Equal(t, []Skipped{{Module: "dns", Reason: "needs mgc 0.70.0 or newer"}}, result.Skipped)

	_, err = installFor(t, cfg, tenantA, "dns")
	assert.ErrorContains(t, err, "needs mgc 0.70.0 or newer", "a module named is not skipped in silence")

	cfg.CLIVersion = "dev"
	result, err = installFor(t, cfg, tenantA)
	require.NoError(t, err)
	assert.Equal(t, []string{"dns@1"}, result.Installed, "a build without a version takes everything")
}

func TestInstallWithoutNamesLeavesWhatIsInstalled(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"), b.spec("dns", "1", "dns spec v1"))
	_, err := installFor(t, cfg, tenantA, "iam", "dns@1")
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, b.spec("iam", "2", "iam spec v2"), b.spec("dns", "2", "dns spec v2"), b.spec("vpn", "1", "vpn spec"))
	result, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	assert.Equal(t, []string{"vpn@1"}, result.Installed)

	assert.Equal(t, "iam spec v1", specOf(t, cfg.Dir, "iam"), "only an upgrade moves what is installed")
	assert.Equal(t, "dns spec v1", specOf(t, cfg.Dir, "dns"))
	assert.True(t, indexOf(t, cfg.Dir)["dns"].Pinned, "a pin survives")
}

func TestInstallChangesNothingWhenOneModuleFails(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("dns", "1", "dns spec"), b.spec("iam", "1", "iam spec", func(s *ManifestSpec) { s.SHA256 = checksum("other") }))

	_, err := installFor(t, cfg, tenantA)
	assert.ErrorContains(t, err, "does not match its checksum")
	_, err = os.Stat(cfg.Dir)
	assert.ErrorIs(t, err, os.ErrNotExist, "dns came first, and is not installed either")
}

func TestInstallAgainOnlyDownloadsWhatChanged(t *testing.T) {
	b, cfg := newBucket(t)
	iam := b.spec("iam", "1", "iam spec")
	b.publish(t, EnvProd, tenantA, iam)
	_, err := installFor(t, cfg, tenantA, "iam")
	require.NoError(t, err)

	b.hits = nil
	_, err = installFor(t, cfg, tenantA, "iam")
	require.NoError(t, err)
	assert.NotContains(t, b.hits, iam.Path, "the same spec is not downloaded again")

	require.NoError(t, os.WriteFile(specPath(t, cfg.Dir, "iam"), []byte("edited by hand"), 0o600))
	_, err = installFor(t, cfg, tenantA, "iam")
	require.NoError(t, err)
	assert.Contains(t, b.hits, iam.Path)
	assert.Equal(t, "iam spec", specOf(t, cfg.Dir, "iam"), "installing again repairs it")
	assertOnlyIndexedSpecs(t, cfg.Dir)
}

func TestUpgradeLeavesPinnedModulesAlone(t *testing.T) {
	b, cfg := newBucket(t)
	iam1, dns1 := b.spec("iam", "1", "iam spec v1"), b.spec("dns", "1", "dns spec v1")
	b.publish(t, EnvProd, tenantA, iam1, dns1)
	_, err := installFor(t, cfg, tenantA, "iam", "dns@1")
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, iam1, dns1, b.spec("iam", "2", "iam spec v2"), b.spec("dns", "2", "dns spec v2"))
	result, err := Upgrade(context.Background(), cfg, nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"iam@2"}, result.Installed)
	assert.Equal(t, "dns spec v1", specOf(t, cfg.Dir, "dns"))

	result, err = Upgrade(context.Background(), cfg, []string{"dns"})
	require.NoError(t, err)
	assert.Equal(t, []string{"dns@2"}, result.Installed)
	assert.False(t, indexOf(t, cfg.Dir)["dns"].Pinned, "a module upgraded by name is no longer pinned")

	result, err = Upgrade(context.Background(), cfg, nil)
	require.NoError(t, err)
	assert.Empty(t, result.Installed)
}

func TestUpgradeRemovesWhatWasTakenAway(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"), b.spec("dns", "1", "dns spec"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))
	result, err := Upgrade(context.Background(), cfg, nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"dns"}, result.Removed)
	assert.NotContains(t, indexOf(t, cfg.Dir), "dns")
	assertOnlyIndexedSpecs(t, cfg.Dir)

	delete(b.files, manifestPath(t, EnvProd, tenantA))
	result, err = Upgrade(context.Background(), cfg, nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"iam"}, result.Removed, "a manifest removed takes everything it gave")
}

func TestUpgradeRemovesAPinnedModuleTakenAway(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"), b.spec("dns", "1", "dns spec"))
	_, err := installFor(t, cfg, tenantA, "iam@1", "dns@1")
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, b.spec("dns", "1", "dns spec"), b.spec("dns", "2", "dns spec v2"))
	result, err := Upgrade(context.Background(), cfg, nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"iam"}, result.Removed)
	assert.Empty(t, result.Installed, "dns keeps its pin")
}

func TestUpgradeRemovesAPinnedVersionTakenAway(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec v1"), b.spec("dns", "1", "dns spec v1"))
	_, err := installFor(t, cfg, tenantA, "iam@1", "dns@1")
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, b.spec("iam", "2", "iam spec v2"), b.spec("dns", "2", "dns spec v2"))
	result, err := Upgrade(context.Background(), cfg, []string{"dns"})
	require.NoError(t, err)
	assert.Equal(t, []string{"dns@2"}, result.Installed, "named, it moves to the newest")

	result, err = Upgrade(context.Background(), cfg, nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"iam"}, result.Removed, "the version it is pinned to was taken away")
	assert.Empty(t, result.Installed)
	assertOnlyIndexedSpecs(t, cfg.Dir)
}

func TestUpgradeChangesNothingWhenOneModuleFails(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("dns", "1", "dns spec"))
	b.publish(t, EnvProd, tenantB, b.spec("iam", "1", "iam spec"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	_, err = installFor(t, cfg, tenantB)
	require.NoError(t, err)
	before := installed(t, cfg.Dir)

	b.publish(t, EnvProd, tenantA)
	b.status[manifestPath(t, EnvProd, tenantB)] = http.StatusInternalServerError
	_, err = Upgrade(context.Background(), cfg, nil)
	assert.Error(t, err)
	assert.Equal(t, before, installed(t, cfg.Dir), "dns was taken away, but is removed only when everything else works")
}

func TestUpgradeChangesNothingWhenAWriteFails(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("dns", "1", "dns spec v1"), b.spec("iam", "1", "iam spec v1"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	before := installed(t, cfg.Dir)

	iam2 := b.spec("iam", "2", "iam spec v2")
	b.publish(t, EnvProd, tenantA, b.spec("dns", "2", "dns spec v2"), iam2)
	// dns is written first; a directory where the spec of iam goes fails it.
	blocked := filepath.Join(cfg.Dir, specFile(Installed{ManifestSpec: iam2}))
	require.NoError(t, os.Mkdir(blocked, 0o700))
	_, err = Upgrade(context.Background(), cfg, nil)
	assert.Error(t, err)

	require.NoError(t, os.Remove(blocked))
	assert.Equal(t, before, installed(t, cfg.Dir), "the spec of dns written before is removed")
}

func TestUpgradeTakesEachNameOnce(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA)
	result, err := Upgrade(context.Background(), cfg, []string{"iam", "iam"})
	require.NoError(t, err)
	assert.Equal(t, []string{"iam"}, result.Removed)
}

func TestUpgradeReadsTheManifestOfTheInstall(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))
	b.publish(t, EnvProd, tenantB, b.spec("dns", "1", "dns spec"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	_, err = installFor(t, cfg, tenantB)
	require.NoError(t, err)

	b.hits = nil
	_, err = Upgrade(context.Background(), cfg, nil)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{manifestPath(t, EnvProd, tenantA), manifestPath(t, EnvProd, tenantB)}, b.hits,
		"each module goes back to the tenant that installed it, and each manifest is read once")
}

func TestUpgradeKeepsEverythingWhenTheSourceFails(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	before := installed(t, cfg.Dir)

	b.status[manifestPath(t, EnvProd, tenantA)] = http.StatusForbidden
	_, err = Upgrade(context.Background(), cfg, nil)
	assert.Error(t, err)
	assert.Equal(t, before, installed(t, cfg.Dir))

	_, err = Upgrade(context.Background(), cfg, []string{"dns"})
	assert.EqualError(t, err, `alpha module "dns" is not installed`)
}

func TestUninstall(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"), b.spec("dns", "1", "dns spec"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)

	require.NoError(t, Uninstall(cfg, "iam"))
	assert.Equal(t, []string{"dns"}, keysOf(indexOf(t, cfg.Dir)))
	assertOnlyIndexedSpecs(t, cfg.Dir)

	assert.EqualError(t, Uninstall(cfg, "iam"), `alpha module "iam" is not installed`)
}

func TestListTellsTheStatusOfEachModule(t *testing.T) {
	b, cfg := newBucket(t)
	iam1, iam2 := b.spec("iam", "1", "iam spec v1"), b.spec("iam", "2", "iam spec v2")
	dns, connect1, foo := b.spec("dns", "1", "dns spec"), b.spec("connect", "1", "connect spec v1"), b.spec("foo", "1", "foo spec")
	kms1, lb1 := b.spec("kms", "1", "kms spec v1"), b.spec("lb", "1", "lb spec v1")
	b.publish(t, EnvProd, tenantA, iam1, iam2, dns, connect1, foo, kms1, lb1)
	b.publish(t, EnvProd, tenantB, b.spec("bar", "1", "bar spec"))
	_, err := installFor(t, cfg, tenantA, "iam@1", "dns", "connect", "foo", "kms@1", "lb")
	require.NoError(t, err)
	_, err = installFor(t, cfg, tenantB)
	require.NoError(t, err)

	b.publish(t, EnvProd, tenantA, iam1, iam2, dns, connect1, b.spec("connect", "2", "connect spec v2"), b.spec("vpn", "1", "vpn spec"),
		b.spec("kms", "2", "kms spec v2"), lb1, b.spec("lb", "2", "lb spec v2", minVersion("0.70.0")),
		b.spec("zone", "1", "zone spec", minVersion("0.70.0")), b.spec("bar", "2", "bar spec v2"))
	b.hits = nil
	listed, err := List(context.Background(), cfg, manifestPath(t, EnvProd, tenantA))
	require.NoError(t, err)
	assert.Equal(t, []Listed{
		{Module: "bar", Status: "other tenant", Installed: "1", Latest: "2", Versions: []string{"2"}},
		{Module: "connect", Status: "update available", Installed: "1", Latest: "2", Versions: []string{"2", "1"}},
		{Module: "dns", Status: "up to date", Installed: "1", Latest: "1", Versions: []string{"1"}},
		{Module: "foo", Status: "taken away", Installed: "1"},
		{Module: "iam", Status: "pinned", Installed: "1", Latest: "2", Versions: []string{"2", "1"}},
		{Module: "kms", Status: "taken away", Installed: "1", Latest: "2", Versions: []string{"2"}},
		{Module: "lb", Status: "needs newer mgc", Installed: "1", Latest: "2", Versions: []string{"2", "1"}},
		{Module: "vpn", Status: "not installed", Latest: "1", Versions: []string{"1"}},
		{Module: "zone", Status: "needs newer mgc", Latest: "1", Versions: []string{"1"}},
	}, listed)
	assert.Equal(t, []string{manifestPath(t, EnvProd, tenantA)}, b.hits, "bar came from another tenant, whose manifest is not read")
}

func TestInstalledStoreChecksEveryRead(t *testing.T) {
	b, cfg := newBucket(t)
	b.publish(t, EnvProd, tenantA, b.spec("iam", "1", "iam spec"))
	_, err := installFor(t, cfg, tenantA)
	require.NoError(t, err)
	store := InstalledStore(cfg.Dir, cfg.Keys)

	names, err := store.Names()
	require.NoError(t, err)
	assert.Equal(t, []string{"iam"}, names)
	pkg, err := store.Read("iam")
	require.NoError(t, err)
	assert.Equal(t, "iam spec", string(pkg.Spec))
	_, err = InstalledStore(cfg.Dir, Keys{}).Read("iam")
	assert.ErrorContains(t, err, "which this mgc does not know", "a key the CLI no longer trusts")

	require.NoError(t, os.WriteFile(specPath(t, cfg.Dir, "iam"), []byte("changed"), 0o600))
	_, err = store.Read("iam")
	assert.ErrorContains(t, err, "changed since it was installed")
	require.NoError(t, os.Remove(specPath(t, cfg.Dir, "iam")))
	_, err = store.Read("iam")
	assert.ErrorContains(t, err, "its spec is missing: run 'mgc preview install iam'")

	require.NoError(t, os.WriteFile(filepath.Join(cfg.Dir, "dns.openapi.yaml"), []byte("copied by hand"), 0o600))
	names, err = InstalledStore(cfg.Dir, cfg.Keys).Names()
	require.NoError(t, err)
	assert.Equal(t, []string{"iam"}, names, "only what the index lists")
	_, err = store.Read("dns")
	assert.EqualError(t, err, `alpha module "dns" is not installed`)
}

func keysOf(index map[string]Installed) []string {
	var keys []string
	for key := range index {
		keys = append(keys, key)
	}
	return keys
}
