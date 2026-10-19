package main

import (
	"bufio"
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview/alphasync"
	"github.com/google/uuid"
)

const planFile = "plan.tsv"

type tenantFlags []string

func (t *tenantFlags) String() string     { return strings.Join(*t, ",") }
func (t *tenantFlags) Set(v string) error { *t = append(*t, v); return nil }

type publishing struct {
	keys    alphasync.Keys
	env     string
	baseURL string
	out     string
	client  *http.Client
	plan    []string
}

func newPublishing(name string, args []string, extra func(*flag.FlagSet)) (*publishing, []string, error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	keysFile := flags.String("keys", "", "trusted keys file the CLI embeds")
	baseURL := flags.String("base-url", alphasync.DefaultBaseURL, "where the bucket is read from")
	out := flags.String("out", "", "new dir for the files to upload and "+planFile)
	env := flags.String("env", "", "environment of the tenants: prod or pre-prod")
	tenantsFile := flags.String("tenants", "", "file with one tenant id per line; '#' starts a comment")
	var given tenantFlags
	flags.Var(&given, "tenant", "tenant id; may be repeated")
	extra(flags)
	if err := flags.Parse(args); err != nil {
		return nil, nil, err
	}
	if *keysFile == "" || *out == "" || *env == "" {
		return nil, nil, fmt.Errorf("%s needs --keys, --out and --env", name)
	}
	if err := alphasync.CheckEnv(*env); err != nil {
		return nil, nil, err
	}

	encoded, err := os.ReadFile(*keysFile)
	if err != nil {
		return nil, nil, err
	}
	keys, err := alphasync.ParseKeys(encoded)
	if err != nil {
		return nil, nil, err
	}

	tenants, err := readTenants(*tenantsFile, given)
	if err != nil {
		return nil, nil, err
	}
	if err := os.Mkdir(*out, 0o755); err != nil {
		return nil, nil, err
	}

	p := &publishing{
		keys:    keys,
		env:     *env,
		baseURL: strings.TrimSuffix(*baseURL, "/"),
		out:     *out,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
	return p, tenants, nil
}

func grant(args []string, stdout io.Writer) error {
	var specFile, itemFile *string
	p, tenants, err := newPublishing("grant", args, func(flags *flag.FlagSet) {
		specFile = flags.String("spec", "", "compiled spec")
		itemFile = flags.String("item", "", "signed item printed by sign")
	})
	if err != nil {
		return err
	}
	if *specFile == "" || *itemFile == "" {
		return errors.New("grant needs --spec and --item")
	}

	item, err := p.readItem(*itemFile)
	if err != nil {
		return err
	}
	if item.Env != p.env {
		return fmt.Errorf("%s was signed for %q, not for %q: sign it again with --env %s", *itemFile, item.Env, p.env, p.env)
	}
	spec, err := os.ReadFile(*specFile)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(spec)
	if hex.EncodeToString(sum[:]) != strings.ToLower(item.SHA256) {
		return fmt.Errorf("%s is not the spec %s signed: sign it again", *specFile, *itemFile)
	}

	published, err := p.specPublished(item, spec)
	if err != nil {
		return err
	}
	if published {
		fmt.Fprintf(stdout, "spec %s is already published\n", item.Path)
	} else {
		if err := p.add(item.Path, spec); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "spec %s will be uploaded\n", item.Path)
	}

	return p.updateManifests(tenants, stdout, func(specs []alphasync.ManifestSpec) ([]alphasync.ManifestSpec, error) {
		published, ok := alphasync.Manifest{Specs: specs}.Find(item.Module, item.Version)
		switch {
		case !ok:
			return append(specs, item), nil
		case published == item:
			return specs, nil
		default:
			return nil, fmt.Errorf("%s@%s is already published with other content: sign it with a new --version", item.Module, item.Version)
		}
	})
}

func revoke(args []string, stdout io.Writer) error {
	var module, version *string
	p, tenants, err := newPublishing("revoke", args, func(flags *flag.FlagSet) {
		module = flags.String("module", "", "module to take away")
		version = flags.String("version", "", "only this version; every version when empty")
	})
	if err != nil {
		return err
	}
	if *module == "" {
		return errors.New("revoke needs --module")
	}

	return p.updateManifests(tenants, stdout, func(specs []alphasync.ManifestSpec) ([]alphasync.ManifestSpec, error) {
		return slices.DeleteFunc(specs, func(s alphasync.ManifestSpec) bool {
			return s.Module == *module && (*version == "" || s.Version == *version)
		}), nil
	})
}

