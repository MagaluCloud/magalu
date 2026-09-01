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
// de outra conta.
func TestUnsetProjectOnTenantChange(t *testing.T) {
	cfg := newTestConfig(t)
	if err := cfg.Set(mgcConfigPkg.ProjectKey, "p"); err != nil {
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
}

// Sem projeto nenhum não há o que avisar, e limpar não pode dar erro — o
// setTenant chama isso em toda troca. `default` conta como "nenhum": avisar
// que ele foi limpo seria ruído em toda troca de quem nunca escolheu projeto.
func TestUnsetProjectOnTenantChangeWithNothingSet(t *testing.T) {
	for name, stored := range map[string]string{"chave ausente": "", "literal default": mgcConfigPkg.ProjectDefault} {
		t.Run(name, func(t *testing.T) {
			cfg := newTestConfig(t)
			if stored != "" {
				if err := cfg.Set(mgcConfigPkg.ProjectKey, stored); err != nil {
					t.Fatal(err)
				}
			}
			ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

			had, err := unsetProjectsForTenantChange(ctx)
			if err != nil {
				t.Fatalf("limpar não pode dar erro: %v", err)
			}
			if had {
				t.Error("sem projeto específico não deveria avisar nada")
			}
		})
	}
}

// Contexto sem config é erro de programação, não algo a engolir em silêncio:
// engolir deixaria o usuário no tenant novo com o projeto do antigo.
func TestUnsetProjectOnTenantChangeWithoutConfig(t *testing.T) {
	if _, err := unsetProjectsForTenantChange(context.Background()); err == nil {
		t.Error("contexto sem config deveria falhar")
	}
}
