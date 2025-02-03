package cmd

import (
	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const watchFlag = "cli.watch"

func newWatchFlag() *flag.Flag {
	watchFlagSchema := mgcSchemaPkg.NewBooleanSchema()
	watchFlagSchema.Description = `Wait until the operation is completed by calling the 'get' link and waiting until termination. Akin to '! get -w'`

	flag := schema_flags.NewSchemaFlag(
		mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			watchFlag: watchFlagSchema,
		}, nil),
		watchFlag,
		watchFlag,
		false,
		false,
		false,
	)
	flag.NoOptDefVal = "true"
	return flag
}

func getWatchFlag(cmd *cobra.Command) bool {
	w, err := cmd.PersistentFlags().GetBool(watchFlag)
	if err != nil {
		return false
	}
	return w
}
