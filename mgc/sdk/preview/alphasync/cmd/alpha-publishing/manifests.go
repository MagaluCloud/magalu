package main

import (
	"bufio"
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
)

const planFile = "plan.tsv"

type inviteFlags []string

func (i *inviteFlags) String() string     { return strings.Join(*i, ",") }
func (i *inviteFlags) Set(v string) error { *i = append(*i, v); return nil }

type publishing struct {
	keys    alphasync.Keys
	baseURL string
	out     string
	client  *http.Client
	plan    []string
}

func newPublishing(name string, args []string, extra func(*flag.FlagSet)) (*publishing, []string, *flag.FlagSet, error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	keysFile := flags.String("keys", "", "trusted keys file the CLI embeds")
	baseURL := flags.String("base-url", alphasync.DefaultBaseURL, "where the bucket is read from")
	out := flags.String("out", "", "new dir for the files to upload and "+planFile)
	invitesFile := flags.String("invites", "", "file with one invite per line; '#' starts a comment")
	var given inviteFlags
	flags.Var(&given, "invite", "invite; may be repeated")
	extra(flags)
	if err := flags.Parse(args); err != nil {
		return nil, nil, nil, err
	}
	if *keysFile == "" || *out == "" {
		return nil, nil, nil, fmt.Errorf("%s needs --keys and --out", name)
	}

	encoded, err := os.ReadFile(*keysFile)
	if err != nil {
		return nil, nil, nil, err
	}
	keys, err := alphasync.ParseKeys(encoded)
	if err != nil {
		return nil, nil, nil, err
	}

	invites, err := readInvites(*invitesFile, given)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := os.Mkdir(*out, 0o755); err != nil {
		return nil, nil, nil, err
	}

	p := &publishing{
		keys:    keys,
		baseURL: strings.TrimSuffix(*baseURL, "/"),
		out:     *out,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
	return p, invites, flags, nil
}

func grant(args []string, stdout io.Writer) error {
	var specFile, itemFile *string
	p, invites, _, err := newPublishing("grant", args, func(flags *flag.FlagSet) {
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

	return p.updateManifests(invites, stdout, func(specs []alphasync.ManifestSpec) []alphasync.ManifestSpec {
		return append(withoutModule(specs, item.Module), item)
	})
}

func revoke(args []string, stdout io.Writer) error {
	var module *string
	p, invites, _, err := newPublishing("revoke", args, func(flags *flag.FlagSet) {
		module = flags.String("module", "", "module to take away")
	})
	if err != nil {
		return err
	}
	if *module == "" {
		return errors.New("revoke needs --module")
	}

	return p.updateManifests(invites, stdout, func(specs []alphasync.ManifestSpec) []alphasync.ManifestSpec {
		return withoutModule(specs, *module)
	})
}

func (p *publishing) updateManifests(invites []string, stdout io.Writer, change func([]alphasync.ManifestSpec) []alphasync.ManifestSpec) error {
	for _, invite := range invites {
		path := alphasync.ManifestPath(invite)
		current, err := p.readManifest(path)
		if err != nil {
			return fmt.Errorf("invite %s: %w", invite, err)
		}

		next := change(slices.Clone(current.Specs))
		slices.SortFunc(next, func(a, b alphasync.ManifestSpec) int { return strings.Compare(a.Module, b.Module) })

		for _, spec := range next {
			if err := p.keys.Verify(spec); err != nil {
				return fmt.Errorf("invite %s: %w", invite, err)
			}
		}

		before, after := modules(current.Specs), modules(next)
		if slices.EqualFunc(current.Specs, next, func(a, b alphasync.ManifestSpec) bool { return a == b }) {
			fmt.Fprintf(stdout, "invite %s unchanged: %s\n", invite, before)
			continue
		}

		data, err := json.MarshalIndent(alphasync.Manifest{Specs: next}, "", "  ")
		if err != nil {
			return err
		}
		if err := p.add(path, append(data, '\n')); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "invite %s: %s -> %s (%s)\n", invite, before, after, path)
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
	return manifest, nil
}

func (p *publishing) fetch(path string) ([]byte, bool, error) {
	resp, err := p.client.Get(p.baseURL + "/" + path)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		return body, true, err
	case http.StatusNotFound, http.StatusForbidden:
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

func readInvites(file string, given []string) ([]string, error) {
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

	var invites []string
	for _, line := range raw {
		invite, err := alphasync.NormalizeInvite(line)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(invites, invite) {
			invites = append(invites, invite)
		}
	}
	if len(invites) == 0 {
		return nil, errors.New("no invite: use --invites or --invite")
	}
	return invites, nil
}

func withoutModule(specs []alphasync.ManifestSpec, module string) []alphasync.ManifestSpec {
	return slices.DeleteFunc(specs, func(s alphasync.ManifestSpec) bool { return s.Module == module })
}

func modules(specs []alphasync.ManifestSpec) string {
	names := []string{}
	for _, spec := range specs {
		names = append(names, spec.Module)
	}
	return "[" + strings.Join(names, " ") + "]"
}
