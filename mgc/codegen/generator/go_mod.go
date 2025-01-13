package generator

import (
	_ "embed"
	"path"
	"strings"
	"text/template"

	mgcSdkPkg "github.com/MagaluCloud/magalu/mgc/sdk"
)

type modTemplateData struct {
	ModuleName  string
	CoreVersion string
	SdkVersion  string
}

var (
	//go:embed go.mod.template
	modTemplateContents string
	modTemplate         *template.Template
)

func init() {
	modTemplate = templateMust("go.mod.template", modTemplateContents)
}

func getVersion(rawVersion string) string {
	parts := strings.Split(rawVersion, " ")
	return parts[0]
}

func generateGoMod(dirname string, sdk *mgcSdkPkg.Sdk, ctx *GeneratorContext) (err error) {
	return templateWrite(
		ctx,
		path.Join(dirname, "go.mod"),
		modTemplate,
		modTemplateData{
			ModuleName:  ctx.ModuleName,
			CoreVersion: getVersion(mgcSdkPkg.Version), // assume the same
			SdkVersion:  getVersion(mgcSdkPkg.Version),
		},
	)
}
