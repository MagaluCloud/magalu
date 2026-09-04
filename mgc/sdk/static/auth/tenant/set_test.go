package tenant

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
)

func newTestConfig(t *testing.T) *mgcConfigPkg.Config {
	t.Helper()
	pm, _ := profile_manager.NewInMemoryProfileManager()
	return mgcConfigPkg.New(pm)
}

// tokenWithTenant monta um access token sintetico. O parser da CLI usa
// ParseUnverified, entao a assinatura nao precisa ser valida — o que importa e a
// claim `tenant`, no formato <tipo>.<uuid> que o IDP devolve.
func tokenWithTenant(t *testing.T, claim string) string {
	t.Helper()

	enc := func(s string) string {
		return base64.RawURLEncoding.EncodeToString([]byte(s))
	}

	payload := `{}`
	if claim != "" {
		payload = `{"tenant":"` + claim + `"}`
	}

	return enc(`{"alg":"none","typ":"JWT"}`) + "." + enc(payload) + ".x"
}

// newTestAuth devolve um Auth carregado com o token dado, sem tocar em rede nem
// no perfil real.
func newTestAuth(t *testing.T, claim string) *mgcAuthPkg.Auth {
	t.Helper()

	// InitTokensFromFile deixa MGC_ACCESS_TOKEN sobrescrever o token; zerar isola
	// o teste do ambiente de quem roda.
	t.Setenv("MGC_ACCESS_TOKEN", "")

	pm, _ := profile_manager.NewInMemoryProfileManager()
	auth := mgcAuthPkg.New(
		map[string]mgcAuthPkg.Config{},
		&http.Client{},
		pm,
		mgcConfigPkg.New(pm),
	)

	if claim != "" {
		if err := auth.SetTokens(&mgcAuthPkg.LoginResult{AccessToken: tokenWithTenant(t, claim)}); err != nil {
			t.Fatal(err)
		}
	}

	return auth
}

// O projeto passa a ser o default do tenant, que no IAM tem o mesmo id do
// tenant. O tipo que prefixa a claim nao entra: so o id.
func TestSelectDefaultProjectWritesTenantID(t *testing.T) {
	cfg := newTestConfig(t)
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	if err := SelectDefaultProject(ctx, newTestAuth(t, "business.tenant-alpha")); err != nil {
		t.Fatalf("SelectDefaultProject: %v", err)
	}

	if got, want := cfg.Project(), "tenant-alpha"; got != want {
		t.Errorf("project = %q, quer %q", got, want)
	}
}

// O ponto da funcao: um projeto escolhido para o tenant anterior nao pode
// sobreviver a troca, senao a CLI manda o x-project-id de um tenant no token de
// outro.
func TestSelectDefaultProjectReplacesPreviousSelection(t *testing.T) {
	cfg := newTestConfig(t)
	if err := cfg.SetProject("projeto-do-tenant-antigo"); err != nil {
		t.Fatal(err)
	}
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	if err := SelectDefaultProject(ctx, newTestAuth(t, "business.tenant-beta")); err != nil {
		t.Fatalf("SelectDefaultProject: %v", err)
	}

	if got, want := cfg.Project(), "tenant-beta"; got != want {
		t.Errorf("project = %q, quer %q", got, want)
	}
}

// Token sem claim de tenant nao deveria acontecer, mas se acontecer o menos
// errado e ficar sem projeto — nunca manter o do tenant anterior.
func TestSelectDefaultProjectWithoutTenantClaimClearsIt(t *testing.T) {
	cfg := newTestConfig(t)
	if err := cfg.SetProject("projeto-do-tenant-antigo"); err != nil {
		t.Fatal(err)
	}
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	if err := SelectDefaultProject(ctx, newTestAuth(t, "")); err != nil {
		t.Fatalf("sem claim nao pode ser erro: %v", err)
	}

	if got := cfg.Project(); got != "" {
		t.Errorf("project deveria estar limpo, veio %q", got)
	}
}

func TestSelectDefaultProjectWithoutConfig(t *testing.T) {
	if err := SelectDefaultProject(context.Background(), newTestAuth(t, "business.tenant-alpha")); err == nil {
		t.Error("contexto sem config deveria falhar")
	}
}
