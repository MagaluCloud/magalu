package generator

import (
	_ "embed"
	"fmt"

	"path"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	mgcSdkPkg "magalu.cloud/sdk"
)

type groupTemplateData struct {
	ModuleName    string
	PackageName   string
	PackageImport string
	RefPath       core.RefPath
	core.DescriptorSpec
}

type serviceTemplateData struct {
	ModuleName    string
	PackageName   string
	PackageImport string
	RefPath       core.RefPath
	ClientImport  string
	core.DescriptorSpec
	ExecutorsData []executorTemplateData
}

var (
	//go:embed group.go.template
	groupTemplateContents string
	//go:embed service.go.template
	serviceTemplateContents string
	groupTemplate           *template.Template
	serviceTemplate         *template.Template
)

func init() {
	groupTemplate = templateMust("group.go.template", groupTemplateContents)
	serviceTemplate = templateMust("service.go.template", serviceTemplateContents)
}

func getGroupNames(name string) (fileName string, goName string) {
	name = strings.ReplaceAll(name, " ", "_")
	fileName = strcase.SnakeCase(name)
	goName = strcase.LowerCamelCase(name)
	return
}

var step11 = []string{"auth", "block_storage", "config", "container_registry", "dbaas", "kubernetes", "network"}
var step1 = []string{"virtual-machine-xaas"}
var step13 = []string{"object-storage", "profile", "virtual-machine"}

func generateGroup(dirname string, relPath string, refPath core.RefPath, group core.Grouper, ctx *GeneratorContext) (err error) {
	can_continue := refPath == ""

	for _, step := range step1 {
		if can_continue {
			break
		}
		xpp := strings.Split(string(refPath), "/")
		can_continue = xpp[1] == step
		if can_continue {
			break
		}
	}

	if can_continue {
		fmt.Println(refPath)

		groupDirName, groupGoName := getGroupNames(group.Name())
		p := path.Join(dirname, groupDirName)
		err = createDir(ctx, p)
		if err != nil {
			return
		}

		groupTemplateData := &groupTemplateData{
			ModuleName:     ctx.ModuleName,
			PackageName:    groupGoName,
			PackageImport:  path.Join(ctx.ModuleName, relPath, groupDirName),
			DescriptorSpec: group.DescriptorSpec(),
		}

		serviceTemplateData := &serviceTemplateData{
			ModuleName:     ctx.ModuleName,
			PackageName:    groupGoName,
			PackageImport:  path.Join(ctx.ModuleName, relPath, groupDirName),
			DescriptorSpec: group.DescriptorSpec(),
			ExecutorsData:  []executorTemplateData{},
		}

		err = templateWrite(
			ctx,
			path.Join(p, "doc.go"),
			groupTemplate,
			groupTemplateData,
		)
		if err != nil {
			return
		}

		childRelPath := path.Join(relPath, groupDirName)
		_, err = group.VisitChildren(func(child core.Descriptor) (run bool, err error) {
			if child.IsInternal() {
				return true, nil
			}

			childRefPath := refPath.Add(child.Name())
			switch c := child.(type) {
			case core.Grouper:
				err = generateGroup(p, childRelPath, childRefPath, c, ctx)
				if err != nil {
					return false, &utils.ChainedError{Name: child.Name(), Err: err}
				}
				return true, nil

			case core.Executor:
				execData, err := generateExecutor(p, groupTemplateData, childRefPath, c, ctx)
				if err != nil {
					return false, &utils.ChainedError{Name: child.Name(), Err: err}
				}
				serviceTemplateData.ClientImport = ctx.ModuleName
				serviceTemplateData.ExecutorsData = append(serviceTemplateData.ExecutorsData, execData)
				return true, nil

			default:
				return false, &utils.ChainedError{Name: child.Name(), Err: fmt.Errorf("child %v not group/executor", child)}
			}
		})

		if len(serviceTemplateData.ExecutorsData) > 0 {
			err = templateWrite(
				ctx,
				path.Join(p, "service.go"),
				serviceTemplate,
				serviceTemplateData,
			)
			if err != nil {
				return
			}
		}
	}
	return
}

func generateGroups(dirname string, sdk *mgcSdkPkg.Sdk, ctx *GeneratorContext) (err error) {
	root := sdk.Group()
	err = generateGroup(dirname, "", core.RefPath(""), root, ctx)
	if err != nil {
		return &utils.ChainedError{Name: root.Name(), Err: err}
	}
	return
}
