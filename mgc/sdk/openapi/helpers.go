package openapi

import (
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/getkin/kin-openapi/openapi3"
)

func getExtension(prefix *string, name string, extensions map[string]any, def any) (value any, ok bool) {
	if prefix == nil || *prefix == "" {
		return def, false
	}
	key := *prefix + "-" + name
	value, ok = extensions[key]
	if !ok {
		value = def
	}
	return
}

func getExtensionString(prefix *string, name string, extensions map[string]any, def string) (str string, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	str, ok = value.(string)
	return
}

func getExtensionObject(prefix *string, name string, extensions map[string]any, def map[string]any) (m map[string]any, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	m, ok = value.(map[string]any)
	return
}

func getExtensionArray(prefix *string, name string, extensions map[string]any, def []any) (m []any, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	m, ok = value.([]any)
	return
}

func getExtensionBool(prefix *string, name string, extensions map[string]any, def bool) (b bool, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	b, ok = value.(bool)
	return
}

func getNameExtension(prefix *string, extensions map[string]any, def string) string {
	str, _ := getExtensionString(prefix, "name", extensions, def)
	return str
}

func getDescriptionExtension(prefix *string, extensions map[string]any, def string) string {
	str, _ := getExtensionString(prefix, "description", extensions, def)
	return str
}

func getHiddenExtension(prefix *string, extensions map[string]any) bool {
	b, _ := getExtensionBool(prefix, "hidden", extensions, false)
	return b
}

// getProjectScopeExtension lê `x-mgc-project-scope` do bloco `servers` da spec.
// A declaração é por SERVIÇO, não por operação: todo endpoint de um produto
// escopável participa, então marcar um a um seria repetição que envelhece mal —
// endpoint novo nasceria de fora sem ninguém notar.
//
// Valores: "project" (lê a chave `project`) e "iam" (lê `iamProject`). Qualquer
// outro é ignorado, e ausente significa produto não escopável.
func getProjectScopeExtension(prefix *string, servers openapi3.Servers) core.ProjectScope {
	if len(servers) == 0 {
		return ""
	}
	v, _ := getExtensionString(prefix, "project-scope", servers[0].Extensions, "")
	switch v {
	case string(core.ProjectScopeProduct):
		return core.ProjectScopeProduct
	case string(core.ProjectScopeIAM):
		return core.ProjectScopeIAM
	default:
		return ""
	}
}

// getScopeRequiredExtension diz se a operação exige escopo explícito.
//
// A declaração é do SERVIÇO — `x-mgc-scope-required: [post, patch, put, delete]`
// no bloco `servers` — e vale por MÉTODO: escrita exige, leitura nunca. Assim
// endpoint novo do produto nasce protegido, sem ninguém marcar.
//
// A operação pode se declarar isenta (`x-mgc-scope-required: false`), para a
// escrita que só existe no tenant. Mas NÃO pode ligar a exigência: se o serviço
// não declarou nada, operação nenhuma exige. O opt-out é o que faz esquecer uma
// marcação produzir ruído visível em vez de exposição silenciosa.
func getScopeRequiredExtension(prefix *string, servers openapi3.Servers, opExtensions map[string]any, method string) bool {
	if len(servers) == 0 {
		return false
	}

	raw, ok := getExtension(prefix, "scope-required", servers[0].Extensions, nil)
	if !ok || raw == nil {
		return false // serviço não declarou: ninguém exige
	}

	methods, ok := raw.([]any)
	if !ok {
		return false
	}
	required := false
	for _, m := range methods {
		if s, ok := m.(string); ok && strings.EqualFold(s, method) {
			required = true
			break
		}
	}
	if !required {
		return false
	}

	// Isenção da operação: só desliga.
	if optOut, ok := getExtensionBool(prefix, "scope-required", opExtensions, true); ok && !optOut {
		return false
	}
	return true
}
