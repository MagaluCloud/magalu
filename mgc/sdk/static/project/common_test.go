package project

import (
	"context"
	"errors"
	"strings"
	"testing"

	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
)

func TestHostForEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		expected string
	}{
		{name: "empty falls back to prod", env: "", expected: prodHost},
		{name: "default falls back to prod", env: "default", expected: prodHost},
		{name: "prod", env: "prod", expected: prodHost},
		{name: "pre-prod", env: "pre-prod", expected: preProdHost},
		{name: "raw pre-prod host", env: preProdHost, expected: preProdHost},
		{name: "unknown falls back to prod", env: "staging", expected: prodHost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hostForEnv(tt.env); got != tt.expected {
				t.Errorf("hostForEnv(%q) = %q, expected %q", tt.env, got, tt.expected)
			}
		})
	}
}

func TestBuildProjectsURL(t *testing.T) {
	tests := []struct {
		name      string
		serverUrl string
		env       string
		expected  string
	}{
		{
			name:     "prod default",
			env:      "",
			expected: "https://api.magalu.cloud/iam/api/v1/projects",
		},
		{
			name:     "pre-prod",
			env:      "pre-prod",
			expected: "https://api.pre-prod.jaxyendy.com:8443/iam/api/v1/projects",
		},
		{
			name:      "serverUrl replaces the whole base, including the /iam prefix",
			serverUrl: "http://localhost:8080",
			env:       "prod",
			expected:  "http://localhost:8080/api/v1/projects",
		},
		{
			name:      "serverUrl trailing slash",
			serverUrl: "http://localhost:8080/",
			expected:  "http://localhost:8080/api/v1/projects",
		},
		{
			name:      "serverUrl wins over env",
			serverUrl: "https://example.com/iam",
			env:       "pre-prod",
			expected:  "https://example.com/iam/api/v1/projects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildProjectsURL(tt.serverUrl, tt.env)
			if got != tt.expected {
				t.Errorf("buildProjectsURL(%q, %q) = %q, expected %q", tt.serverUrl, tt.env, got, tt.expected)
			}
		})
	}
}

func TestResolveProjectByIDAndName(t *testing.T) {
	available := []projectResult{
		{ID: "id-alpha", Name: "alpha", Type: "default"},
		{ID: "id-beta", Name: "Beta", Type: "managed"},
	}

	cases := []struct {
		name    string
		query   string
		wantID  string
		wantErr string
	}{
		{name: "por id", query: "id-beta", wantID: "id-beta"},
		{name: "por nome", query: "alpha", wantID: "id-alpha"},
		{name: "nome ignora caixa", query: "BETA", wantID: "id-beta"},
		{name: "desconhecido", query: "gama", wantErr: "not found"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := resolveProject(available, c.query)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("erro = %v, quer conter %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveProject(%q): %v", c.query, err)
			}
			if got.ID != c.wantID {
				t.Errorf("id = %q, quer %q", got.ID, c.wantID)
			}
		})
	}
}

func TestResolveProjectAmbiguousName(t *testing.T) {
	available := []projectResult{
		{ID: "id-1", Name: "repetido"},
		{ID: "id-2", Name: "Repetido"},
	}

	_, err := resolveProject(available, "repetido")
	if err == nil {
		t.Fatal("nome ambíguo deveria falhar")
	}
	for _, want := range []string{"id-1", "id-2"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("erro deveria listar %q para desempate, veio: %v", want, err)
		}
	}
}

func TestResolveProjectIDWinsOverName(t *testing.T) {
	available := []projectResult{
		{ID: "colisao", Name: "alpha"},
		{ID: "id-beta", Name: "colisao"},
	}

	got, err := resolveProject(available, "colisao")
	if err != nil {
		t.Fatalf("resolveProject: %v", err)
	}
	if got.ID != "colisao" {
		t.Errorf("id = %q, quer colisao (match por id vence)", got.ID)
	}
}

func TestResolveProjectEmptyList(t *testing.T) {
	_, err := resolveProject(nil, "alpha")
	if err == nil {
		t.Fatal("lista vazia deveria falhar")
	}
	if !strings.Contains(err.Error(), "project create") {
		t.Errorf("erro deveria sugerir criar um projeto, veio: %v", err)
	}
}

func TestCurrentReadsConfig(t *testing.T) {
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := mgcConfigPkg.New(pm)
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	got, err := current(ctx, struct{}{}, projectConfig{})
	if err != nil {
		t.Fatalf("current sem projeto setado não é erro: %v", err)
	}
	if got.ID != "" || got.Warning != warnNoProject {
		t.Errorf("sem projeto deveria avisar, veio %+v", got)
	}

	if err := cfg.Set(mgcConfigPkg.ProjectKey, "id-alpha"); err != nil {
		t.Fatal(err)
	}
	got, err = current(ctx, struct{}{}, projectConfig{})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "id-alpha" {
		t.Errorf("id = %q, quer id-alpha", got.ID)
	}
}

