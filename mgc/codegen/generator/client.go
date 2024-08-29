package generator

import (
	_ "embed"
	"path"
	"text/template"
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

func generateClient(dirname string, ctx *GeneratorContext) (err error) {
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
