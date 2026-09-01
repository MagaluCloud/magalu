package project

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
)

const (
	prodHost    = "api.magalu.cloud"
	preProdHost = "api.pre-prod.jaxyendy.com:8443"

	// Mirrors the IAM spec server entry ("https://{env}/iam"): 'serverUrl' replaces the
	// whole base and the operation path is appended to whatever base ends up in use.
	basePath     = "/iam"
	projectsPath = "/api/v1/projects"
)

type projectConfig struct {
	mgcConfigPkg.NetworkConfig `json:",squash"` // nolint
}

type projectResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type projectCreate struct {
	Name string `json:"name"`
}

type createResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TenantID string `json:"tenant_id"`
}

func hostForEnv(env string) string {
	switch env {
	case "pre-prod", preProdHost:
		return preProdHost
	default:
		return prodHost
	}
}

func buildProjectsURL(serverUrl, env string) string {
	base := serverUrl
	if base == "" {
		base = "https://" + hostForEnv(env) + basePath
	}
	return strings.TrimSuffix(base, "/") + projectsPath
}

func envFromContext(ctx context.Context) string {
	c := mgcConfigPkg.FromContext(ctx)
	if c == nil {
		return ""
	}

	var env string
	if err := c.Get("env", &env); err != nil {
		return ""
	}

	return env
}

func projectsURL(ctx context.Context, cfg projectConfig) string {
	return buildProjectsURL(cfg.ServerUrl, envFromContext(ctx))
}

func authenticatedClient(ctx context.Context) (*mgcHttpPkg.Client, *mgcAuthPkg.Auth, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, nil, fmt.Errorf("programming error: unable to retrieve auth configuration from context")
	}

	if _, err := auth.AccessToken(ctx); err != nil {
		return nil, nil, fmt.Errorf("you are not logged in. To authenticate, please run 'mgc auth login'")
	}

	client := auth.AuthenticatedHttpClientFromContext(ctx)
	if client == nil {
		return nil, nil, fmt.Errorf("programming error: unable to retrieve HTTP Client from context")
	}

	return client, auth, nil
}

func newRequest(ctx context.Context, auth *mgcAuthPkg.Auth, method, url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// The XaaS gateway expects the tenant header alongside the bearer token, the same
	// way OpenAPI-backed operations handle it.
	if auth.CurrentSecurityMethod() == mgcAuthPkg.XTenantID.String() {
		xTenantID, err := auth.XTenantID(ctx)
		if err != nil {
			return nil, err
		}
		r.Header.Set("x-tenant-id", xTenantID)
	}

	return r, nil
}
