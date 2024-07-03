module magalu.cloud/lib

go 1.22.3

require (
	magalu.cloud/core v0.19.4-rc6-internal
	magalu.cloud/sdk  v0.19.4-rc6-internal
)

replace magalu.cloud/core => ../core

replace magalu.cloud/sdk => ../sdk
