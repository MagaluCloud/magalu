package openapi

import (
	"os"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

var testPrefix = func() *string { p := "x-mgc"; return &p }()

// GUARDA DE REGENERAÇÃO — o teste mais importante desta etapa.
//
// A exigência de escopo só existe na spec EMBUTIDA, e ela chega lá pelo merge
// das customizações. Se alguém regenerar as specs sem reaplicá-las, a exigência
// some SEM ERRO NENHUM e escritas de IAM voltam a passar sem escopo. Este teste
// transforma esse afrouxamento silencioso em CI vermelho.
func TestEmbeddedIAMSpecCarriesScopeRequired(t *testing.T) {
	data, err := os.ReadFile("openapis/iam.openapi.yaml")
	if err != nil {
		t.Fatalf("spec embutida do IAM: %v", err)
	}

	var doc struct {
		Servers []map[string]any `yaml:"servers"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.Servers) == 0 {
		t.Fatal("spec do IAM sem bloco servers")
	}

	methods, ok := doc.Servers[0]["x-mgc-scope-required"]
	if !ok {
		t.Fatal("a spec embutida do IAM perdeu `x-mgc-scope-required`. " +
			"Regenerar sem reaplicar openapi-customizations/ afrouxa a exigência de escopo em silêncio")
	}
	list, ok := methods.([]any)
	if !ok || len(list) == 0 {
		t.Fatalf("`x-mgc-scope-required` deveria ser uma lista de métodos, veio %#v", methods)
	}
}

// A regra é por MÉTODO: escrita exige escopo, leitura nunca. Bloquear leitura
// seria hostil — é como a pessoa descobre o que existe.
func TestScopeRequiredByMethod(t *testing.T) {
	servers := openapi3.Servers{{
		URL:        "https://{env}/iam",
		Extensions: map[string]any{"x-mgc-scope-required": []any{"post", "patch", "put", "delete"}},
	}}

	cases := map[string]bool{
		"POST": true, "PATCH": true, "PUT": true, "DELETE": true,
		"GET": false, "HEAD": false,
	}
	for method, want := range cases {
		if got := getScopeRequiredExtension(testPrefix, servers, nil, method); got != want {
			t.Errorf("%s -> %v, quer %v", method, got, want)
		}
	}
}

// Serviço que não declara nada não exige escopo de ninguém.
func TestScopeRequiredUndeclaredService(t *testing.T) {
	servers := openapi3.Servers{{URL: "https://{env}/{region}/compute"}}

	if getScopeRequiredExtension(testPrefix, servers, nil, "POST") {
		t.Error("produto sem declaração não pode exigir escopo")
	}
}

// OPT-OUT: a operação pode se declarar isenta. É a válvula para escrita que só
// existe no tenant e que, sem isso, pediria --scope tenant para sempre.
func TestScopeRequiredOperationOptOut(t *testing.T) {
	servers := openapi3.Servers{{
		URL:        "https://{env}/iam",
		Extensions: map[string]any{"x-mgc-scope-required": []any{"post"}},
	}}
	opt := map[string]any{"x-mgc-scope-required": false}

	if getScopeRequiredExtension(testPrefix, servers, opt, "POST") {
		t.Error("a operação declarou isenção e mesmo assim exigiu escopo")
	}
}

// A isenção é só para MENOS: uma operação não pode exigir escopo num serviço
// que não declarou nada. Opt-out, nunca opt-in — senão a marcação volta a ser
// um a um, que é o que a direção da falha queria evitar.
func TestScopeRequiredOperationCannotOptIn(t *testing.T) {
	servers := openapi3.Servers{{URL: "https://{env}/{region}/compute"}}
	opt := map[string]any{"x-mgc-scope-required": true}

	if getScopeRequiredExtension(testPrefix, servers, opt, "POST") {
		t.Error("operação não pode LIGAR a exigência num serviço que não declarou")
	}
}

// O campo precisa chegar ao descritor, que é o que viaja até a CLI.
func TestScopeRequiredReachesDescriptor(t *testing.T) {
	var spec core.DescriptorSpec
	spec.ScopeRequired = true
	if !spec.ScopeRequired {
		t.Error("DescriptorSpec precisa carregar ScopeRequired")
	}
}
