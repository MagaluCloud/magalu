package openapi

import (
	"fmt"
	"os"

	"github.com/invopop/yaml"
	"magalu.cloud/core"
	"magalu.cloud/core/dataloader"
	"magalu.cloud/core/utils"
)

type indexModuleSpec struct {
	core.DescriptorSpec
	Url           string
	Path          string
	SuperInternal bool `json:"super_internal,omitempty"`
}

type indexFileSpec struct {
	Version string
	Modules []indexModuleSpec
}

const indexFileName = "index.openapi.yaml"
const indexVersion = "1.0.0"

// Source -> Module -> Resource -> Operation

// -- ROOT: Source

func NewSource(loader dataloader.Loader, extensionPrefix *string) *core.SimpleGrouper[core.Grouper] {
	return core.NewSimpleGrouper(
		core.DescriptorSpec{
			Name:        "OpenApis",
			Description: fmt.Sprintf("OpenApis loaded using %v", loader),
		},
		func() (modules []core.Grouper, err error) {
			data, err := loader.Load(indexFileName)
			if err != nil {
				return nil, err
			}

			var index indexFileSpec
			err = yaml.Unmarshal(data, &index)
			if err != nil {
				return nil, err
			}
			if index.Version != indexVersion {
				return nil, fmt.Errorf("unsupported %q version %q, expected %q", indexFileName, index.Version, indexVersion)
			}

			var finalIndexModule []indexModuleSpec
			if os.Getenv("IGNORE_SUPER_HIDDEN") == "true" {
				finalIndexModule = index.Modules
			} else {
				for _, mod := range index.Modules {
					if !mod.SuperInternal {
						finalIndexModule = append(finalIndexModule, mod)
					}
				}
			}

			modules = make([]core.Grouper, len(finalIndexModule))
			refResolver := core.NewMultiRefPathResolver()
			for i := range finalIndexModule {
				var module core.Grouper
				module, err = newModule(
					&finalIndexModule[i],
					extensionPrefix,
					loader,
					logger(),
					refResolver,
				)
				if err != nil {
					err = &utils.ChainedError{Name: finalIndexModule[i].Path, Err: err}
					return
				}
				modules[i] = module
			}

			return modules, nil
		},
	)
}
