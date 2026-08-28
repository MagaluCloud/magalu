package cmd

import (
	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const scopeFlag = "scope"

// newScopeFlag monta a flag --scope. Ela NÃO é persistente do root: entra por
// comando via addExtraFlag, do mesmo jeito que --cli.watch entra só onde há link
// `get` terminador. Assim ela não polui o help dos produtos que não a entendem,
// e usá-la onde não existe já erra com "unknown flag", sem validação nossa.
//
// Só o IAM a recebe: os demais produtos não têm nível de tenant, e para eles
// "default" já é a ausência de configuração — um enum de um valor só.
func newScopeFlag() *flag.Flag {
	schema := mgcSchemaPkg.NewStringSchema()
	schema.Description = "Scope this command applies to: 'default' for the tenant's default project, 'tenant' for the entire tenant. The IAM API encodes the default project as the tenant id"
	schema.Enum = []any{scopeDefault, scopeTenant}

	return schema_flags.NewSchemaFlag(
		mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{scopeFlag: schema}, nil),
		scopeFlag,
		scopeFlag,
		false,
		false,
		false,
	)
}

// hasScopeFlag diz se este executor deve receber --scope.
func hasScopeFlag(exec core.Executor) bool {
	return exec.DescriptorSpec().ProjectScope == core.ProjectScopeIAM
}

func getScopeFlag(cmd *cobra.Command) string {
	f := cmd.Flags().Lookup(scopeFlag)
	if f == nil {
		return ""
	}
	return f.Value.String()
}
