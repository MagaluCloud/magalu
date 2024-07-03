package generator

import (
	_ "embed"
	"path"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type executorTemplateData struct {
	ModuleName          string
	PackageName         string
	PackageImport       string
	ClientImport        string
	HelpersImport       string
	GoName              string
	TerminatorExecutor  bool
	ConfirmableExecutor bool
	RefPath             core.RefPath
	Types               *templateTypes
	core.DescriptorSpec
}

type templateTypes struct {
	Parameters string
	Configs    string
	Result     string
	Types      generatorTemplateTypes
}

var (
	//go:embed executor.go.template
	executorTemplateContents string
	executorTemplate         *template.Template
)

func init() {
	executorTemplate = templateMust("executor.go.template", executorTemplateContents)
}

func getExecutorNames(name string) (fileName string, goName string) {
	name = strings.ReplaceAll(name, " ", "_")
	fileName = strcase.SnakeCase(name)
	goName = strcase.UpperCamelCase(name)
	return
}

func generateTypes(exec core.Executor, execGoName string) (t *templateTypes, err error) {
	t = &templateTypes{}

	if schema := exec.ParametersSchema(); schema != nil && len(schema.Properties) > 0 {
		t.Parameters = execGoName + "Parameters"
		_, err = t.Types.addSchemaOrAlias(t.Parameters, schema)
		if err != nil {
			err = &utils.ChainedError{Name: "parameters", Err: err}
			return
		}
	}

	if schema := exec.ConfigsSchema(); schema != nil && len(schema.Properties) > 0 {
		t.Configs = execGoName + "Configs"
		_, err = t.Types.addSchemaOrAlias(t.Configs, schema)
		if err != nil {
			err = &utils.ChainedError{Name: "configs", Err: err}
			return
		}
	}

	if schema := exec.ResultSchema(); schema != nil && !mgcSchemaPkg.CheckSimilarJsonSchemas(schema, nullSchema) {
		t.Result = execGoName + "Result"
		_, err = t.Types.addSchemaOrAlias(t.Result, schema)
		if err != nil {
			err = &utils.ChainedError{Name: "result", Err: err}
			return
		}
	}

	return
}

func generateExecutor(dirname string, groupTemplateData *groupTemplateData, refPath core.RefPath, exec core.Executor, ctx *GeneratorContext) (executorTemplateData, error) {
	execDirName, execGoName := getExecutorNames(exec.Name())
	_, isTerminatorExecutor := exec.(core.TerminatorExecutor)
	_, isConfirmableExecutor := exec.(core.ConfirmableExecutor)
	types, err := generateTypes(exec, execGoName)
	if err != nil {
		return executorTemplateData{}, err
	}

	execData := executorTemplateData{
		ModuleName:          ctx.ModuleName,
		PackageName:         groupTemplateData.PackageName,
		PackageImport:       groupTemplateData.PackageImport,
		ClientImport:        ctx.ModuleName,
		HelpersImport:       path.Join(ctx.ModuleName, helpersPackage),
		GoName:              execGoName,
		TerminatorExecutor:  isTerminatorExecutor,
		ConfirmableExecutor: isConfirmableExecutor,
		RefPath:             refPath,
		Types:               types,
		DescriptorSpec:      exec.DescriptorSpec(),
	}

	err = templateWrite(
		ctx,
		path.Join(dirname, execDirName+".go"),
		executorTemplate,
		execData,
	)

	return execData, err
}
