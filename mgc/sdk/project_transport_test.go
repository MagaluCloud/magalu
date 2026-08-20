package sdk

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

// headersFor roda uma request pelo transport com os dois escopos fixos e
// devolve os headers que chegariam ao servidor.
func headersFor(t *testing.T, project, iamProject, url string) http.Header {
	t.Helper()
	capture := &captureTransport{}
	transport := newProjectTransport(capture,
		func() string { return project },
		func() string { return iamProject },
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := transport.RoundTrip(req); err != nil {
		t.Fatal(err)
	}
	return capture.last
}

// Sem projeto configurado o header NÃO é enviado — nunca vazio. É o contrato:
// ausência significa "a API decide o escopo".
func TestNoProjectOmitsHeader(t *testing.T) {
	h := headersFor(t, "", "", "https://api.magalu.cloud/compute/v1/instances")

	if v, ok := h[projectHeader]; ok {
		t.Errorf("header não deveria existir sem projeto, veio %q", v)
	}
}

// Endpoint comum usa o projeto da CLI.
func TestNonIamUsesProject(t *testing.T) {
	h := headersFor(t, "proj-cli", "proj-iam", "https://api.magalu.cloud/compute/v1/instances")

	if got, want := h.Get(projectHeader), "proj-cli"; got != want {
		t.Errorf("header = %q, quer %q", got, want)
	}
}

// Endpoint do IAM usa o OUTRO escopo. É o coração do desenho: o projeto do IAM
// é escopo de configuração e não pode ser herdado do projeto da CLI.
func TestIamUsesIamProject(t *testing.T) {
	h := headersFor(t, "proj-cli", "proj-iam", "https://api.magalu.cloud/iam/api/v1/projects")

	if got, want := h.Get(projectHeader), "proj-iam"; got != want {
		t.Errorf("header no IAM = %q, quer %q", got, want)
	}
}

// IAM sem escopo próprio não cai no projeto da CLI: o header some, e a ação
// vale para o tenant inteiro (que é justamente o caso que pede confirmação).
func TestIamWithoutIamProjectOmitsHeaderEvenWithCliProject(t *testing.T) {
	h := headersFor(t, "proj-cli", "", "https://api.magalu.cloud/iam/api/v1/roles")

	if v, ok := h[projectHeader]; ok {
		t.Errorf("IAM não pode herdar o projeto da CLI, veio %q", v)
	}
}

// O inverso: endpoint comum não herda o escopo do IAM.
func TestNonIamWithoutProjectOmitsHeaderEvenWithIamProject(t *testing.T) {
	h := headersFor(t, "", "proj-iam", "https://api.magalu.cloud/compute/v1/instances")

	if v, ok := h[projectHeader]; ok {
		t.Errorf("endpoint comum não pode herdar o escopo do IAM, veio %q", v)
	}
}

// O reconhecimento do IAM é pela FORMA da rota (/iam/api/v<N>/...), não por
// "contém iam". Cobre prod, pre-prod e fake local, e recusa o que só se parece.
func TestIamDetectionMatchesRouteShape(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want string
	}{
		{"prod", "https://api.magalu.cloud/iam/api/v1/roles", "proj-iam"},
		{"pre-prod", "https://api.pre-prod.jaxyendy.com:8443/iam/api/v1/projects", "proj-iam"},
		{"fake local", "http://127.0.0.1:8080/iam/api/v1/projects", "proj-iam"},
		{"versão futura", "http://127.0.0.1:8080/iam/api/v2/projects", "proj-iam"},

		{"bucket iam no OS real", "https://br-se1.magaluobjects.com/iam/objeto.txt", "proj-cli"},
		{"bucket iam com api/v1 no OS real", "https://br-se1.magaluobjects.com/iam/api/v1/x", "proj-cli"},
		{"bucket iam em OS apontado p/ fake", "http://127.0.0.1:8080/iam/objeto.txt", "proj-cli"},
		{"iam no meio do path", "https://api.magalu.cloud/compute/v1/iam", "proj-cli"},
		{"prefixo parecido", "https://api.magalu.cloud/xiam/api/v1/roles", "proj-cli"},
		{"versão não-numérica", "https://api.magalu.cloud/iam/api/vX/roles", "proj-cli"},
		{"sem o segmento de versão", "https://api.magalu.cloud/iam/roles", "proj-cli"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := headersFor(t, "proj-cli", "proj-iam", c.url)
			if got := h.Get(projectHeader); got != c.want {
				t.Errorf("%s -> header %q, quer %q", c.url, got, c.want)
			}
		})
	}
}

// Limitação conhecida e aceita: object storage apontado para um fake, com bucket
// chamado "iam" e objeto sob "api/v<N>/", é indistinguível de uma rota do IAM.
// Exige as três coincidências ao mesmo tempo e o header é ignorado pelo S3.
// Este teste existe para que a escolha fique registrada — se algum dia ela pesar,
// o caminho é a identidade vir do contexto de quem monta a request, não da URL.
func TestIamDetectionKnownLimitation(t *testing.T) {
	h := headersFor(t, "proj-cli", "proj-iam", "http://127.0.0.1:8080/iam/api/v1/objeto.txt")

	if got := h.Get(projectHeader); got != "proj-iam" {
		t.Errorf("comportamento mudou: %q. Se foi de propósito, atualize o comentário", got)
	}
}

// O transport é montado uma vez, mas o valor é lido a cada request: nada de
// congelar o projeto na construção do cliente (a flag é aplicada depois).
func TestProjectIsReadPerRequest(t *testing.T) {
	current := "primeiro"
	capture := &captureTransport{}
	transport := newProjectTransport(capture,
		func() string { return current },
		func() string { return "" },
	)

	send := func() string {
		req, _ := http.NewRequest(http.MethodGet, "https://api.magalu.cloud/compute/v1/instances", nil)
		if _, err := transport.RoundTrip(req); err != nil {
			t.Fatal(err)
		}
		return capture.last.Get(projectHeader)
	}

	if got := send(); got != "primeiro" {
		t.Fatalf("header = %q, quer primeiro", got)
	}
	current = "segundo"
	if got := send(); got != "segundo" {
		t.Fatalf("header = %q, quer segundo (lido por request)", got)
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

	client := &http.Client{Transport: newProjectTransport(http.DefaultTransport,
		func() string { return "proj-cli" },
		func() string { return "" },
	)}
	resp, err := client.Get(srv.URL + "/compute/v1/instances")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("status = %d, quer 418 (resposta de baixo preservada)", resp.StatusCode)
	}
	if seen != "proj-cli" {
		t.Errorf("servidor recebeu %q, quer proj-cli", seen)
	}
}
