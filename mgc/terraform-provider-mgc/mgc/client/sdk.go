package client

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"

	mgcSdk "magalu.cloud/lib"
	"magalu.cloud/sdk"
)

type SDKFrom interface {
	resource.ConfigureRequest | datasource.ConfigureRequest
}

func NewSDKClient[T SDKFrom, G tfutil.ResponseFrom](req T, resp *G) (*mgcSdk.Client, error, error) {
	var config tfutil.ProviderConfig

	devErrMsg := "fail to parse provider config"
	switch tp := any(req).(type) {
	case resource.ConfigureRequest:
		if cfg, ok := tp.ProviderData.(tfutil.ProviderConfig); ok {
			config = cfg
			break
		}
		return nil, fmt.Errorf("%s", devErrMsg), fmt.Errorf("unexpected Resource Configure Type")
	case datasource.ConfigureRequest:
		if cfg, ok := tp.ProviderData.(tfutil.ProviderConfig); ok {
			config = cfg
			break
		}
		return nil, fmt.Errorf("%s", devErrMsg), fmt.Errorf("unexpected Data Source Configure Type")
	default:
		return nil, fmt.Errorf("%s", devErrMsg), fmt.Errorf("provider data is null")
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
	} else {
		// Remove this block when the provider does not support CLI Auth
		tfutil.AddCLIAuthWarning(resp)
	}

	if config.ObjectStorage != nil && config.ObjectStorage.ObjectKeyPair != nil {
		sdkClient.Sdk().Config().AddTempKeyPair("apikey", config.ObjectStorage.ObjectKeyPair.KeyID.ValueString(),
			config.ObjectStorage.ObjectKeyPair.KeySecret.ValueString())
	}

	return sdkClient, nil, nil
}