func TestSetResolvesAProjectNamedDefault(t *testing.T) {
	available := []projectResult{
		{ID: "id-do-projeto", Name: "default", Type: "managed"},
		{ID: "id-alpha", Name: "alpha", Type: "default"},
	}

	got, err := resolveProject(available, "default")
	if err != nil {
		t.Fatalf("resolveProject: %v", err)
	}
	if got.ID != "id-do-projeto" {
		t.Errorf("id = %q, quer id-do-projeto (o projeto real, não o sentinela)", got.ID)
	}
}

func TestResolveCurrentFound(t *testing.T) {
	available := []projectResult{
		{ID: "id-alpha", Name: "alpha", Type: "default"},
		{ID: "id-beta", Name: "beta", Type: "managed"},
	}

	got := resolveCurrent("id-beta", available, nil)
	if got.ID != "id-beta" || got.Name != "beta" || got.Type != "managed" {
		t.Errorf("resolveCurrent = %+v, quer id-beta/beta/managed", got)
	}
	if got.Warning != "" {
		t.Errorf("achou o projeto, nao deveria avisar nada: %q", got.Warning)
	}
}

func TestResolveCurrentNotInTenant(t *testing.T) {
	got := resolveCurrent("id-fantasma", []projectResult{{ID: "id-alpha", Name: "alpha"}}, nil)

	if got.ID != "id-fantasma" {
		t.Errorf("o id da config tem de sobreviver: %q", got.ID)
	}
	if got.Name != "" {
		t.Errorf("sem correspondencia nao ha nome: %q", got.Name)
	}
	if !strings.Contains(got.Warning, "not found") {
		t.Errorf("deveria avisar que nao encontrou: %q", got.Warning)
	}
}

func TestResolveCurrentLookupFailed(t *testing.T) {
	got := resolveCurrent("id-alpha", nil, errors.New("no network"))

	if got.ID != "id-alpha" {
		t.Errorf("o id tem de sobreviver a falha de consulta: %q", got.ID)
	}
	if got.Warning == "" {
		t.Error("degradar em silencio esconderia que a consulta falhou")
	}
	if strings.Contains(got.Warning, "not found") {
		t.Errorf("falha de consulta nao e 'nao encontrado': %q", got.Warning)
	}
}

func TestCurrentWithoutProjectDoesNotLookUp(t *testing.T) {
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := mgcConfigPkg.New(pm)
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	got, err := currentScope(ctx, mgcConfigPkg.ProjectKey, projectConfig{})
	if err != nil {
		t.Fatalf("sem projeto nao e erro: %v", err)
	}
	if got.ID != "" || got.Warning != warnNoProject {
		t.Errorf("deveria avisar que nao ha projeto, veio %+v", got)
	}
}

func TestCurrentDegradesWithoutAuth(t *testing.T) {
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := mgcConfigPkg.New(pm)
	if err := cfg.Set(mgcConfigPkg.ProjectKey, "id-alpha"); err != nil {
		t.Fatal(err)
	}
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	got, err := currentScope(ctx, mgcConfigPkg.ProjectKey, projectConfig{})
	if err != nil {
		t.Fatalf("nunca falhar: %v", err)
	}
	if got.ID != "id-alpha" {
		t.Errorf("id = %q, quer id-alpha", got.ID)
	}
	if got.Warning == "" {
		t.Error("deveria avisar que nao conseguiu resolver o nome")
	}
}

func TestSetWithoutLoginAcceptsUUID(t *testing.T) {
	t.Setenv("MGC_PROJECT", "")
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := mgcConfigPkg.New(pm)
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	got, err := set(ctx, setParams{IDOrName: "3F2A9C1E-5B7D-4E8A-9C0F-1A2B3C4D5E6F"}, projectConfig{})
	if err != nil {
		t.Fatalf("set sem login com UUID: %v", err)
	}

	want := "3f2a9c1e-5b7d-4e8a-9c0f-1a2b3c4d5e6f"
	if got.ID != want || cfg.Project() != want {
		t.Errorf("id = %q, config = %q, quer %q", got.ID, cfg.Project(), want)
	}
	if got.Warning != warnNotVerified {
		t.Errorf("deveria avisar que nao verificou, veio %q", got.Warning)
	}
}

func TestSetWithoutLoginRejectsName(t *testing.T) {
	t.Setenv("MGC_PROJECT", "")
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := mgcConfigPkg.New(pm)
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	_, err := set(ctx, setParams{IDOrName: "my-project"}, projectConfig{})
	if err == nil || !strings.Contains(err.Error(), "UUID") {
		t.Fatalf("nome sem login deveria pedir UUID, veio: %v", err)
	}
	if got := cfg.Project(); got != "" {
		t.Errorf("nada deveria ser gravado, veio %q", got)
	}
}

func TestUnsetClearsProject(t *testing.T) {
	t.Setenv("MGC_PROJECT", "")
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := mgcConfigPkg.New(pm)
	if err := cfg.SetProject("id-alpha"); err != nil {
		t.Fatal(err)
	}
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	for i := 0; i < 2; i++ {
		if _, err := unset(ctx); err != nil {
			t.Fatalf("unset #%d: %v", i+1, err)
		}
		if got := cfg.Project(); got != "" {
			t.Errorf("unset #%d: project = %q, quer vazio", i+1, got)
		}
	}
}
