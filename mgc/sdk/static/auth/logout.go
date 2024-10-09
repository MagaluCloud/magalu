package auth

import (
	"context"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	"magalu.cloud/core/utils"
)

var getLogout = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "logout",
			Description: "Run logout",
		},
		func(ctx context.Context, parameters accessTokenParameters, _ struct{}) (output string, err error) {
			auth := mgcAuthPkg.FromContext(ctx)
			err = auth.Logout()
			if err != nil {
				return "fail", err
			}
			return "success", nil
		},
	)
})
