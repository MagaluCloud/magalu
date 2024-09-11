package workspace

import (
	"context"
	"errors"

	"magalu.cloud/core"
	"magalu.cloud/core/profile_manager"
	"magalu.cloud/core/utils"
)

type deleteParams struct {
	Name string `json:"name" jsonschema_description:"workspace name" mgc:"positional"`
}

var getDelete = utils.NewLazyLoader[core.Executor](func() core.Executor {
	exec := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "delete",
			Description: "Deletes the workspace with the specified name",
		},
		delete,
	)

	msg := "This operation will permanently delete workspace {{.parameters.name}} and it's contents. Do you wish to continue?"

	cExecutor := core.NewConfirmableExecutor(
		exec,
		core.ConfirmPromptWithTemplate(msg),
	)

	return core.NewExecuteResultOutputOptions(cExecutor, func(exec core.Executor, result core.Result) string {
		return "template=Deleted workspace {{.name}}\n"
	})
})

func delete(ctx context.Context, params deleteParams, _ struct{}) (*profile_manager.Profile, error) {
	m := profile_manager.FromContext(ctx)
	if m == nil {
		return nil, WorkspaceError{Name: "", Err: errors.New("couldn't get ProfileManager from context")}
	}

	p, err := m.Get(params.Name)
	if err != nil {
		return nil, WorkspaceError{Name: params.Name, Err: err}
	}

	err = m.Delete(p)
	if err != nil {
		return p, WorkspaceError{Name: params.Name, Err: err}
	}

	return p, nil
}
