package project

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type listParams struct {
	// TODO: add a description once IAM documents what 'managed' filters
	Managed *bool `json:"managed,omitempty"`
}

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:         "list",
			ProjectScope: core.ProjectScopeIAM,
			Summary:      "List the projects of your tenant",
			Description:  "List all projects available for the currently authenticated tenant",
		},
		list,
	)

	return core.NewHumanIdentifiableFieldsExecutor(exec, []string{"name"})
})

func list(ctx context.Context, params listParams, cfg projectConfig) ([]projectResult, error) {
	client, auth, err := authenticatedClient(ctx)
	if err != nil {
		return nil, err
	}

	r, err := newRequest(ctx, auth, http.MethodGet, projectsURL(ctx, cfg), nil)
	if err != nil {
		return nil, err
	}

	if params.Managed != nil {
		q := r.URL.Query()
		q.Set("managed", strconv.FormatBool(*params.Managed))
		r.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	defer resp.Body.Close()
	var result []projectResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
