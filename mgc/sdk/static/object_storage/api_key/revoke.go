package api_key

import (
	"context"
	"fmt"
	"net/http"

	"magalu.cloud/core"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"

	mgcAuthPkg "magalu.cloud/core/auth"
)

type revokeParams struct {
	UUID string `json:"uuid" jsonschema_description:"UUID of api key to revoke" mgc:"positional"`
}

var getRevoke = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{scope_PA},
			Name:        "revoke",
			Description: "Revoke credentials used in Object Storage requests",
		},
		revoke,
	)

	msg := "This operation will permanently revoke the api-key {{.parameters.uuid}}. Do you wish to continue?"

	cExecutor := core.NewConfirmableExecutor(
		exec,
		core.ConfirmPromptWithTemplate(msg),
	)

	return core.NewExecuteResultOutputOptions(cExecutor, func(exec core.Executor, result core.Result) string {
		return "template=Api-key {{.uuid}} revoked!\n"
	})
})

func revoke(ctx context.Context, parameter revokeParams, _ struct{}) (*revokeParams, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: unable to auth from context")
	}

	httpClient := auth.AuthenticatedHttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("programming error: unable to get HTTP Client from context")
	}

	url := fmt.Sprintf("%s/%s/revoke", auth.GetConfig().ApiKeysUrlV1, parameter.UUID)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, mgcHttpPkg.NewHttpErrorFromResponse(resp, r)
	}

	return &parameter, nil
}
