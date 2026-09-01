package config

import (
	"context"
	"strings"
)

// O escopo efetivo da request viaja no contexto, não é deduzido da URL. Quem
// monta a request sabe de que produto ela é e qual chave de config vale; o
// transport só carimba o que recebeu. É o único canal que chega intacto até o
// RoundTripper, porque as requests são criadas com NewRequestWithContext.
type projectScopeKey struct{}

// NewProjectScopeContext carrega o id de projeto já resolvido para esta request.
// String vazia significa "sem escopo": o header não é enviado.
//
// Na prática quem resolve nunca devolve vazio para produto escopável: sem nada
// selecionado vem ProjectDefault, que a API entende igual à ausência.
func NewProjectScopeContext(parent context.Context, projectID string) context.Context {
	return context.WithValue(parent, projectScopeKey{}, projectID)
}

// ProjectScopeFromContext devolve o id resolvido, ou "" quando não há escopo.
func ProjectScopeFromContext(ctx context.Context) string {
	id, _ := ctx.Value(projectScopeKey{}).(string)
	return id
}

// ProjectKey é o projeto selecionado, usado por TODO produto escopável — IAM
// incluído. Houve uma segunda chave (`iamProject`) enquanto o IAM tinha o seu
// próprio comando de seleção; com `mgc iam project` removido, ninguém a
// gravava, e `mgc project set` não chegava aos comandos de IAM.
//
// Também houve dois COMPORTAMENTOS, `core.ProjectScope` "project" e "iam", pela
// mesma razão. Quem diz o alcance de um comando de IAM hoje é a --parent-type
// dele, então sobrou um booleano: DescriptorSpec.ProjectScoped.
const ProjectKey = "project"

// ProjectDefault é o projeto default do tenant, escrito por extenso. A API o
// aceita no lugar de um uuid, então mandá-lo é equivalente a omitir o escopo —
// e gravá-lo é melhor que gravar vazio: `mgc project current` responde onde
// você está em vez de responder nada, e a configuração diz o que faz em vez de
// depender de quem a lê saber o que a ausência significa.
const ProjectDefault = "default"

// IsProjectDefault diz se o valor representa o projeto default — vazio incluído,
// porque configuração nova e `auth tenant set` deixam a chave ausente.
func IsProjectDefault(id string) bool {
	return id == "" || strings.EqualFold(id, ProjectDefault)
}

// Project é o valor cru da configuração, ou "" quando não há nenhum. Quem
// precisa do escopo efetivo deve tratar vazio e ProjectDefault como o mesmo
// caso — use IsProjectDefault.
func (c *Config) Project() string {
	return c.getString(ProjectKey)
}

// UnsetProject remove a chave. Não é mais o caminho de 'mgc project default',
// que grava ProjectDefault; sobrou para a troca de tenant: o projeto pertence ao tenant, então mantê-lo apontaria
// para um escopo de outra conta. Mesma razão pela qual a API key é limpa em
// 'auth tenant set'.
func (c *Config) UnsetProject() error {
	return c.unsetKeys(ProjectKey)
}

func (c *Config) unsetKeys(keys ...string) error {
	for _, key := range keys {
		if c.getString(key) == "" {
			continue // já limpo: nada a remover do arquivo
		}
		if err := c.Delete(key); err != nil {
			return err
		}
		// Delete reescreve o arquivo, mas não desfaz o override em memória que
		// Config.Set deixa no viper (viper.Set grava na camada de maior
		// precedência, que sobrevive ao ReadConfig). Sem isto, uma leitura no
		// mesmo processo ainda enxergaria o projeto antigo. O temp config é
		// consultado antes do viper, então mascara o resíduo.
		if err := c.SetTempConfig(key, ""); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) getString(key string) string {
	var s string
	if err := c.Get(key, &s); err != nil {
		return ""
	}
	return s
}
