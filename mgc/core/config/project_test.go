package config

import "testing"

func TestProjectAbsentIsEmpty(t *testing.T) {
	c, _ := setupWithoutFile("")

	if got := c.Project(); got != "" {
		t.Errorf("Project() sem config = %q, quer vazio", got)
	}
}

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

func TestUnsetProjectIsIdempotent(t *testing.T) {
	c, _ := setupWithoutFile("")

	if err := c.UnsetProject(); err != nil {
		t.Fatalf("UnsetProject em config vazia deveria ser no-op, veio: %v", err)
	}
}
