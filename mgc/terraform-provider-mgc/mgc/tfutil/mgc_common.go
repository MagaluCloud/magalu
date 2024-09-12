package tfutil

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	"magalu.cloud/sdk"
)

type GenericIDNameModel struct {
	Name types.String `tfsdk:"name"`
	ID   types.String `tfsdk:"id"`
}

type GenericIDModel struct {
	ID types.String `tfsdk:"id"`
}

type SDKFrom interface {
	resource.ConfigureRequest | datasource.ConfigureRequest
}

func SDKClientGenerator[T SDKFrom](req T) (*mgcSdk.Client, error) {
	var config ProviderConfig

	switch tp := any(req).(type) {
	case resource.ConfigureRequest:
		if cfg, ok := tp.ProviderData.(ProviderConfig); ok {
			config = cfg
			break
		}
		return nil, fmt.Errorf("unexpected Resource Configure Type")
	case datasource.ConfigureRequest:
		if cfg, ok := tp.ProviderData.(ProviderConfig); ok {
			config = cfg
			break
		}
		return nil, fmt.Errorf("unexpected Data Source Configure Type")
	default:
		return nil, fmt.Errorf("")
	}

	local_sdk := sdk.NewSdk()
	sdkClient := mgcSdk.NewClient(local_sdk)

	if config.Region.ValueString() != "" {
		_ = sdkClient.Sdk().Config().SetTempConfig("region", config.Region.ValueString())
	}
	if config.Env.ValueString() != "" {
		_ = sdkClient.Sdk().Config().SetTempConfig("env", config.Env.ValueString())
	}
	if config.ApiKey.ValueString() != "" {
		_ = sdkClient.Sdk().Auth().SetAPIKey(config.ApiKey.ValueString())
	}

	if config.ObjectStorage != nil && config.ObjectStorage.ObjectKeyPair != nil {
		sdkClient.Sdk().Config().AddTempKeyPair("apikey", config.ObjectStorage.ObjectKeyPair.KeyID.ValueString(),
			config.ObjectStorage.ObjectKeyPair.KeySecret.ValueString())
	}

	return sdkClient, nil
}
