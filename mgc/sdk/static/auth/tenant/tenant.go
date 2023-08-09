package tenant

import "magalu.cloud/core"

func NewTenant() core.Grouper {
	return core.NewStaticGroup(
		"tenant",
		"",
		"Tenant-related operations",
		[]core.Descriptor{
			newTenantList(),
			newTenantSelect(),
			newTenantCurrent(),
		},
	)
}
