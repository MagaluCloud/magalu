package project

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
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

// Nome duplicado não pode escolher um "vencedor" arbitrário: erra pedindo o id.
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

// Se um nome coincidir com o id de OUTRO projeto, o id vence — é o
// identificador não-ambíguo.
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

// Tenant sem projeto: a mensagem tem de dizer o que fazer.
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
	// Sem nada selecionado a resposta é o literal, não vazio: o comando existe
	// para dizer onde você está, e "em lugar nenhum" não é resposta.
	if got.ID != mgcConfigPkg.ProjectDefault {
		t.Errorf("sem projeto o id deveria ser %q, veio %q", mgcConfigPkg.ProjectDefault, got.ID)
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

// O default precisa ser confirmável e a mensagem não pode ser vazia: mensagem
// vazia faz o handleExecutor pular o prompt silenciosamente.
func TestDefaultIsConfirmableAndExplainsFallback(t *testing.T) {
	cExec, ok := core.ExecutorAs[core.ConfirmableExecutor](getDefault())
	if !ok {
		t.Fatal("default deveria implementar ConfirmableExecutor")
	}

	msg := cExec.ConfirmPrompt(core.Parameters{}, core.Configs{})
	if msg == "" {
		t.Fatal("mensagem vazia pularia o prompt")
	}
	for _, want := range []string{"default project", "tenant"} {
		if !strings.Contains(msg, want) {
			t.Errorf("confirmação deveria mencionar %q, veio: %q", want, msg)
		}
	}
}

// Trocar de escopo é reversível; só a remoção pede confirmação.
func TestSetAndCurrentAreNotConfirmable(t *testing.T) {
	for name, exec := range map[string]core.Executor{"set": getSet(), "current": getCurrent()} {
		if _, ok := core.ExecutorAs[core.ConfirmableExecutor](exec); ok {
			t.Errorf("%s não deveria pedir confirmação", name)
		}
	}
}

// default GRAVA o literal, não apaga a chave. Apagar deixava a configuração
// muda: quem a lesse depois precisava saber que a ausência significa "o projeto
// default".
func TestDefaultWritesTheLiteral(t *testing.T) {
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := mgcConfigPkg.New(pm)
	if err := cfg.Set(mgcConfigPkg.ProjectKey, "p"); err != nil {
		t.Fatal(err)
	}
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	if _, err := useDefault(ctx, struct{}{}, struct{}{}); err != nil {
		t.Fatalf("useDefault: %v", err)
	}
	if got, want := cfg.Project(), mgcConfigPkg.ProjectDefault; got != want {
		t.Errorf("project = %q, quer %q", got, want)
	}
}

// `default` NÃO é palavra reservada no positional do set: o projeto que a
// pessoa criou com esse nome tem de ser selecionável pelo nome, como qualquer
// outro. Era o contrário antes, e o projeto real ficava inalcançável em
// silêncio — o comando selecionava o default do tenant sem dizer nada.
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

// Sem projeto com esse nome — o caso comum — o erro tem de apontar a saída, ou
// a pessoa fica sabendo só que não achou.
func TestSetDefaultWithoutSuchProjectPointsToTheCommand(t *testing.T) {
	_, err := resolveProject([]projectResult{{ID: "id-alpha", Name: "alpha"}}, "default")
	if err == nil {
		t.Fatal("sem projeto chamado default deveria falhar")
	}
	if !strings.Contains(err.Error(), "mgc project default") {
		t.Errorf("o erro deveria ensinar o comando, veio: %v", err)
	}
}

// --- current enriquecido: id da config + nome vindo da lista ---
//
// Nao ha GET /{id} na API de projects, entao "resolver o nome" e listar e
// procurar. A decisao de o que devolver em cada desfecho fica isolada aqui,
// testavel sem rede; currentScope so faz a ligacao.

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

// Config apontando para projeto que nao existe mais (apagado, ou de outro
// tenant) e informacao util — hoje isso so aparece como 403 em outro comando.
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

// Nunca falhar: sem rede ou deslogado, o comando ainda responde o essencial —
// o id, que e local. O nome e enriquecimento.
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

// Sem projeto configurado nao ha o que resolver — e nao se gasta uma requisicao.
func TestCurrentWithoutProjectDoesNotLookUp(t *testing.T) {
	pm, _ := profile_manager.NewInMemoryProfileManager()
	cfg := mgcConfigPkg.New(pm)
	// Contexto SEM auth: se currentScope tentasse listar, viria erro.
	ctx := mgcConfigPkg.NewContext(context.Background(), cfg)

	got, err := currentScope(ctx, mgcConfigPkg.ProjectKey, projectConfig{})
	if err != nil {
		t.Fatalf("sem projeto nao e erro: %v", err)
	}
	if got.ID != mgcConfigPkg.ProjectDefault || got.Warning != "" {
		t.Errorf("deveria devolver o literal sem aviso, veio %+v", got)
	}
}

// Com projeto setado e sem auth no contexto, a consulta falha — e mesmo assim
// o comando responde, com aviso.
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
