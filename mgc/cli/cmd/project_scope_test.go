package cmd

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
)

func newScopeTestConfig(t *testing.T) *config.Config {
	t.Helper()
	pm, _ := profile_manager.NewInMemoryProfileManager()
	return config.New(pm)
}

// Sem nada selecionado o escopo é o literal, não vazio: é o que faz
// `mgc project current` responder onde você está, e a configuração dizer o que
// faz em vez de depender de quem a lê saber o que a ausência significa.
func TestResolveScopeUnsetIsTheLiteralDefault(t *testing.T) {
	if got, want := resolveProjectScope(newScopeTestConfig(t)), config.ProjectDefault; got != want {
		t.Errorf("sem projeto configurado = %q, quer %q", got, want)
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

	if got, want := resolveProjectScope(cfg), "proj-1"; got != want {
		t.Errorf("escopo = %q, quer %q", got, want)
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

	if got, want := resolveProjectScope(cfg), "da-flag"; got != want {
		t.Errorf("escopo = %q, quer %q", got, want)
	}
}

// `--project-id default` é aceito e sobrevive intacto: mandá-lo é equivalente a
// omitir, e quem traduz isso para cada produto é quem monta a request.
func TestResolveScopeAcceptsTheLiteralDefault(t *testing.T) {
	cfg := newScopeTestConfig(t)
	if err := cfg.SetTempConfig(config.ProjectKey, config.ProjectDefault); err != nil {
		t.Fatal(err)
	}

	if got, want := resolveProjectScope(cfg), config.ProjectDefault; got != want {
		t.Errorf("escopo = %q, quer %q", got, want)
	}
}
