package alphasync

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type securityGroupsConfigs struct {
	Region string `json:"region,omitempty" jsonschema:"enum=br-ne1,enum=br-se1"`
}

// networkTree answers 'network security-groups list' with the groups of each
// region.
func networkTree(t *testing.T, byRegion map[string]any, regions *[]string) core.Grouper {
	t.Helper()
	list := core.NewStaticExecute(
		core.DescriptorSpec{Name: "list", Description: "List security groups"},
		func(_ context.Context, _ struct{}, configs securityGroupsConfigs) (map[string]any, error) {
			*regions = append(*regions, configs.Region)
			switch value := byRegion[configs.Region].(type) {
			case error:
				return nil, value
			case []map[string]any:
				groups := []any{}
				for _, group := range value {
					groups = append(groups, group)
				}
				return map[string]any{"security_groups": groups}, nil
			default:
				return map[string]any{"security_groups": []any{}}, nil
			}
		},
	)
	group := func(name string, children ...core.Descriptor) core.Grouper {
		return core.NewStaticGroup(core.DescriptorSpec{Name: name, Description: name}, func() []core.Descriptor { return children })
	}
	return group("root", group("network", group("security_groups", list)))
}

func withTree(ctx context.Context, root core.Grouper) context.Context {
	return core.NewGrouperContext(ctx, func() core.Grouper { return root })
}

func TestAPIKeyTenantFromTheDefaultSecurityGroup(t *testing.T) {
	var regions []string
	root := networkTree(t, map[string]any{
		"br-ne1": []map[string]any{
			{"name": "web", "is_default": false, "tenant_id": tenantB},
			{"name": "default", "is_default": true, "tenant_id": tenantA},
		},
	}, &regions)

	tenant, err := apiKeyTenant(withTree(context.Background(), root))
	require.NoError(t, err)
	assert.Equal(t, tenantA, tenant, "only the default group counts")
	assert.Equal(t, []string{"br-ne1"}, regions)
}

func TestAPIKeyTenantTriesEveryRegion(t *testing.T) {
	var regions []string
	root := networkTree(t, map[string]any{
		"br-se1": []map[string]any{{"name": "default", "is_default": true, "tenant_id": tenantA}},
	}, &regions)

	tenant, err := apiKeyTenant(withTree(context.Background(), root))
	require.NoError(t, err)
	assert.Equal(t, tenantA, tenant)
	assert.Equal(t, []string{"br-ne1", "br-se1"}, regions)

	regions = nil
	root = networkTree(t, map[string]any{}, &regions)
	_, err = apiKeyTenant(withTree(context.Background(), root))
	assert.ErrorContains(t, err, "no default security group in any region")
	assert.Equal(t, []string{"br-ne1", "br-se1"}, regions)
}

func TestAPIKeyTenantErrors(t *testing.T) {
	var regions []string
	forbidden := &mgcHttpPkg.HttpError{Code: http.StatusForbidden, Status: "403 Forbidden"}
	root := networkTree(t, map[string]any{"br-ne1": forbidden}, &regions)
	_, err := apiKeyTenant(withTree(context.Background(), root))
	assert.ErrorContains(t, err, "network.read")
	assert.Equal(t, []string{"br-ne1"}, regions, "a key without the scope fails the same in every region")

	regions = nil
	failing := &mgcHttpPkg.HttpError{Code: http.StatusBadGateway, Status: "502 Bad Gateway"}
	root = networkTree(t, map[string]any{
		"br-ne1": failing,
		"br-se1": []map[string]any{{"is_default": true, "tenant_id": tenantA}},
	}, &regions)
	tenant, err := apiKeyTenant(withTree(context.Background(), root))
	require.NoError(t, err, "a region that fails does not hide the others")
	assert.Equal(t, tenantA, tenant)

	_, err = apiKeyTenant(withTree(context.Background(), core.NewStaticGroup(
		core.DescriptorSpec{Name: "root", Description: "root"},
		func() []core.Descriptor { return nil },
	)))
	assert.ErrorContains(t, err, "'mgc network security-groups list'")
}

// fakeJWT is unsigned: the CLI reads its claims without verifying them.
func fakeJWT(t *testing.T, tenant string, scopes string) string {
	t.Helper()
	claims, err := json.Marshal(map[string]any{
		"tenant": "pf." + tenant,
		"scope":  scopes,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	require.NoError(t, err)
	encode := base64.RawURLEncoding.EncodeToString
	return encode([]byte(`{"alg":"none","typ":"JWT"}`)) + "." + encode(claims) + ".sig"
}

// newAuthContext builds a context like the SDK's, with an in-memory config.
func newAuthContext(t *testing.T, token string, configs map[string]mgcAuthPkg.Config, transport http.RoundTripper) (context.Context, *mgcAuthPkg.Auth) {
	t.Helper()
	t.Setenv("MGC_ACCESS_TOKEN", token)
	manager, _ := profile_manager.NewInMemoryProfileManager()
	config := mgcConfigPkg.New(manager)
	if transport == nil {
		transport = http.DefaultTransport
	}
	auth := mgcAuthPkg.New(configs, &http.Client{Transport: transport}, manager, config)

	ctx := mgcConfigPkg.NewContext(context.Background(), config)
	ctx = mgcAuthPkg.NewContext(ctx, auth)
	ctx = mgcHttpPkg.NewClientContext(ctx, mgcHttpPkg.NewClient(transport))
	return ctx, auth
}

func TestResolveTenant(t *testing.T) {
	var regions []string
	root := networkTree(t, map[string]any{
		"br-ne1": []map[string]any{{"is_default": true, "tenant_id": tenantB}},
	}, &regions)

	ctx, _ := newAuthContext(t, fakeJWT(t, tenantA, "openid"), nil, nil)
	tenant, err := resolveTenant(ctx)
	require.NoError(t, err)
	assert.Equal(t, tenantA, tenant, "the tenant of the login")

	ctx, auth := newAuthContext(t, fakeJWT(t, tenantA, "openid"), nil, nil)
	require.NoError(t, auth.SetAPIKey("key"))
	tenant, err = resolveTenant(withTree(ctx, root))
	require.NoError(t, err)
	assert.Equal(t, tenantB, tenant, "the API key wins over the login")

	ctx, _ = newAuthContext(t, "", nil, nil)
	_, err = resolveTenant(ctx)
	assert.EqualError(t, err, "you are not logged in: run 'mgc auth login', or pass --api-key")
}

func TestAPIKeyTenantFromAWrappedResult(t *testing.T) {
	result := core.NewSimpleResult(core.ResultSource{}, nil, map[string]any{
		"security_groups": []any{map[string]any{"is_default": true, "tenant_id": tenantA}},
	})
	assert.Equal(t, tenantA, defaultGroupTenant(core.NewResultWithOriginalSource(core.ResultSource{}, result)))
}
