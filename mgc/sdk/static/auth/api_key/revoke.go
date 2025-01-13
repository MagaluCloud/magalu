package api_key

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/utils"

	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
)

type revokeParams struct {
	ID string `json:"id" jsonschema_description:"ID of api key to revoke" mgc:"positional"`
}

var getRevoke = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Scopes:      core.Scopes{scope_PA},
			Name:        "revoke",
			Description: "Revoke an API key by its ID",
		},
		revoke,
	)

	msg := "This operation will permanently revoke the api-key {{.parameters.id}}. Do you wish to continue?"

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

	url := fmt.Sprintf("%s/%s/revoke", auth.GetConfig().ApiKeysUrlV1, parameter.ID)
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
