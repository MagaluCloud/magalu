package config

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
