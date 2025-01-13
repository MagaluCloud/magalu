package cmd

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/MagaluCloud/magalu/mgc/codegen/generator"
	"github.com/MagaluCloud/magalu/mgc/codegen/ui"
	mgcSdkPkg "github.com/MagaluCloud/magalu/mgc/sdk"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewRoot() (rootCmd *cobra.Command) {
	rootCmd = &cobra.Command{
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"outputDir"},
		Use:       os.Args[0],
		Short:     "Outputs a human-friendly SDK to access the internal runtime",
		RunE:      run,
	}
	addVerboseFlag(rootCmd)
	addModuleNameFlag(rootCmd)

	return
}

func run(cmd *cobra.Command, args []string) (err error) {
	outputDir, err := filepath.Abs(args[0])
	if err != nil {
		return
	}
	sdk := mgcSdkPkg.NewSdk()

	var out io.Writer
	if getVerboseFlag(cmd) {
		out = os.Stdout
	}
	ctx := &generator.GeneratorContext{
		ModuleName: getModuleNameFlag(cmd),
		Reporter:   ui.NewReporter(path.Dir(outputDir), out, os.Stderr),
		FS:         afero.NewOsFs(),
	}

	return generator.GenerateSdk(outputDir, sdk, ctx)
}
