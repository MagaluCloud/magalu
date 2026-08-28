package config

import (
	"context"
	"testing"
)

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

// --- escopo efetivo viajando no contexto ---
//
// O transport nao decide mais nada: quem sabe qual produto e qual escopo e a
// camada que monta a request. O contexto e o unico canal, porque e o que chega
// intacto ate o RoundTripper (http.NewRequestWithContext).

func TestProjectScopeContextRoundTrip(t *testing.T) {
	ctx := context.Background()

	if got := ProjectScopeFromContext(ctx); got != "" {
		t.Errorf("contexto limpo deveria devolver vazio, veio %q", got)
	}

	ctx = NewProjectScopeContext(ctx, "id-alpha")
	if got := ProjectScopeFromContext(ctx); got != "id-alpha" {
		t.Errorf("escopo = %q, quer id-alpha", got)
	}
}

// Escopo vazio nao pode virar header vazio: quem nao tem escopo nao carimba
// nada, e o contexto tem de reportar isso como ausencia.
func TestProjectScopeContextEmptyStaysEmpty(t *testing.T) {
	ctx := NewProjectScopeContext(context.Background(), "")

	if got := ProjectScopeFromContext(ctx); got != "" {
		t.Errorf("escopo vazio = %q, quer vazio", got)
	}
}

// Escopo mais interno vence: uma sub-request com escopo proprio nao herda o de
// fora sem querer.
func TestProjectScopeContextInnermostWins(t *testing.T) {
	ctx := NewProjectScopeContext(context.Background(), "de-fora")
	ctx = NewProjectScopeContext(ctx, "de-dentro")

	if got := ProjectScopeFromContext(ctx); got != "de-dentro" {
		t.Errorf("escopo = %q, quer de-dentro", got)
	}
}
