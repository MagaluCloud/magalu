package cmd

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
)

func newScopeTestConfig(t *testing.T) *config.Config {
	t.Helper()
	pm, _ := profile_manager.NewInMemoryProfileManager()
	return config.New(pm)
}

// Produto que não declara escopo não recebe header nenhum — é o caso de 6 dos
// 10 produtos, e o motivo de a declaração ser opt-in.
func TestResolveScopeUndeclaredIsEmpty(t *testing.T) {
	cfg := newScopeTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "proj-1"); err != nil {
		t.Fatal(err)
	}

	if got := resolveProjectScope(cfg, ""); got != "" {
		t.Errorf("produto sem declaração não pode carimbar, veio %q", got)
	}
}

// Há uma chave só: `mgc project set` vale para produtos E para o IAM. Antes
// eram duas, e o projeto selecionado simplesmente não chegava aos comandos de
// IAM — que liam uma chave que nenhum comando gravava.
func TestResolveScopeReadsTheSingleKey(t *testing.T) {
	cfg := newScopeTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "proj-1"); err != nil {
		t.Fatal(err)
	}

	for _, scope := range []core.ProjectScope{core.ProjectScopeProduct, core.ProjectScopeIAM} {
		if got, want := resolveProjectScope(cfg, scope), "proj-1"; got != want {
			t.Errorf("escopo %q = %q, quer %q", scope, got, want)
		}
	}
}

// A flag vence a configuração: --project-id grava no temp config, que Config.Get
// consulta antes do viper.
func TestResolveScopeFlagWinsOverConfig(t *testing.T) {
	cfg := newScopeTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "do-arquivo"); err != nil {
		t.Fatal(err)
	}
	if err := cfg.SetTempConfig(config.ProjectKey, "da-flag"); err != nil {
		t.Fatal(err)
	}

	if got, want := resolveProjectScope(cfg, core.ProjectScopeIAM), "da-flag"; got != want {
		t.Errorf("escopo = %q, quer %q", got, want)
	}
}

// Sem projeto selecionado não há escopo — e isso não é erro. Para produtos
// significa "o projeto default"; para o IAM, quem diz o alcance é a
// --parent-type, não este valor.
func TestResolveScopeUnsetIsEmpty(t *testing.T) {
	if got := resolveProjectScope(newScopeTestConfig(t), core.ProjectScopeIAM); got != "" {
		t.Errorf("sem projeto configurado o escopo é vazio, veio %q", got)
	}
}
