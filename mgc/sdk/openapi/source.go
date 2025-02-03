package openapi

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/dataloader"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/invopop/yaml"
)

type indexModuleSpec struct {
	core.DescriptorSpec
	Url  string
	Path string
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
			modules = make([]core.Grouper, len(index.Modules))
			refResolver := core.NewMultiRefPathResolver()
			for i := range index.Modules {
				var module core.Grouper
				module, err = newModule(
					&index.Modules[i],
					extensionPrefix,
					loader,
					logger(),
					refResolver,
				)
				if err != nil {
					err = &utils.ChainedError{Name: index.Modules[i].Path, Err: err}
					return
				}
				modules[i] = module
			}

			return modules, nil
		},
	)
}
