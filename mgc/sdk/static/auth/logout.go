package auth

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
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
