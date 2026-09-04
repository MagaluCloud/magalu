package config

import "testing"

// Ausência é o caso normal: sem projeto configurado o header não é enviado, e
// a API decide o escopo (projeto default do tenant).
func TestProjectAbsentIsEmpty(t *testing.T) {
	c, _ := setupWithoutFile("")

	if got := c.Project(); got != "" {
		t.Errorf("Project() sem config = %q, quer vazio", got)
	}
}

// A flag global sobrescreve o arquivo: é o mesmo mecanismo do --api-key, que
// grava um temp config consultado antes do viper.
func TestProjectTempConfigOverridesFile(t *testing.T) {
	c, _ := setupWithoutFile("")

	if err := c.Set(ProjectKey, "do-arquivo"); err != nil {
		t.Fatal(err)
	}
	if err := c.SetTempConfig(ProjectKey, "da-flag"); err != nil {
		t.Fatal(err)
	}

	if got, want := c.Project(), "da-flag"; got != want {
		t.Errorf("flag deveria vencer o arquivo: %q, quer %q", got, want)
	}
}

// O env é lido pelo viper (AutomaticEnv + prefixo MGC), então MGC_PROJECT vale
// sem nenhum wiring próprio — e a flag ainda vence o env.
func TestProjectEnvIsReadAndFlagStillWins(t *testing.T) {
	t.Setenv("MGC_PROJECT", "do-env")
	c, _ := setupWithoutFile("")

	if got, want := c.Project(), "do-env"; got != want {
		t.Errorf("MGC_PROJECT = %q, quer %q", got, want)
	}

	if err := c.SetTempConfig(ProjectKey, "da-flag"); err != nil {
		t.Fatal(err)
	}
	if got, want := c.Project(), "da-flag"; got != want {
		t.Errorf("flag deveria vencer o env: %q, quer %q", got, want)
	}
}

// Limpar o que nunca foi setado não pode ser erro — `auth tenant set` chama isso
// sempre, tenha o usuário escolhido projeto ou não.
func TestUnsetProjectIsIdempotent(t *testing.T) {
	c, _ := setupWithoutFile("")

	if err := c.UnsetProject(); err != nil {
		t.Fatalf("UnsetProject em config vazia deveria ser no-op, veio: %v", err)
	}
}
