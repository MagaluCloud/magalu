package generator

import (
	_ "embed"
	"path"
	"text/template"

	mgcSdkPkg "magalu.cloud/sdk"
)

type helpersTemplateData struct {
	PackageName  string
	ModuleName   string
	ClientImport string
}

const (
	helpersPackage = "helpers"
)

var (
	//go:embed helpers_converters.go.template
	helpersConvertersTemplateContents string
	helpersConvertersTemplate         *template.Template

	//go:embed helpers_executors.go.template
	helpersExecutorsTemplateContents string
	helpersExecutorsTemplate         *template.Template
)

func init() {
	helpersConvertersTemplate = templateMust("helpers_converters.go.template", helpersConvertersTemplateContents)
	helpersExecutorsTemplate = templateMust("helpers_executors.go.template", helpersExecutorsTemplateContents)
}

func generateHelpers(dirname string, sdk *mgcSdkPkg.Sdk, ctx *GeneratorContext) (err error) {
	p := path.Join(dirname, helpersPackage)
	err = createDir(ctx, p)
	if err != nil {
		return err
	}

	data := helpersTemplateData{
		PackageName:  helpersPackage,
		ModuleName:   ctx.ModuleName,
		ClientImport: ctx.ModuleName,
	}

	err = templateWrite(
		ctx,
		path.Join(p, "converters.go"),
		helpersConvertersTemplate,
		data,
	)
	if err != nil {
		return
	}

	return templateWrite(
		ctx,
		path.Join(p, "executors.go"),
		helpersExecutorsTemplate,
		data,
	)
}
