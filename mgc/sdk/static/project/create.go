package project

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type createParams struct {
	Name string `json:"name" jsonschema:"description=Name of the new project,required,example=my-project" mgc:"positional"`
}

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:          "create",
			ProjectScoped: true,
			Summary:       "Create a new project",
			Description:   "Create a new project under the currently authenticated tenant",
		},
		create,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Project created successfully. Id={{.id}} Name={{.name}}\n"
	})
})

func create(ctx context.Context, params createParams, cfg projectConfig) (*createResult, error) {
	client, auth, err := authenticatedClient(ctx)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(projectCreate{Name: params.Name}); err != nil {
		return nil, err
	}

	r, err := newRequest(ctx, auth, http.MethodPost, projectsURL(ctx, cfg), &buf)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	// IAM answers 200 on creation, not 201.
	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result createResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
