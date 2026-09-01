package sdk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
)

const projectHeader = "X-Project-Id"

// captureTransport guarda os headers da última request, sem sair para a rede.
type captureTransport struct {
	last http.Header
}

func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.last = req.Header.Clone()
	return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Request: req}, nil
}

// headersFor roda uma request pelo transport, com ou sem escopo no contexto.
func headersFor(t *testing.T, scope, url string) http.Header {
	t.Helper()
	capture := &captureTransport{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	if scope != "" {
		req = req.WithContext(mgcConfigPkg.NewProjectScopeContext(req.Context(), scope))
	}
	if _, err := newProjectTransport(capture).RoundTrip(req); err != nil {
		t.Fatal(err)
	}
	return capture.last
}

// Sem escopo no contexto o header NÃO é enviado — nunca vazio. É o que faz um
// produto não escopável não receber nada.
func TestNoScopeOmitsHeader(t *testing.T) {
	h := headersFor(t, "", "https://api.magalu.cloud/compute/v1/instances")

	if v, ok := h[projectHeader]; ok {
		t.Errorf("header não deveria existir sem escopo, veio %q", v)
	}
}

// Com escopo no contexto, carimba — e não olha a URL para decidir.
func TestScopeFromContextIsSent(t *testing.T) {
	h := headersFor(t, "proj-alpha", "https://api.magalu.cloud/compute/v1/instances")

	if got, want := h.Get(projectHeader), "proj-alpha"; got != want {
		t.Errorf("header = %q, quer %q", got, want)
	}
}

// A URL deixou de participar da decisão: a MESMA rota carimba ou não conforme
// o contexto. É o que apaga o regex de IAM, a exclusão de object storage e a
// colisão do bucket chamado "iam".
func TestURLDoesNotAffectDecision(t *testing.T) {
	cases := []struct {
		name string
		url  string
	}{
		{"iam", "https://api.magalu.cloud/iam/api/v1/roles"},
		{"object storage", "https://br-se1.magaluobjects.com/iam/objeto.txt"},
		{"bucket api/v1", "https://br-se1.magaluobjects.com/iam/api/v1/x"},
		{"fake local", "http://127.0.0.1:8080/v0/vpcs"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := headersFor(t, "proj-x", c.url).Get(projectHeader); got != "proj-x" {
				t.Errorf("com escopo, %s deveria carimbar proj-x, veio %q", c.url, got)
			}
			if v, ok := headersFor(t, "", c.url)[projectHeader]; ok {
				t.Errorf("sem escopo, %s não deveria carimbar, veio %q", c.url, v)
			}
		})
	}
}

// O transport encadeia: o header chega ao servidor e a resposta de baixo passa
// intacta.
func TestProjectTransportDelegates(t *testing.T) {
	var seen string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get(projectHeader)
		w.WriteHeader(http.StatusTeapot)
	}))
	t.Cleanup(srv.Close)

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/qualquer", nil)
	req = req.WithContext(mgcConfigPkg.NewProjectScopeContext(req.Context(), "proj-alpha"))

	resp, err := (&http.Client{Transport: newProjectTransport(http.DefaultTransport)}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("status = %d, quer 418 (resposta de baixo preservada)", resp.StatusCode)
	}
	if seen != "proj-alpha" {
		t.Errorf("servidor recebeu %q, quer proj-alpha", seen)
	}
}
