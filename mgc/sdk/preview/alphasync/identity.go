package alphasync

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
)

// An API key does not say which tenant it belongs to, but the default security
// group of the tenant does: mgc runs this command in the place of the user.
// The tree names it as the spec does, and the command line in kebab case.
var SecurityGroupsList = []string{"network", "security_groups", "list"}

const securityGroupsCommand = "'mgc network security-groups list'"

// resolveTenant finds whose manifest to read: the tenant of the API key or,
// without one, the tenant of the login.
func resolveTenant(ctx context.Context) (string, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return "", errors.New("programming error: unable to retrieve auth configuration from context")
	}
	if apiKey, _ := auth.ApiKey(ctx); apiKey != "" {
		return apiKeyTenant(ctx)
	}
	if login, _ := auth.CurrentTenantID(); login != "" {
		return login, nil
	}
	return "", errors.New("you are not logged in: run 'mgc auth login', or pass --api-key")
}

func apiKeyTenant(ctx context.Context) (string, error) {
	command := securityGroupsCommand
	exec, err := findExecutor(core.GrouperFromContext(ctx), SecurityGroupsList)
	if err != nil {
		return "", fmt.Errorf("finding the tenant of an API key needs %s: %w", command, err)
	}

	configs := configsOf(ctx, exec)
	var lastErr error
	for _, region := range regionsOf(exec, configs) {
		if region != nil {
			configs["region"] = region
		}
		result, err := exec.Execute(ctx, core.Parameters{}, configs)
		if err != nil {
			var httpErr *mgcHttpPkg.HttpError
			if errors.As(err, &httpErr) && (httpErr.Code == http.StatusUnauthorized || httpErr.Code == http.StatusForbidden) {
				return "", fmt.Errorf("mgc finds the tenant of an API key with %s, so the key needs the network.read scope: %w", command, err)
			}
			lastErr = fmt.Errorf("%s in region %v: %w", command, region, err)
			continue
		}
		if tenant := defaultGroupTenant(result); tenant != "" {
			return tenant, nil
		}
	}

	if lastErr != nil {
		return "", lastErr
	}
	return "", errors.New("the tenant of the API key was not found: it has no default security group in any region")
}

func findExecutor(root core.Grouper, path []string) (core.Executor, error) {
	if root == nil {
		return nil, errors.New("programming error: unable to retrieve the commands from context")
	}

	var node core.Descriptor = root
	for _, name := range path {
		group, ok := node.(core.Grouper)
		if !ok {
			return nil, fmt.Errorf("%q has no subcommands", node.Name())
		}
		child, err := group.GetChildByName(name)
		if err != nil {
			return nil, err
		}
		node = child
	}

	exec, ok := node.(core.Executor)
	if !ok {
		return nil, fmt.Errorf("%q is not a command", node.Name())
	}
	return exec, nil
}

// configsOf reads the configs of a command from the CLI config, as the command
// line does, so the call goes to the same environment and server.
func configsOf(ctx context.Context, exec core.Executor) core.Configs {
	configs := core.Configs{}
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return configs
	}
	for name := range exec.ConfigsSchema().Properties {
		var value any
		if err := config.Get(name, &value); err == nil && value != nil {
			configs[name] = value
		}
	}
	return configs
}

// regionsOf lists the configured region first, then every other one the
// command knows: a tenant may have no security group where the CLI points.
func regionsOf(exec core.Executor, configs core.Configs) []any {
	var regions []any
	if region, ok := configs["region"]; ok {
		regions = append(regions, region)
	}
	if prop, ok := exec.ConfigsSchema().Properties["region"]; ok && prop.Value != nil {
		for _, region := range prop.Value.Enum {
			if !slices.Contains(regions, region) {
				regions = append(regions, region)
			}
		}
	}
	if len(regions) == 0 {
		return []any{nil}
	}
	return regions
}

func defaultGroupTenant(result core.Result) string {
	withValue, ok := core.ResultAs[core.ResultWithValue](result)
	if !ok {
		return ""
	}
	value, _ := withValue.Value().(map[string]any)
	groups, _ := value["security_groups"].([]any)
	for _, item := range groups {
		group, _ := item.(map[string]any)
		if isDefault, _ := group["is_default"].(bool); !isDefault {
			continue
		}
		if tenant, _ := group["tenant_id"].(string); tenant != "" {
			return tenant
		}
	}
	return ""
}
