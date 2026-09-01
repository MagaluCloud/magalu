package cmd

import (
	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const projectFlag = "project-id"

// newProjectIDFlag monta a --project-id. Como a --scope, ela NÃO é persistente
// do root: entra por comando via addExtraFlag, só onde o produto declara escopo.
//
// Metade dos produtos não suporta projeto — uma flag global apareceria no help
// e nos docs de todos eles sem fazer nada. E usá-la onde não existe passa a
// errar com "unknown flag", do próprio cobra, sem validação nossa.
func newProjectIDFlag() *flag.Flag {
	schema := mgcSchemaPkg.NewStringSchema()
	schema.Description = "Project to scope this command to, or 'default' for the tenant's default project. Overrides the configured project for this invocation"

	return schema_flags.NewSchemaFlag(
		mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{projectFlag: schema}, nil),
		projectFlag,
		projectFlag,
		false,
		false,
		false,
	)
}

// hasProjectFlag diz se este executor deve receber --project-id: todo produto
// escopável, IAM incluído.
func hasProjectFlag(exec core.Executor) bool {
	return exec.DescriptorSpec().ProjectScoped
}

func getProjectFlag(cmd *cobra.Command, name string) string {
	f := cmd.Flags().Lookup(name)
	if f == nil {
		return ""
	}
	return f.Value.String()
}

// applyProjectFlags dá à flag precedência sobre env e arquivo, gravando no temp
// config — que Config.Get consulta antes do viper. Mesmo mecanismo do --api-key.
//
// Há uma chave só: `--project-id` significa "nesta invocação, use este
// projeto", e vale para todo produto escopável, IAM incluído.
func applyProjectFlags(cmd *cobra.Command, cfg *config.Config) {
	value := getProjectFlag(cmd, projectFlag)
	if value == "" {
		return
	}
	_ = cfg.SetTempConfig(config.ProjectKey, value)
}
