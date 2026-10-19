package alphasync

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
)

// Target is a module to install, at a version or, when empty, its newest one.
type Target struct {
	Module  string
	Version string
}

func ParseTarget(arg string) (Target, error) {
	module, version, pinned := strings.Cut(arg, "@")
	if err := CheckModule(module); err != nil {
		return Target{}, err
	}
	if pinned {
		if err := CheckVersion(version); err != nil {
			return Target{}, err
		}
	}
	return Target{Module: module, Version: version}, nil
}

type Result struct {
	Installed []string  `json:"installed"`
	Removed   []string  `json:"removed,omitempty"`
	Skipped   []Skipped `json:"skipped,omitempty"`
}

type Skipped struct {
	Module string `json:"module"`
	Reason string `json:"reason"`
}

// Listed is a module the tenant was given, installed or not, or one installed
// from another tenant.
type Listed struct {
	Module string `json:"module"`
	Status string `json:"status"`
	// Always present, even empty, so that a table keeps its rows aligned.
	Installed string   `json:"installed"`
	Latest    string   `json:"latest"`
	Versions  []string `json:"versions,omitempty"`
}

const (
	statusNotInstalled    = "not installed"
	statusUpToDate        = "up to date"
	statusUpdateAvailable = "update available"
	statusNeedsNewerCLI   = "needs newer mgc"
	statusPinned          = "pinned"
	statusTakenAway       = "taken away"
	statusOtherTenant     = "other tenant"
)

// The same answer whether the module exists or not: an unreleased module is
// never confirmed to anyone it was not given to.
func notAvailable(target Target) error {
	if target.Version != "" {
		return fmt.Errorf("alpha module %s@%s is not available for your tenant", target.Module, target.Version)
	}
	return fmt.Errorf("alpha module %q is not available for your tenant", target.Module)
}

// Install installs modules from the manifest of a tenant: the ones named, or
// every one not installed yet. A module named with a version stays pinned to
// it; what is installed already only moves when named, or on upgrade.
func Install(ctx context.Context, cfg Config, manifestPath, env string, targets []Target) (*Result, error) {
	if cfg.BaseURL == "" {
		return nil, errNoSource
	}
	manifest, err := fetchManifest(ctx, cfg, manifestPath)
	if err != nil {
		return nil, err
	}
	index, err := readIndex(cfg.Dir)
	if err != nil {
		return nil, err
	}

	result := &Result{Installed: []string{}}
	var steps []step
	add := func(entry Installed) error {
		next, err := prepare(ctx, cfg, index, entry, env)
		if err != nil {
			return err
		}
		steps = append(steps, next)
		result.Installed = append(result.Installed, entry.label())
		return nil
	}

	if len(targets) == 0 {
		for _, spec := range manifest.Latest() {
			if _, ok := index[spec.Module]; ok {
				continue
			}
			if reason := tooOld(cfg.CLIVersion, spec.MinCLIVersion); reason != "" {
				result.Skipped = append(result.Skipped, Skipped{Module: spec.Module, Reason: reason})
				continue
			}
			if err := add(Installed{spec, manifestPath, false}); err != nil {
				return nil, err
			}
		}
	}
	named := map[string]Target{}
	for _, target := range targets {
		if previous, ok := named[target.Module]; ok {
			if previous != target {
				return nil, fmt.Errorf("alpha module %q is named more than once", target.Module)
			}
			continue
		}
		named[target.Module] = target

		spec, ok := manifest.Find(target.Module, target.Version)
		if !ok {
			return nil, notAvailable(target)
		}
		if reason := tooOld(cfg.CLIVersion, spec.MinCLIVersion); reason != "" {
			return nil, fmt.Errorf("alpha module %s %s", spec.label(), reason)
		}
		if err := add(Installed{spec, manifestPath, target.Version != ""}); err != nil {
			return nil, err
		}
	}
	return result, apply(cfg.Dir, index, steps)
}

