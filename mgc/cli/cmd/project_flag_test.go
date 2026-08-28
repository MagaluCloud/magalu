package cmd

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	"github.com/spf13/cobra"
)

// newFlagTestCmd monta um comando com a --project-id anexada como ela passa a
// ser na prática: por comando, não persistente do root.
func newFlagTestCmd(t *testing.T) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "list"}
	cmd.Flags().AddFlag(newProjectIDFlag())
	return cmd
}

func newTestConfig(t *testing.T) *config.Config {
	t.Helper()
	pm, _ := profile_manager.NewInMemoryProfileManager()
	return config.New(pm)
}

// Sem a flag, nada é gravado no temp config — o valor do arquivo/env segue
// valendo.
func TestProjectFlagsAbsentDoNotOverride(t *testing.T) {
	cmd := newFlagTestCmd(t)
	cfg := newTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "do-arquivo"); err != nil {
		t.Fatal(err)
	}

	applyProjectFlags(cmd, cfg)

	if got, want := cfg.Project(), "do-arquivo"; got != want {
		t.Errorf("Project() = %q, quer %q (flag ausente não sobrescreve)", got, want)
	}
}

// A flag vence o arquivo, igual ao --api-key.
func TestProjectFlagOverridesConfig(t *testing.T) {
	cmd := newFlagTestCmd(t)
	cfg := newTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "do-arquivo"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set(projectFlag, "da-flag"); err != nil {
		t.Fatal(err)
	}

	applyProjectFlags(cmd, cfg)

	if got, want := cfg.Project(), "da-flag"; got != want {
		t.Errorf("Project() = %q, quer %q", got, want)
	}
}

// `--project-id` vale para QUALQUER produto escopável, IAM incluído — é a flag
// única. Gravar só em `project` a deixaria inerte nos comandos de IAM.
func TestProjectFlagReachesBothScopes(t *testing.T) {
	cmd := newFlagTestCmd(t)
	cfg := newTestConfig(t)
	if err := cmd.Flags().Set(projectFlag, "da-flag"); err != nil {
		t.Fatal(err)
	}

	applyProjectFlags(cmd, cfg)

	if got, want := cfg.Project(), "da-flag"; got != want {
		t.Errorf("Project() = %q, quer %q", got, want)
	}
	if got, want := cfg.IamProject(), "da-flag"; got != want {
		t.Errorf("IamProject() = %q, quer %q", got, want)
	}
}

// A flag vence o env, que já vence o arquivo (o env é lido pelo viper).
func TestProjectFlagOverridesEnv(t *testing.T) {
	t.Setenv("MGC_PROJECT", "do-env")
	cmd := newFlagTestCmd(t)
	cfg := newTestConfig(t)

	applyProjectFlags(cmd, cfg)
	if got, want := cfg.Project(), "do-env"; got != want {
		t.Fatalf("sem flag deveria valer o env: %q, quer %q", got, want)
	}

	if err := cmd.Flags().Set(projectFlag, "da-flag"); err != nil {
		t.Fatal(err)
	}
	applyProjectFlags(cmd, cfg)

	if got, want := cfg.Project(), "da-flag"; got != want {
		t.Errorf("flag deveria vencer o env: %q, quer %q", got, want)
	}
}

// Comando de produto que não declara escopo não recebe a flag — e ler uma flag
// inexistente não pode explodir.
func TestProjectFlagAbsentFromUnscopedCommand(t *testing.T) {
	cmd := &cobra.Command{Use: "list"} // sem AddFlag
	cfg := newTestConfig(t)

	applyProjectFlags(cmd, cfg) // não deve entrar em pânico

	if got := cfg.Project(); got != "" {
		t.Errorf("nada deveria ter sido gravado, veio %q", got)
	}
}
