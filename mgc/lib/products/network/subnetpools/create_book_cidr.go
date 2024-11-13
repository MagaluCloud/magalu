/*
Executor: create-book-cidr

# Summary

# Book Subnetpool

# Description

# Booking a CIDR range from a subnetpool

Version: 1.141.3

import "magalu.cloud/lib/products/network/subnetpools"
*/
package subnetpools

import (
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type CreateBookCidrParameters struct {
	Cidr         *string `json:"cidr,omitempty"`
	Mask         *int    `json:"mask,omitempty"`
	SubnetpoolId string  `json:"subnetpool_id"`
}

type CreateBookCidrConfigs struct {
	Env       *string `json:"env,omitempty"`
	Region    *string `json:"region,omitempty"`
	ServerUrl *string `json:"serverUrl,omitempty"`
}

type CreateBookCidrResult struct {
	Cidr *string `json:"cidr,omitempty"`
}

func (s *service) CreateBookCidr(
	parameters CreateBookCidrParameters,
	configs CreateBookCidrConfigs,
) (
	result CreateBookCidrResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("CreateBookCidr", mgcCore.RefPath("/network/subnetpools/create-book-cidr"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateBookCidrParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateBookCidrConfigs](configs); err != nil {
		return
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateBookCidrResult](r)
}

// Context from caller is used to allow cancellation of long-running requests
func (s *service) CreateBookCidrContext(
	ctx context.Context,
	parameters CreateBookCidrParameters,
	configs CreateBookCidrConfigs,
) (
	result CreateBookCidrResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("CreateBookCidr", mgcCore.RefPath("/network/subnetpools/create-book-cidr"), s.client, ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[CreateBookCidrParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs
	if c, err = mgcHelpers.ConvertConfigs[CreateBookCidrConfigs](configs); err != nil {
		return
	}

	sdkConfig := s.client.Sdk().Config().TempConfig()
	if c["serverUrl"] == nil && sdkConfig["serverUrl"] != nil {
		c["serverUrl"] = sdkConfig["serverUrl"]
	}

	if c["env"] == nil && sdkConfig["env"] != nil {
		c["env"] = sdkConfig["env"]
	}

	if c["region"] == nil && sdkConfig["region"] != nil {
		c["region"] = sdkConfig["region"]
	}

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[CreateBookCidrResult](r)
}

// TODO: links
// TODO: related
