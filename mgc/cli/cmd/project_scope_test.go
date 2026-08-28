package cmd

import (
	"strings"
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
	if err := cfg.Set(config.ProjectKey, "proj-cli"); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Set(config.ProjectKey, "proj-iam"); err != nil {
		t.Fatal(err)
	}

	if got := resolveProjectScope(cfg, "", ""); got != "" {
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
		if got, want := resolveProjectScope(cfg, scope, ""), "proj-1"; got != want {
			t.Errorf("escopo %q = %q, quer %q", scope, got, want)
		}
	}
}

// `--scope tenant` omite o header: é assim que o IAM entende "o tenant inteiro".
func TestResolveScopeTenantOmitsHeader(t *testing.T) {
	cfg := newScopeTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "proj-iam"); err != nil {
		t.Fatal(err)
	}

	if got := resolveProjectScope(cfg, core.ProjectScopeIAM, scopeTenant); got != "" {
		t.Errorf("--scope tenant deveria omitir, veio %q", got)
	}
}

// `--scope default` manda o id do TENANT: é a codificação que o IAM usa para o
// projeto default. Quirk da API, escondido do usuário.
func TestResolveScopeDefaultSendsTenantID(t *testing.T) {
	cfg := newScopeTestConfig(t)

	got := resolveProjectScopeWithTenant(cfg, core.ProjectScopeIAM, scopeDefault, "tenant-123")
	if got != "tenant-123" {
		t.Errorf("--scope default = %q, quer o id do tenant", got)
	}
}

// A flag vence a configuração, pelos dois lados.
func TestResolveScopeFlagWinsOverConfig(t *testing.T) {
	cfg := newScopeTestConfig(t)
	if err := cfg.Set(config.ProjectKey, "do-arquivo"); err != nil {
		t.Fatal(err)
	}
	if err := cfg.SetTempConfig(config.ProjectKey, "da-flag"); err != nil {
		t.Fatal(err)
	}

	if got, want := resolveProjectScope(cfg, core.ProjectScopeIAM, ""), "da-flag"; got != want {
		t.Errorf("escopo = %q, quer %q", got, want)
	}
}

// --- exigência de escopo nas escritas de IAM ---

// O caso que a feature existe para evitar: escrita de IAM sem escopo nenhum.
// Omitir o header, ali, significa "o tenant inteiro".
func TestScopeRequiredMissingIsReported(t *testing.T) {
	cfg := newScopeTestConfig(t)

	msg := checkScopeRequired(cfg, true, core.ProjectScopeIAM, "")
	if msg == "" {
		t.Fatal("escrita de IAM sem escopo deveria ser reportada")
	}
	// A mensagem tem de ensinar as três saídas, senão o usuário fica preso.
	for _, want := range []string{"--scope default", "--scope tenant", "--project-id"} {
		if !strings.Contains(msg, want) {
			t.Errorf("mensagem deveria citar %q, veio: %q", want, msg)
		}
	}
}

// Qualquer uma das três formas satisfaz.
func TestScopeRequiredSatisfied(t *testing.T) {
	cfgComProjeto := newScopeTestConfig(t)
	if err := cfgComProjeto.Set(config.ProjectKey, "proj-iam"); err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		cfg  *config.Config
		flag string
	}{
		"projeto configurado": {cfgComProjeto, ""},
		"--scope tenant":      {newScopeTestConfig(t), scopeTenant},
		"--scope default":     {newScopeTestConfig(t), scopeDefault},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if msg := checkScopeRequired(c.cfg, true, core.ProjectScopeIAM, c.flag); msg != "" {
				t.Errorf("não deveria reportar nada, veio: %q", msg)
			}
		})
	}
}

// Leitura nunca exige: bloquear `iam roles list` seria hostil, é como a pessoa
// descobre o que existe.
func TestScopeNotRequiredForReads(t *testing.T) {
	if msg := checkScopeRequired(newScopeTestConfig(t), false, core.ProjectScopeIAM, ""); msg != "" {
		t.Errorf("leitura não pode exigir escopo, veio: %q", msg)
	}
}

// Produto que não é do IAM nunca exige, mesmo em escrita: para eles omitir
// significa "projeto default", que é seguro.
func TestScopeNotRequiredForProducts(t *testing.T) {
	if msg := checkScopeRequired(newScopeTestConfig(t), true, core.ProjectScopeProduct, ""); msg != "" {
		t.Errorf("produto comum não pode exigir escopo, veio: %q", msg)
	}
}
