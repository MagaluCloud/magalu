package alphasync

import (
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/Masterminds/semver/v3"
)

// Released binaries keep this name: the bucket must never be deleted, or
// whoever creates one with the same name would serve these CLIs.
const DefaultBaseURL = "https://mgc-cli-preview.br-se1.magaluobjects.com"

const (
	maxManifest = 1 << 20
	maxSpec     = 32 << 20
)

var (
	logger       = mgcLoggerPkg.NewLazy[Config]()
	errNoSource  = errors.New("this build of mgc has no alpha source")
	validVersion = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`)
)

type Config struct {
	BaseURL string
	// Dir holds the alpha modules installed on this machine for the
	// environment the CLI points at.
	Dir        string
	CLIVersion string
	// HTTPClient must not authenticate: the bucket is public, and a token
	// sent to it would leak.
	HTTPClient *http.Client
	Check      func(preview.Package) error
	Keys       Keys
	// Group builds the commands of a store, to learn the scopes they need.
	Group func(preview.Store) core.Grouper
}

// Manifest lists every version of every module given to a tenant.
type Manifest struct {
	Specs []ManifestSpec `json:"specs"`
}

type ManifestSpec struct {
	Module  string `json:"module"`
	Version string `json:"version"`
	// RFC 3339 in UTC, so that comparing the text compares the dates.
	PublishedAt   string `json:"published_at"`
	Env           string `json:"env"`
	Path          string `json:"path"`
	SHA256        string `json:"sha256"`
	MinCLIVersion string `json:"min_cli_version,omitempty"`
	KeyID         string `json:"key_id"`
	Signature     string `json:"signature"`
}

func (m Manifest) Validate() error {
	seen := map[string]bool{}
	for _, spec := range m.Specs {
		if err := CheckModule(spec.Module); err != nil {
			return err
		}
		if err := CheckVersion(spec.Version); err != nil {
			return fmt.Errorf("alpha module %q: %w", spec.Module, err)
		}
		if at, err := time.Parse(time.RFC3339, spec.PublishedAt); err != nil || at.UTC().Format(time.RFC3339) != spec.PublishedAt {
			return fmt.Errorf("alpha module %s has no publication date in UTC", spec.label())
		}
		if seen[spec.label()] {
			return fmt.Errorf("alpha module %s is listed twice", spec.label())
		}
		seen[spec.label()] = true
	}
	return nil
}

// Find gives a version of a module, or its newest one when version is empty.
func (m Manifest) Find(module, version string) (ManifestSpec, bool) {
	var found ManifestSpec
	ok := false
	for _, spec := range m.Specs {
		if spec.Module == module && (version == "" || spec.Version == version) && (!ok || spec.newerThan(found)) {
			found, ok = spec, true
		}
	}
	return found, ok
}

// Latest gives the newest version of every module, sorted by module.
func (m Manifest) Latest() []ManifestSpec {
	newest := map[string]ManifestSpec{}
	for _, spec := range m.Specs {
		if current, ok := newest[spec.Module]; !ok || spec.newerThan(current) {
			newest[spec.Module] = spec
		}
	}
	latest := slices.Collect(maps.Values(newest))
	slices.SortFunc(latest, func(a, b ManifestSpec) int { return strings.Compare(a.Module, b.Module) })
	return latest
}

// Versions gives the versions of a module, newest first.
func (m Manifest) Versions(module string) []string {
	var specs []ManifestSpec
	for _, spec := range m.Specs {
		if spec.Module == module {
			specs = append(specs, spec)
		}
	}
	slices.SortFunc(specs, func(a, b ManifestSpec) int {
		return cmp.Or(strings.Compare(b.PublishedAt, a.PublishedAt), strings.Compare(b.Version, a.Version))
	})
	versions := make([]string, len(specs))
	for i, spec := range specs {
		versions[i] = spec.Version
	}
	return versions
}

func (s ManifestSpec) newerThan(other ManifestSpec) bool {
	if s.PublishedAt != other.PublishedAt {
		return s.PublishedAt > other.PublishedAt
	}
	return s.Version > other.Version
}

func (s ManifestSpec) label() string {
	return s.Module + "@" + s.Version
}

// A module is named on the command line as module@version, so '@' never
// appears in its name.
func CheckModule(module string) error {
	if module == "" || module != filepath.Base(module) || strings.HasPrefix(module, ".") || strings.Contains(module, "@") {
		return fmt.Errorf("alpha module %q is not a valid name", module)
	}
	return nil
}

func CheckVersion(version string) error {
	if !validVersion.MatchString(version) {
		return fmt.Errorf("version %q is not valid: use up to 64 letters, digits, '.', '-' and '_'", version)
	}
	return nil
}

// fetchManifest gives an empty manifest to a tenant never invited.
func fetchManifest(ctx context.Context, cfg Config, path string) (*Manifest, error) {
	body, status, err := get(ctx, cfg, path, maxManifest)
	if err != nil {
		return nil, err
	}
	// The bucket answers 404 for a missing manifest. Any other answer, a 403
	// from a proxy included, may hide one that exists.
	switch status {
	case http.StatusOK:
	case http.StatusNotFound:
		return &Manifest{}, nil
	default:
		return nil, fmt.Errorf("alpha source answered %d for the manifest of your tenant", status)
	}

	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("manifest of your tenant does not parse: %w", err)
	}
	if err := manifest.Validate(); err != nil {
		return nil, fmt.Errorf("manifest of your tenant is not valid: %w", err)
	}
	return &manifest, nil
}

// manifestCache fetches each manifest once, however many modules come from it.
func manifestCache(ctx context.Context, cfg Config) func(path string) (*Manifest, error) {
	fetched := map[string]*Manifest{}
	return func(path string) (*Manifest, error) {
		if manifest, ok := fetched[path]; ok {
			return manifest, nil
		}
		manifest, err := fetchManifest(ctx, cfg, path)
		if err != nil {
			return nil, err
		}
		fetched[path] = manifest
		return manifest, nil
	}
}

// download gives the spec of a signed item once it is safe to install.
func download(ctx context.Context, cfg Config, spec ManifestSpec, env string) ([]byte, error) {
	if err := cfg.Keys.Verify(spec); err != nil {
		return nil, err
	}
	if spec.Env != env {
		return nil, fmt.Errorf("alpha module %s was published for %q, not for %q", spec.label(), spec.Env, env)
	}
	if spec.SHA256 == "" {
		return nil, fmt.Errorf("alpha module %s has no checksum in its manifest", spec.label())
	}

	body, status, err := get(ctx, cfg, spec.Path, maxSpec)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("alpha source answered %d for module %s", status, spec.label())
	}
	sum := sha256.Sum256(body)
	if !strings.EqualFold(hex.EncodeToString(sum[:]), spec.SHA256) {
		return nil, fmt.Errorf("alpha module %s does not match its checksum", spec.label())
	}
	if err := cfg.Check(preview.Package{Name: spec.Module, Spec: body}); err != nil {
		return nil, err
	}
	return body, nil
}

func get(ctx context.Context, cfg Config, path string, limit int64) ([]byte, int, error) {
	url := strings.TrimSuffix(cfg.BaseURL, "/") + "/" + strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	client := cfg.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, 0, err
	}
	if int64(len(body)) > limit {
		return nil, 0, fmt.Errorf("%s is larger than %d bytes", url, limit)
	}
	return body, resp.StatusCode, nil
}

func tooOld(cliVersion, minVersion string) string {
	if minVersion == "" {
		return ""
	}
	least, err := semver.NewVersion(minVersion)
	if err != nil {
		return fmt.Sprintf("its minimum mgc version %q is not valid", minVersion)
	}
	current, err := semver.NewVersion(cliVersion)
	if err != nil || !current.LessThan(least) {
		return ""
	}
	return fmt.Sprintf("needs mgc %s or newer", least)
}
