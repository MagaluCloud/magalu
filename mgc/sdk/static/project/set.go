package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/google/uuid"
)

type setParams struct {
	IDOrName string `json:"id" jsonschema:"description=UUID or name of the project to use. Run 'mgc alpha project list' to see the available ones,required,example=my-project" mgc:"positional"`
}

type setResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Warning string `json:"warning,omitempty"`
}

const warnNotVerified = "not verified: without a user login the API checks the project on the next request"

var getSet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:         "set",
			Summary:      "Set the project used by the CLI",
			Description:  "Subsequent Virtual Machine, Kubernetes and Block Storage requests are scoped to this project. To go back to the tenant's default, run 'mgc alpha project unset'",
			Observations: "Changing the tenant keeps the selected project, and the API refuses a project from another tenant: run 'mgc alpha project unset' or pick another one. Without a user login, as with an API key, only a project id is accepted, and the API checks it on the next request.",
		},
		set,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Success! Current project changed to {{if .name}}{{.name}} ({{.id}}){{else}}{{.id}}{{end}}\n{{if .warning}}Warning: {{.warning}}\n{{end}}"
	})
})

func set(ctx context.Context, params setParams, cfg projectConfig) (*setResult, error) {
	return setScope(ctx, mgcConfigPkg.ProjectKey, params.IDOrName, cfg)
}

func setScope(ctx context.Context, key, query string, cfg projectConfig) (*setResult, error) {
	if !hasAccessToken(ctx) {
		return applyUnverified(ctx, key, query)
	}

	available, err := list(ctx, listParams{}, cfg)
	if err != nil {
		return nil, err
	}
	return applyScope(ctx, key, query, available)
}

func hasAccessToken(ctx context.Context) bool {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return false
	}
	_, err := auth.AccessToken(ctx)
	return err == nil
}

func applyUnverified(ctx context.Context, key, query string) (*setResult, error) {
	id, err := uuid.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("without a user login only a project id (UUID) is accepted, since names are resolved by the API. Run 'mgc auth login' to select a project by name")
	}

	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}
	if err := config.Set(key, id.String()); err != nil {
		return nil, err
	}

	return &setResult{ID: id.String(), Warning: warnNotVerified}, nil
}

func applyScope(ctx context.Context, key, query string, available []projectResult) (*setResult, error) {
	found, err := resolveProject(available, query)
	if err != nil {
		return nil, err
	}

	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return nil, fmt.Errorf("programming error: unable to retrieve config from context")
	}
	if err := config.Set(key, found.ID); err != nil {
		return nil, err
	}

	return &setResult{ID: found.ID, Name: found.Name}, nil
}

func resolveProject(available []projectResult, query string) (*projectResult, error) {
	if len(available) == 0 {
		return nil, fmt.Errorf("no projects available for this tenant. Create one with 'mgc alpha project create'")
	}

	for i := range available {
		if available[i].ID == query {
			return &available[i], nil
		}
	}

	var matches []projectResult
	for _, p := range available {
		if strings.EqualFold(p.Name, query) {
			matches = append(matches, p)
		}
	}

	switch len(matches) {
	case 1:
		return &matches[0], nil
	case 0:
		return nil, fmt.Errorf("project %q not found.\nAvailable projects:%s", query, projectLines(available))
	default:
		return nil, fmt.Errorf("more than one project named %q. Use the id instead:%s", query, projectLines(matches))
	}
}

func projectLines(projects []projectResult) string {
	var b strings.Builder
	for _, p := range projects {
		fmt.Fprintf(&b, "\n  %s  %s", p.ID, p.Name)
	}
	return b.String()
}
