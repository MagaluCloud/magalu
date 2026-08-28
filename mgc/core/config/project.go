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

const (
	// ProjectKey é o projeto da CLI: vai como 'x-project-id' em toda request,
	// exceto nas do IAM.
	ProjectKey = "project"
	// IamProjectKey é o escopo usado SOMENTE pelos endpoints do IAM. É separado
	// de propósito: uma escrita de IAM sem projeto setado vale para o tenant
	// inteiro, então o usuário precisa optar por esse escopo explicitamente, e
	// não herdá-lo do projeto da CLI.
	IamProjectKey = "iamProject"
)

// Project é o projeto da CLI, ou "" quando não há nenhum configurado (nesse
// caso o header é omitido e a API decide o escopo).
func (c *Config) Project() string {
	return c.getString(ProjectKey)
}

// IamProject é o escopo de projeto dos comandos do IAM, ou "" quando não há
// nenhum — o que significa "o tenant inteiro".
func (c *Config) IamProject() string {
	return c.getString(IamProjectKey)
}

// UnsetProject limpa só o escopo da CLI ('mgc project unset').
func (c *Config) UnsetProject() error {
	return c.unsetKeys(ProjectKey)
}

// UnsetIamProject limpa só o escopo do IAM ('mgc iam project unset').
func (c *Config) UnsetIamProject() error {
	return c.unsetKeys(IamProjectKey)
}

// UnsetProjects limpa os dois escopos. Chamado ao trocar de tenant: o projeto
// pertence ao tenant, então continuar com ele apontaria para um escopo de outra
// conta. Mesma razão pela qual a API key é limpa em 'auth tenant set'.
func (c *Config) UnsetProjects() error {
	return c.unsetKeys(ProjectKey, IamProjectKey)
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
