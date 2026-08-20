package config

import "testing"

// Ausência é o caso normal: sem projeto configurado o header não é enviado, e
// a API decide o escopo (projeto default do tenant).
func TestProjectAbsentIsEmpty(t *testing.T) {
	c, _ := setupWithoutFile("")

	if got := c.Project(); got != "" {
		t.Errorf("Project() sem config = %q, quer vazio", got)
	}
	if got := c.IamProject(); got != "" {
		t.Errorf("IamProject() sem config = %q, quer vazio", got)
	}
}

// Os dois escopos são independentes: setar um não afeta o outro. É o ponto do
// desenho — o projeto do IAM é escopo de configuração, não o da CLI.
func TestProjectAndIamProjectAreIndependent(t *testing.T) {
	c, _ := setupWithoutFile("")

	if err := c.Set(ProjectKey, "proj-cli"); err != nil {
		t.Fatal(err)
	}
	if got, want := c.Project(), "proj-cli"; got != want {
		t.Errorf("Project() = %q, quer %q", got, want)
	}
	if got := c.IamProject(); got != "" {
		t.Errorf("setar project não pode setar iamProject, veio %q", got)
	}

	if err := c.Set(IamProjectKey, "proj-iam"); err != nil {
		t.Fatal(err)
	}
	if got, want := c.IamProject(), "proj-iam"; got != want {
		t.Errorf("IamProject() = %q, quer %q", got, want)
	}
	if got, want := c.Project(), "proj-cli"; got != want {
		t.Errorf("setar iamProject não pode alterar project, veio %q", got)
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

func TestIamProjectEnv(t *testing.T) {
	t.Setenv("MGC_IAMPROJECT", "iam-do-env")
	c, _ := setupWithoutFile("")

	if got, want := c.IamProject(), "iam-do-env"; got != want {
		t.Errorf("MGC_IAMPROJECT = %q, quer %q", got, want)
	}
}

// Trocar de tenant invalida os dois escopos (o projeto pertence ao tenant), do
// mesmo jeito que a API key é limpa hoje pelo `auth tenant set`.
func TestUnsetProjectsClearsBoth(t *testing.T) {
	c, _ := setupWithoutFile("")

	if err := c.Set(ProjectKey, "p"); err != nil {
		t.Fatal(err)
	}
	if err := c.Set(IamProjectKey, "i"); err != nil {
		t.Fatal(err)
	}

	if err := c.UnsetProjects(); err != nil {
		t.Fatalf("UnsetProjects: %v", err)
	}
	if got := c.Project(); got != "" {
		t.Errorf("project deveria estar limpo, veio %q", got)
	}
	if got := c.IamProject(); got != "" {
		t.Errorf("iamProject deveria estar limpo, veio %q", got)
	}
}

// Limpar o que nunca foi setado não pode ser erro — `auth tenant set` chama isso
// sempre, tenha o usuário escolhido projeto ou não.
func TestUnsetProjectsIsIdempotent(t *testing.T) {
	c, _ := setupWithoutFile("")

	if err := c.UnsetProjects(); err != nil {
		t.Fatalf("UnsetProjects em config vazia deveria ser no-op, veio: %v", err)
	}
}

// `mgc project unset` limpa só o escopo da CLI; o do IAM tem comando próprio.
func TestUnsetProjectKeepsIamProject(t *testing.T) {
	c, _ := setupWithoutFile("")
	if err := c.Set(ProjectKey, "p"); err != nil {
		t.Fatal(err)
	}
	if err := c.Set(IamProjectKey, "i"); err != nil {
		t.Fatal(err)
	}

	if err := c.UnsetProject(); err != nil {
		t.Fatalf("UnsetProject: %v", err)
	}

	if got := c.Project(); got != "" {
		t.Errorf("project deveria estar limpo, veio %q", got)
	}
	if got, want := c.IamProject(), "i"; got != want {
		t.Errorf("iamProject não podia ser tocado: %q, quer %q", got, want)
	}
}

// E o inverso, para `mgc iam project unset`.
func TestUnsetIamProjectKeepsProject(t *testing.T) {
	c, _ := setupWithoutFile("")
	if err := c.Set(ProjectKey, "p"); err != nil {
		t.Fatal(err)
	}
	if err := c.Set(IamProjectKey, "i"); err != nil {
		t.Fatal(err)
	}

	if err := c.UnsetIamProject(); err != nil {
		t.Fatalf("UnsetIamProject: %v", err)
	}

	if got := c.IamProject(); got != "" {
		t.Errorf("iamProject deveria estar limpo, veio %q", got)
	}
	if got, want := c.Project(), "p"; got != want {
		t.Errorf("project não podia ser tocado: %q, quer %q", got, want)
	}
}
