package alphasync

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
)

type installParams struct {
	Modules []string `json:"modules,omitempty" jsonschema_description:"Modules as module or module@version; every module not installed yet when left out" mgc:"positional"`
}

type upgradeParams struct {
	Modules []string `json:"modules,omitempty" jsonschema_description:"Modules to upgrade; every module not pinned to a version when left out" mgc:"positional"`
}

type uninstallParams struct {
	Module string `json:"module" jsonschema_description:"Module to remove from this machine" mgc:"positional"`
}

type uninstallResult struct {
	Uninstalled string `json:"uninstalled"`
}

// NewPreviewRoot is 'mgc preview', which manages the alpha modules installed on
// this machine.
func NewPreviewRoot(cfg Config) core.Grouper {
	group := core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "preview",
			GroupID:     "settings",
			Summary:     "Install the preview commands",
			Description: `Lists, installs, upgrades and removes the alpha commands.`,
		},
		func() []core.Descriptor {
			return []core.Descriptor{newList(cfg), newInstall(cfg), newUpgrade(cfg), newUninstall(cfg)}
		},
	)
	return core.NewStaticGroup(
		core.DescriptorSpec{Name: "preview root", Description: "Preview commands that ship with the CLI"},
		func() []core.Descriptor { return []core.Descriptor{group} },
	)
}

func newList(cfg Config) core.Executor {
	exec := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "list",
			Summary:     "List the alpha modules of your tenant",
			Description: "Lists the alpha modules Magalu gave to your tenant, with their versions, and the ones installed on this machine",
		},
		func(ctx context.Context, _ struct{}, _ struct{}) ([]Listed, error) {
			path, _, err := tenantManifest(ctx)
			if err != nil {
				return nil, err
			}
			return List(ctx, cfg, path)
		},
	)
	return core.NewExecuteResultOutputOptions(exec, func(core.Executor, core.Result) string {
		return "table=MODULE:$[*].module,STATUS:$[*].status,INSTALLED:$[*].installed,LATEST:$[*].latest"
	})
}

func newInstall(cfg Config) core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:    "install",
			Summary: "Install alpha modules of your tenant",
			Description: `Installs alpha modules of your tenant at their newest version, or, without
names, every one not installed yet. A module named with a version, as
iam@2026-10-08, stays at it until upgraded by name`,
			Observations: "With an API key, mgc finds its tenant by listing the security groups of the key, which needs the network.read scope.",
		},
		func(ctx context.Context, params installParams, _ struct{}) (*Result, error) {
			targets := make([]Target, 0, len(params.Modules))
			for _, arg := range params.Modules {
				target, err := ParseTarget(arg)
				if err != nil {
					return nil, err
				}
				targets = append(targets, target)
			}
			path, env, err := tenantManifest(ctx)
			if err != nil {
				return nil, err
			}

			result, err := Install(ctx, cfg, path, env, targets)
			if err != nil {
				return nil, err
			}
			addScopes(ctx, cfg, result.Installed)
			return result, nil
		},
	)
}

func newUpgrade(cfg Config) core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:    "upgrade",
			Summary: "Upgrade the installed alpha modules",
			Description: `Moves installed modules to their newest version and removes the ones taken
away: a module no longer given, or the version a module is pinned to. Modules
pinned to a version only move when named`,
		},
		func(ctx context.Context, params upgradeParams, _ struct{}) (*Result, error) {
			result, err := Upgrade(ctx, cfg, params.Modules)
			if err != nil {
				return nil, err
			}
			addScopes(ctx, cfg, result.Installed)
			return result, nil
		},
	)
}

func newUninstall(cfg Config) core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "uninstall",
			Summary:     "Remove an alpha module from this machine",
			Description: "Removes an installed alpha module from this machine. Your tenant keeps access to it, and 'mgc preview install' brings it back",
		},
		func(_ context.Context, params uninstallParams, _ struct{}) (*uninstallResult, error) {
			if err := Uninstall(cfg, params.Module); err != nil {
				return nil, err
			}
			return &uninstallResult{Uninstalled: params.Module}, nil
		},
	)
}

func tenantManifest(ctx context.Context) (path, env string, err error) {
	tenant, err := resolveTenant(ctx)
	if err != nil {
		return "", "", err
	}
	env = envFromContext(ctx)
	path, err = ManifestPath(env, tenant)
	return path, env, err
}

func envFromContext(ctx context.Context) string {
	var env string
	if config := mgcConfigPkg.FromContext(ctx); config != nil {
		_ = config.Get("env", &env)
	}
	return CurrentEnv(env)
}

// addScopes gives the token of the login the scopes the installed modules
// need. An API key cannot change, so the user learns which scopes it must have.
func addScopes(ctx context.Context, cfg Config, labels []string) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil || cfg.Group == nil || len(labels) == 0 {
		return
	}
	modules := make([]string, len(labels))
	for i, label := range labels {
		modules[i], _, _ = strings.Cut(label, "@")
	}
	needed, err := scopesOf(cfg, modules)
	if err != nil || len(needed) == 0 {
		return
	}

	scopes, err := auth.CurrentScopes()
	key, _ := auth.ApiKey(ctx)
	login, _ := auth.CurrentTenantID()
	if key != "" || login == "" || err != nil {
		fmt.Fprintf(os.Stderr, "Warning: these commands need the scopes %s: log in, or use an API key that has them\n", needed.AsScopesString())
		return
	}
	missing := false
	for _, scope := range needed {
		if !slices.Contains(scopes, scope) {
			scopes.Add(scope)
			missing = true
		}
	}
	if !missing {
		return
	}
	if _, err := auth.SetScopes(ctx, scopes); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: these commands need scopes the login lacks, and adding them failed (run 'mgc auth login' again): %v\n", err)
	}
}
