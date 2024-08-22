package workspace

import (
	"context"
	"errors"

	"magalu.cloud/core"
	"magalu.cloud/core/profile_manager"
	"magalu.cloud/core/utils"
)

type createParams struct {
	Name string `json:"name" jsonschema_description:"Workspace name" mgc:"positional"`
	Copy string `json:"copy,omitempty" jsonschema_description:"Name of the workspace to copy content from. If this parameter is passed, the new Workspace will be pre-populated with the contents of the Workspace with the specified name" mgc:"positional"`
}

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	exec := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "create",
			Description: "Creates a new workspace",
		},
		create,
	)

	return core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "template=Created workspace {{.name}}\n"
	})
})

func create(ctx context.Context, params createParams, _ struct{}) (*profile_manager.Profile, error) {
	m := profile_manager.FromContext(ctx)
	if m == nil {
		return nil, ProfileError{Name: "", Err: errors.New("couldn't get ProfileManager from context")}
	}

	p, err := m.Create(params.Name)
	if err != nil {
		return nil, ProfileError{Name: params.Name, Err: err}
	}

	if params.Copy != "" {
		src, err := m.Get(params.Copy)
		if err != nil {
			return p, ProfileError{Name: params.Name, Err: err}
		}

		err = m.Copy(src, p)
		if err != nil {
			return p, ProfileError{Name: params.Name, Err: err}
		}
	}

	return p, nil
}
