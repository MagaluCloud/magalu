package sdk

import (
	"net/http"
	"regexp"
	"strings"
)

var _ http.RoundTripper = (*projectTransport)(nil)

// projectHeaderName é o header de escopo de projeto lido pelas APIs da MGC.
const projectHeaderName = "x-project-id"

// iamPathPattern casa as rotas do IAM: o prefixo de serviço declarado na spec
// (`servers: [{url: "https://{env}/iam"}]`) seguido do versionamento que as 36
// rotas publicadas usam. Ancorado e com barra final para não casar `/xiam/...`;
// `v\d+` cobre um futuro /api/v2.
var iamPathPattern = regexp.MustCompile(`^/iam/api/v\d+/`)

// objectStorageHostSuffix identifica o object storage, único serviço em que o
// primeiro segmento do path é string controlada pelo usuário (o nome do
// bucket) — um bucket chamado "iam" não pode ser lido como rota do IAM.
// Casar por host, e não por uma allowlist de hosts do IAM, mantém o
// MGC_SERVERURL apontando para um fake local funcionando.
const objectStorageHostSuffix = "magaluobjects.com"

// projectTransport injeta o escopo de projeto em toda request, no molde do
// DefaultSdkTransport (que já faz isso com User-Agent e X-Request-Id) — assim
// nenhum comando precisa saber que o header existe.
//
// Há DOIS escopos, e a escolha é pela rota: o IAM lê do seu próprio, porque uma
// escrita de IAM sem projeto vale para o tenant inteiro e esse alcance não pode
// ser herdado sem querer do projeto da CLI.
//
// Os valores são lidos por request (não capturados na construção) porque a flag
// global --project-id é aplicada depois que o cliente HTTP já existe.
type projectTransport struct {
	transport  http.RoundTripper
	project    func() string
	iamProject func() string
}

func newProjectTransport(transport http.RoundTripper, project, iamProject func() string) *projectTransport {
	return &projectTransport{transport: transport, project: project, iamProject: iamProject}
}

func (t *projectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Ausência de projeto = header omitido; a API decide o escopo. Nunca vazio.
	if id := t.projectFor(req); id != "" {
		req.Header.Set(projectHeaderName, id)
	}

	transport := t.transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	return transport.RoundTrip(req)
}

func (t *projectTransport) projectFor(req *http.Request) string {
	if isIamRequest(req) {
		return t.iamProject()
	}
	return t.project()
}

// isIamRequest exige host que não seja de object storage E path no formato das
// rotas do IAM. As duas checagens cobrem casos distintos: o host pega um bucket
// "iam" com objeto sob "api/v1/" no OS real, e o padrão pega um bucket "iam"
// com qualquer outro objeto num OS apontado para fake.
func isIamRequest(req *http.Request) bool {
	if req.URL == nil || isObjectStorageHost(req.URL.Hostname()) {
		return false
	}
	return iamPathPattern.MatchString(req.URL.Path)
}

func isObjectStorageHost(host string) bool {
	return host == objectStorageHostSuffix || strings.HasSuffix(host, "."+objectStorageHostSuffix)
}
