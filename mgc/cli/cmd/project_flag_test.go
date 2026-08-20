package cmd

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	"github.com/spf13/cobra"
)

func newFlagTestCmd(t *testing.T) *cobra.Command {
	t.Helper()
	root := &cobra.Command{Use: "mgc"}
	addProjectFlags(root)
	return root
}

func newTestConfig(t *testing.T) *config.Config {
	t.Helper()
	pm, _ := profile_manager.NewInMemoryProfileManager()
	return config.New(pm)
}

// Sem as flags, nada é gravado no temp config — o valor do arquivo/env segue
// valendo.
func TestProjectFlagsAbsentDoNotOverride(t *testing.T) {
	root := newFlagTestCmd(t)
	cfg := newTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "do-arquivo"); err != nil {
		t.Fatal(err)
	}

	applyProjectFlags(root, cfg)

	if got, want := cfg.Project(), "do-arquivo"; got != want {
		t.Errorf("Project() = %q, quer %q (flag ausente não sobrescreve)", got, want)
	}
}

// A flag global vence o arquivo, igual ao --api-key.
func TestProjectFlagOverridesConfig(t *testing.T) {
	root := newFlagTestCmd(t)
	cfg := newTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "do-arquivo"); err != nil {
		t.Fatal(err)
	}
	if err := root.PersistentFlags().Set(projectFlag, "da-flag"); err != nil {
		t.Fatal(err)
	}

	applyProjectFlags(root, cfg)

	if got, want := cfg.Project(), "da-flag"; got != want {
		t.Errorf("Project() = %q, quer %q", got, want)
	}
}

// O escopo do IAM tem a sua própria flag e não é afetado pela outra.
func TestIamProjectFlagIsIndependent(t *testing.T) {
	root := newFlagTestCmd(t)
	cfg := newTestConfig(t)
	if err := root.PersistentFlags().Set(projectFlag, "da-flag"); err != nil {
		t.Fatal(err)
	}

	applyProjectFlags(root, cfg)

	if got := cfg.IamProject(); got != "" {
		t.Errorf("--project-id não pode setar o escopo do IAM, veio %q", got)
	}

	if err := root.PersistentFlags().Set(iamProjectFlag, "iam-da-flag"); err != nil {
		t.Fatal(err)
	}
	applyProjectFlags(root, cfg)

	if got, want := cfg.IamProject(), "iam-da-flag"; got != want {
		t.Errorf("IamProject() = %q, quer %q", got, want)
	}
	if got, want := cfg.Project(), "da-flag"; got != want {
		t.Errorf("--iam-project-id não pode alterar o projeto da CLI, veio %q", got)
	}
}

// A flag vence o env, que já vence o arquivo (o env é lido pelo viper).
func TestProjectFlagOverridesEnv(t *testing.T) {
	t.Setenv("MGC_PROJECT", "do-env")
	root := newFlagTestCmd(t)
	cfg := newTestConfig(t)

	applyProjectFlags(root, cfg)
	if got, want := cfg.Project(), "do-env"; got != want {
		t.Fatalf("sem flag deveria valer o env: %q, quer %q", got, want)
	}

	if err := root.PersistentFlags().Set(projectFlag, "da-flag"); err != nil {
		t.Fatal(err)
	}
	applyProjectFlags(root, cfg)

	if got, want := cfg.Project(), "da-flag"; got != want {
		t.Errorf("flag deveria vencer o env: %q, quer %q", got, want)
	}
}
