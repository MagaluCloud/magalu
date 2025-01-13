package tenant

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var getList = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        "list",
			Description: "List all available tenants for current login",
		},
		listTenants,
	)

	exec = core.NewHumanIdentifiableFieldsExecutor(exec, []string{"legal_name", "email"})

	return exec
})

func listTenants(ctx context.Context) ([]*mgcAuthPkg.Tenant, error) {
	auth := mgcAuthPkg.FromContext(ctx)
	if auth == nil {
		return nil, fmt.Errorf("programming error: unable to get auth from context")
	}
	return auth.ListTenants(ctx)
}
