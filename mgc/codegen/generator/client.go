package generator

import (
	_ "embed"
	"path"
	"text/template"

	mgcSdkPkg "github.com/MagaluCloud/magalu/mgc/sdk"
)

type clientTemplateData struct {
	PackageName string
	ModuleName  string
}

var (
	//go:embed client.go.template
	clientTemplateContents string
	clientTemplate         *template.Template
)

func init() {
	clientTemplate = templateMust("client.go.template", clientTemplateContents)
}

func generateClient(dirname string, sdk *mgcSdkPkg.Sdk, ctx *GeneratorContext) (err error) {
	return templateWrite(
		ctx,
		path.Join(dirname, "client.go"),
		clientTemplate,
		clientTemplateData{
			PackageName: "client",
			ModuleName:  ctx.ModuleName,
		},
	)
}
