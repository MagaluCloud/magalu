package sdk

import (
	"net/http"

	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
)

var _ http.RoundTripper = (*projectTransport)(nil)

// projectHeaderName é o header de escopo de projeto lido pelas APIs da MGC.
const projectHeaderName = "x-project-id"

// projectTransport carimba o escopo de projeto em toda request que o traga no
// contexto, no molde do DefaultSdkTransport — assim nenhum comando precisa
// saber que o header existe.
//
// Ele NÃO decide nada: não sabe o que é IAM, nem o que é object storage, nem
// lê URL. Quem monta a request sabe de que produto ela é e qual chave de config
// vale, e resolve o id antes; aqui só resta carimbar. Deduzir o produto pela URL
// era heurística — `--server-url` apaga o prefixo do serviço, e no object
// storage o primeiro segmento do path é o nome do bucket.
type projectTransport struct {
	transport http.RoundTripper
}

func newProjectTransport(transport http.RoundTripper) *projectTransport {
	return &projectTransport{transport: transport}
}

func (t *projectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Ausência de escopo = header omitido. Nunca vazio: para os produtos
	// escopáveis omitir significa "projeto default", e para o IAM, "o tenant".
	if id := mgcConfigPkg.ProjectScopeFromContext(req.Context()); id != "" {
		req.Header.Set(projectHeaderName, id)
	}

	transport := t.transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	return transport.RoundTrip(req)
}
