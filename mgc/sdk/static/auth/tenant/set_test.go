package tenant

import (
	"context"
	"testing"

	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
)

func newTestConfig(t *testing.T) *mgcConfigPkg.Config {
	t.Helper()
	pm, _ := profile_manager.NewInMemoryProfileManager()
	return mgcConfigPkg.New(pm)
}

// O projeto pertence ao tenant: mantê-lo após a troca apontaria para um escopo
// de outra conta. Os DOIS escopos caem juntos — nenhum sobrevive à troca.
func TestUnsetProjectsOnTenantChangeClearsBoth(t *testing.T) {
	cfg := newTestConfig(t)
	if err := cfg.Set(mgcConfigPkg.ProjectKey, "p"); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Set(mgcConfigPkg.IamProjectKey, "i"); err != nil {
		t.Fatal(err)
	}
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	had, err := unsetProjectsForTenantChange(ctx)
	if err != nil {
		t.Fatalf("unsetProjectsForTenantChange: %v", err)
	}
	if !had {
		t.Error("deveria reportar que havia projeto selecionado (é o que dispara o aviso)")
	}
	if got := cfg.Project(); got != "" {
		t.Errorf("project deveria estar limpo, veio %q", got)
	}
	if got := cfg.IamProject(); got != "" {
		t.Errorf("iamProject deveria estar limpo, veio %q", got)
	}
}

// Um só dos escopos setado ainda conta como "havia projeto" — o aviso precisa
// aparecer nos dois casos.
func TestUnsetProjectsOnTenantChangeWithOnlyIamScope(t *testing.T) {
	cfg := newTestConfig(t)
	if err := cfg.Set(mgcConfigPkg.IamProjectKey, "i"); err != nil {
		t.Fatal(err)
	}
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	had, err := unsetProjectsForTenantChange(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !had {
		t.Error("só o escopo do IAM setado ainda deveria disparar o aviso")
	}
	if got := cfg.IamProject(); got != "" {
		t.Errorf("iamProject deveria estar limpo, veio %q", got)
	}
}

// Sem projeto nenhum não há o que avisar, e limpar não pode dar erro — o
// setTenant chama isso em toda troca.
func TestUnsetProjectsOnTenantChangeWithNothingSet(t *testing.T) {
	cfg := newTestConfig(t)
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	had, err := unsetProjectsForTenantChange(ctx)
	if err != nil {
		t.Fatalf("limpar config vazia não pode dar erro: %v", err)
	}
	if had {
		t.Error("sem projeto setado não deveria avisar nada")
	}
}

// Contexto sem config é erro de programação, não algo a engolir em silêncio:
// engolir deixaria o usuário no tenant novo com o projeto do antigo.
func TestUnsetProjectsOnTenantChangeWithoutConfig(t *testing.T) {
	if _, err := unsetProjectsForTenantChange(context.Background()); err == nil {
		t.Error("contexto sem config deveria falhar")
	}
}
