package config

import "context"

// O escopo efetivo da request viaja no contexto, não é deduzido da URL. Quem
// monta a request sabe de que produto ela é e qual chave de config vale; o
// transport só carimba o que recebeu. É o único canal que chega intacto até o
// RoundTripper, porque as requests são criadas com NewRequestWithContext.
type projectScopeKey struct{}

// NewProjectScopeContext carrega o id de projeto já resolvido para esta request.
// String vazia significa "sem escopo": o header não é enviado.
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
// O que continua separado é o COMPORTAMENTO, não a chave: `core.ProjectScope`
// diz se omitir significa "o projeto default" (produtos) ou "o tenant inteiro"
// (IAM).
const ProjectKey = "project"

// Project é o projeto da CLI, ou "" quando não há nenhum configurado (nesse
// caso o header é omitido e a API decide o escopo).
func (c *Config) Project() string {
	return c.getString(ProjectKey)
}

// UnsetProject limpa o projeto selecionado. Chamado por 'mgc project default' e
// ao trocar de tenant: o projeto pertence ao tenant, então mantê-lo apontaria
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
