module magalu.cloud/lib

go 1.22.3

require (
	magalu.cloud/core v0.23.0
	magalu.cloud/sdk  v0.23.0
)

replace magalu.cloud/core => ../core

replace magalu.cloud/sdk => ../sdk