func (p *publishing) updateManifests(tenants []string, stdout io.Writer, change func([]alphasync.ManifestSpec) ([]alphasync.ManifestSpec, error)) error {
	for _, tenant := range tenants {
		path, err := alphasync.ManifestPath(p.env, tenant)
		if err != nil {
			return fmt.Errorf("tenant %s: %w", tenant, err)
		}
		current, err := p.readManifest(path)
		if err != nil {
			return fmt.Errorf("tenant %s: %w", tenant, err)
		}

		next, err := change(slices.Clone(current.Specs))
		if err != nil {
			return fmt.Errorf("tenant %s: %w", tenant, err)
		}
		slices.SortFunc(next, func(a, b alphasync.ManifestSpec) int {
			return cmp.Or(strings.Compare(a.Module, b.Module), strings.Compare(a.PublishedAt, b.PublishedAt))
		})
		if err := (alphasync.Manifest{Specs: next}).Validate(); err != nil {
			return fmt.Errorf("tenant %s: %w", tenant, err)
		}

		for _, spec := range next {
			if err := p.keys.Verify(spec); err != nil {
				return fmt.Errorf("tenant %s: %w", tenant, err)
			}
			if spec.Env != p.env {
				return fmt.Errorf("tenant %s: module %q was signed for %q, not for %q", tenant, spec.Module, spec.Env, p.env)
			}
		}

		before, after := labels(current.Specs), labels(next)
		if slices.EqualFunc(current.Specs, next, func(a, b alphasync.ManifestSpec) bool { return a == b }) {
			fmt.Fprintf(stdout, "tenant %s unchanged: %s\n", tenant, before)
			continue
		}

		data, err := json.MarshalIndent(alphasync.Manifest{Specs: next}, "", "  ")
		if err != nil {
			return err
		}
		if err := p.add(path, append(data, '\n')); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "tenant %s: %s -> %s (%s)\n", tenant, before, after, path)
	}

	return os.WriteFile(filepath.Join(p.out, planFile), []byte(strings.Join(p.plan, "")), 0o644)
}

func (p *publishing) readItem(file string) (alphasync.ManifestSpec, error) {
	var item alphasync.ManifestSpec
	data, err := os.ReadFile(file)
	if err != nil {
		return item, err
	}
	if err := json.Unmarshal(data, &item); err != nil {
		return item, fmt.Errorf("%s is not an item printed by sign: %w", file, err)
	}
	return item, p.keys.Verify(item)
}

func (p *publishing) specPublished(item alphasync.ManifestSpec, spec []byte) (bool, error) {
	body, found, err := p.fetch(item.Path)
	if err != nil || !found {
		return false, err
	}
	if !slices.Equal(body, spec) {
		return false, fmt.Errorf("%s already holds another spec: publish this one under a new path", item.Path)
	}
	return true, nil
}

func (p *publishing) readManifest(path string) (alphasync.Manifest, error) {
	var manifest alphasync.Manifest
	body, found, err := p.fetch(path)
	if err != nil || !found {
		return manifest, err
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		return manifest, fmt.Errorf("manifest %s does not parse: %w", path, err)
	}
	if err := manifest.Validate(); err != nil {
		return manifest, fmt.Errorf("manifest %s is not valid: %w", path, err)
	}
	return manifest, nil
}

func (p *publishing) fetch(path string) ([]byte, bool, error) {
	resp, err := p.client.Get(p.baseURL + "/" + path)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	// Only a 404 means missing: taking a 403 for it would write a manifest
	// without what the tenant already has.
	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		return body, true, err
	case http.StatusNotFound:
		return nil, false, nil
	default:
		return nil, false, fmt.Errorf("bucket answered %d for %s", resp.StatusCode, path)
	}
}

func (p *publishing) add(path string, data []byte) error {
	local := filepath.Join(p.out, "files", filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(local, data, 0o644); err != nil {
		return err
	}
	p.plan = append(p.plan, local+"\t"+path+"\n")
	return nil
}

// readTenants refuses an id that is not a UUID v4: the manifest path is only
// as hard to guess as the random bits of the tenant id.
func readTenants(file string, given []string) ([]string, error) {
	raw := slices.Clone(given)
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line, _, _ := strings.Cut(scanner.Text(), "#")
			if line = strings.TrimSpace(line); line != "" {
				raw = append(raw, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	var tenants []string
	for _, line := range raw {
		id, err := uuid.Parse(line)
		if err != nil || id.Version() != 4 || id.Variant() != uuid.RFC4122 {
			return nil, fmt.Errorf("%q is not a tenant id (UUID v4)", line)
		}
		if tenant := id.String(); !slices.Contains(tenants, tenant) {
			tenants = append(tenants, tenant)
		}
	}
	if len(tenants) == 0 {
		return nil, errors.New("no tenant: use --tenants or --tenant")
	}
	return tenants, nil
}

func labels(specs []alphasync.ManifestSpec) string {
	names := []string{}
	for _, spec := range specs {
		names = append(names, spec.Module+"@"+spec.Version)
	}
	return "[" + strings.Join(names, " ") + "]"
}
