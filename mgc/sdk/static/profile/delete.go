package profile

import (
	"context"
	"errors"

	"magalu.cloud/core"
	"magalu.cloud/core/profile_manager"
	"magalu.cloud/core/utils"
)

type deleteParams struct {
	Name string `json:"name" jsonschema_description:"Profile name" mgc:"positional"`
}

var getDelete = utils.NewLazyLoader[core.Executor](func() core.Executor {
	exec := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "delete",
			Description: "Deletes the profile with the specified name",
		},
		delete,
	)

	msg := "This operation will permanently delete profile {{.parameters.name}} and it's contents. Do you wish to continue?"

	cExecutor := core.NewConfirmableExecutor(
		exec,
		core.ConfirmPromptWithTemplate(msg),
	)

	return core.NewExecuteResultOutputOptions(cExecutor, func(exec core.Executor, result core.Result) string {
		return "template=Deleted profile {{.Name}}\n"
	})
})

func delete(ctx context.Context, params deleteParams, _ struct{}) (*profile_manager.Profile, error) {
	m := profile_manager.FromContext(ctx)
	if m == nil {
		return nil, ProfileError{Name: "", Err: errors.New("Couldn't get ProfileManager from context")}
	}

	p, err := m.Get(params.Name)
	if err != nil {
		return nil, ProfileError{Name: params.Name, Err: err}
	}

	err = m.Delete(p)
	if err != nil {
		return p, ProfileError{Name: params.Name, Err: err}
	}

	return p, nil
}
