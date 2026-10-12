package alphasync

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
)

// Released binaries keep this name: the bucket must never be deleted, or
// whoever creates one with the same name would serve these CLIs.
const DefaultBaseURL = "https://mgc-cli-preview-public.br-ne1.magaluobjects.com"

const (
	invitePrefix = "mgc-alpha:"
	inviteFile   = "invite.json"
	indexFile    = "index.json"
	specSuffix   = ".openapi.yaml"
	maxManifest  = 1 << 20
	maxSpec      = 32 << 20
)

var errNoSource = errors.New("this build of mgc has no alpha source")

type Manifest struct {
	Specs []ManifestSpec `json:"specs"`
}

type ManifestSpec struct {
	Module        string `json:"module"`
	Path          string `json:"path"`
	SHA256        string `json:"sha256"`
	MinCLIVersion string `json:"min_cli_version,omitempty"`
	KeyID         string `json:"key_id"`
	Signature     string `json:"signature"`
}

type Config struct {
	BaseURL    string
	Dir        string
	CLIVersion string
	// HTTPClient must not authenticate: the bucket is public, and a token
	// sent to it would leak.
	HTTPClient *http.Client
	Check      func(preview.Package) error
	Keys       Keys
}

type Skipped struct {
	Module string `json:"module"`
	Reason string `json:"reason"`
}

type Result struct {
	Installed []string  `json:"installed"`
	Skipped   []Skipped `json:"-"`
}

func NewInvite() string {
	return uuid.NewString()
}

func NormalizeInvite(invite string) (string, error) {
	id, err := uuid.Parse(strings.TrimSpace(invite))
	if err != nil {
		return "", fmt.Errorf("%q is not an alpha invite", invite)
	}
	return id.String(), nil
}

func ManifestPath(invite string) string {
	sum := sha256.Sum256([]byte(invitePrefix + invite))
	return "alpha/invites/" + hex.EncodeToString(sum[:]) + ".json"
}

func CheckModule(module string) error {
	if module == "" || module != filepath.Base(module) || strings.HasPrefix(module, ".") {
		return fmt.Errorf("alpha module %q is not a valid name", module)
	}
	return nil
}

func SavedInvite(dir string) (string, error) {
	file := filepath.Join(dir, inviteFile)
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	var invite string
	if err := json.Unmarshal(data, &invite); err != nil {
		return "", fmt.Errorf("%s is damaged, sync again with your invite: %w", file, err)
	}
	return invite, nil
}

func Sync(ctx context.Context, cfg Config, invite string) (*Result, error) {
	if cfg.BaseURL == "" {
		return nil, errNoSource
	}

	invite, err := NormalizeInvite(invite)
	if err != nil {
		return nil, err
	}

	manifest, _, err := fetchManifest(ctx, cfg, invite)
	if err != nil {
		return nil, err
	}

	wanted := map[string]ManifestSpec{}
	for _, spec := range manifest.Specs {
		if _, ok := wanted[spec.Module]; ok {
			return nil, fmt.Errorf("alpha module %q is listed twice", spec.Module)
		}
		wanted[spec.Module] = spec
	}

	result := &Result{Installed: []string{}}

	files := map[string][]byte{}
	index := map[string]ManifestSpec{}
	modules := make([]string, 0, len(wanted))
	for module := range wanted {
		modules = append(modules, module)
	}
	slices.Sort(modules)

	for _, module := range modules {
		spec := wanted[module]
		if err := CheckModule(module); err != nil {
			return nil, err
		}
		if err := cfg.Keys.Verify(spec); err != nil {
			return nil, err
		}
		if reason := tooOld(cfg.CLIVersion, spec.MinCLIVersion); reason != "" {
			result.Skipped = append(result.Skipped, Skipped{Module: module, Reason: reason})
			continue
		}

		data, err := fetchSpec(ctx, cfg, spec)
		if err != nil {
			return nil, err
		}
		if err := cfg.Check(preview.Package{Name: module, Spec: data}); err != nil {
			return nil, err
		}
		files[module+specSuffix] = data
		index[module] = spec
		result.Installed = append(result.Installed, module)
	}

	indexJSON, err := json.Marshal(index)
	if err != nil {
		return nil, err
	}
	files[indexFile] = indexJSON

	inviteJSON, err := json.Marshal(invite)
	if err != nil {
		return nil, err
	}
	files[inviteFile] = inviteJSON

	if err := replaceDir(cfg.Dir, files); err != nil {
		return nil, err
	}
	return result, nil
}

func fetchManifest(ctx context.Context, cfg Config, invite string) (*Manifest, bool, error) {
	body, status, err := get(ctx, cfg, ManifestPath(invite), maxManifest)
	if err != nil {
		return nil, false, err
	}
	switch status {
	case http.StatusOK:
	case http.StatusNotFound, http.StatusForbidden:
		return &Manifest{}, false, nil
	default:
		return nil, false, fmt.Errorf("alpha source answered %d for the manifest of your invite", status)
	}

	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, false, fmt.Errorf("manifest of your invite does not parse: %w", err)
	}
	return &manifest, true, nil
}

func fetchSpec(ctx context.Context, cfg Config, spec ManifestSpec) ([]byte, error) {
	if spec.SHA256 == "" {
		return nil, fmt.Errorf("alpha module %q has no checksum in its manifest", spec.Module)
	}

	body, status, err := get(ctx, cfg, spec.Path, maxSpec)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("alpha source answered %d for module %q", status, spec.Module)
	}

	sum := sha256.Sum256(body)
	if !strings.EqualFold(hex.EncodeToString(sum[:]), spec.SHA256) {
		return nil, fmt.Errorf("alpha module %q does not match its checksum", spec.Module)
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
	return fmt.Sprintf("needs mgc %s or newer: update mgc and sync again", least)
}

func replaceDir(dir string, files map[string][]byte) error {
	parent := filepath.Dir(dir)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}

	tmp, err := os.MkdirTemp(parent, "."+filepath.Base(dir)+"-sync-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	for name, data := range files {
		if err := os.WriteFile(filepath.Join(tmp, name), data, 0o600); err != nil {
			return err
		}
	}

	old := ""
	if _, err := os.Stat(dir); err == nil {
		old = tmp + "-old"
		if err := os.Rename(dir, old); err != nil {
			return err
		}
	}
	if err := os.Rename(tmp, dir); err != nil {
		if old != "" {
			_ = os.Rename(old, dir)
		}
		return err
	}
	if old != "" {
		return os.RemoveAll(old)
	}
	return nil
}