// Upgrade moves installed modules to the newest version their manifest lists,
// and removes the ones it takes away. Without names, a pinned module keeps its
// version, unless that version is taken away; a module named is upgraded and
// unpinned.
func Upgrade(ctx context.Context, cfg Config, modules []string) (*Result, error) {
	if cfg.BaseURL == "" {
		return nil, errNoSource
	}
	index, err := readIndex(cfg.Dir)
	if err != nil {
		return nil, err
	}
	named := len(modules) > 0
	for _, module := range modules {
		if _, ok := index[module]; !ok {
			return nil, fmt.Errorf("alpha module %q is not installed", module)
		}
	}
	if !named {
		modules = slices.Collect(maps.Keys(index))
	}
	slices.Sort(modules)
	modules = slices.Compact(modules)

	manifests := manifestCache(ctx, cfg)
	result := &Result{Installed: []string{}}
	var steps []step
	for _, module := range modules {
		entry := index[module]
		manifest, err := manifests(entry.Manifest)
		if err != nil {
			return nil, err
		}

		newest, listed := manifest.Find(module, "")
		reason := tooOld(cfg.CLIVersion, newest.MinCLIVersion)
		switch {
		case !listed, !named && takenAway(manifest, entry):
			steps = append(steps, step{module: module})
			result.Removed = append(result.Removed, module)
		case entry.Pinned && !named:
		case newest == entry.ManifestSpec:
			if entry.Pinned {
				entry.Pinned = false
				steps = append(steps, step{module: module, entry: &entry})
			}
		case reason != "":
			result.Skipped = append(result.Skipped, Skipped{Module: module, Reason: reason})
		default:
			next, err := prepare(ctx, cfg, index, Installed{newest, entry.Manifest, false}, entry.Env)
			if err != nil {
				return nil, err
			}
			steps = append(steps, next)
			result.Installed = append(result.Installed, newest.label())
		}
	}
	return result, apply(cfg.Dir, index, steps)
}

// step is one change to what is installed. Every download and check happens
// before the first change, so a command that fails leaves everything as it was.
type step struct {
	module string
	// Nil removes the module.
	entry *Installed
	// Nil keeps the installed spec, which is the same.
	spec []byte
}

func prepare(ctx context.Context, cfg Config, index map[string]Installed, entry Installed, env string) (step, error) {
	if current, ok := index[entry.Module]; ok && current.ManifestSpec == entry.ManifestSpec && intact(cfg.Dir, current) {
		return step{module: entry.Module, entry: &entry}, nil
	}
	spec, err := download(ctx, cfg, entry.ManifestSpec, env)
	if err != nil {
		return step{}, err
	}
	return step{module: entry.Module, entry: &entry, spec: spec}, nil
}

func Uninstall(cfg Config, module string) error {
	index, err := readIndex(cfg.Dir)
	if err != nil {
		return err
	}
	if _, ok := index[module]; !ok {
		return fmt.Errorf("alpha module %q is not installed", module)
	}
	return apply(cfg.Dir, index, []step{{module: module}})
}

// takenAway tells whether a manifest no longer gives an installed module: none
// of its versions, or not the one it is pinned to.
func takenAway(manifest *Manifest, entry Installed) bool {
	version := ""
	if entry.Pinned {
		version = entry.Version
	}
	_, ok := manifest.Find(entry.Module, version)
	return !ok
}

// List shows what the manifest of a tenant gives and what is installed. A module
// installed from that manifest is compared with it, as upgrade does; one from
// another tenant's manifest is only told apart.
func List(ctx context.Context, cfg Config, manifestPath string) ([]Listed, error) {
	if cfg.BaseURL == "" {
		return nil, errNoSource
	}
	index, err := readIndex(cfg.Dir)
	if err != nil {
		return nil, err
	}
	given, err := fetchManifest(ctx, cfg, manifestPath)
	if err != nil {
		return nil, err
	}

	listed := map[string]Listed{}
	for _, spec := range given.Latest() {
		item := Listed{Module: spec.Module, Status: statusNotInstalled, Latest: spec.Version, Versions: given.Versions(spec.Module)}
		entry, installed := index[spec.Module]
		switch {
		case installed && entry.Manifest != manifestPath:
			item.Status = statusOtherTenant
		case installed && takenAway(given, entry):
			item.Status = statusTakenAway
		case installed && entry.Pinned:
			item.Status = statusPinned
		case installed && spec.Version == entry.Version:
			item.Status = statusUpToDate
		case tooOld(cfg.CLIVersion, spec.MinCLIVersion) != "":
			item.Status = statusNeedsNewerCLI
		case installed:
			item.Status = statusUpdateAvailable
		}
		item.Installed = entry.Version
		listed[spec.Module] = item
	}
	for module, entry := range index {
		if _, ok := listed[module]; ok {
			continue
		}
		item := Listed{Module: module, Status: statusTakenAway, Installed: entry.Version}
		if entry.Manifest != manifestPath {
			item.Status = statusOtherTenant
		}
		listed[module] = item
	}

	items := []Listed{}
	for _, module := range slices.Sorted(maps.Keys(listed)) {
		items = append(items, listed[module])
	}
	return items, nil
}

func scopesOf(cfg Config, modules []string) (core.Scopes, error) {
	store := someModules{InstalledStore(cfg.Dir, cfg.Keys), modules}
	scopes := core.Scopes{}
	_, err := core.VisitAllExecutors(cfg.Group(store), nil, false, func(exec core.Executor, _ []string) (bool, error) {
		scopes.Add(exec.Scopes()...)
		return true, nil
	})
	return scopes, err
}

type someModules struct {
	preview.Store
	names []string
}

func (s someModules) Names() ([]string, error) {
	return s.names, nil
}
