package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	mgcSdk "magalu.cloud/lib"
	"magalu.cloud/sdk"
)

var _ provider.Provider = (*MgcProvider)(nil)

const providerTypeName = "mgc"

type MgcProvider struct {
	version string
	commit  string
	date    string
	sdk     *sdk.Sdk
}

type ProviderConfig struct {
	Region types.String `tfsdk:"region"`
	Env    types.String `tfsdk:"env"`
}

func (p *MgcProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	tflog.Debug(ctx, "setting provider metadata")
	resp.TypeName = providerTypeName
	resp.Version = p.version
}

//https://github.com/hashicorp/terraform-plugin-framework-validators

func (p *MgcProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	tflog.Debug(ctx, "setting provider schema")

	resp.Schema = schema.Schema{
		Description: "Terraform Provider for Magalu Cloud",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				MarkdownDescription: "Region",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("br-ne1", "br-se1", "br-mgl1"),
				},
			},
			"env": schema.StringAttribute{
				MarkdownDescription: "Environment",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("prod", "pre-prod"),
				},
			},
		},
	}
}

func (p *MgcProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "configuring MGC provider")

	var data ProviderConfig

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "fail to get configs from provider")
	}

	if !data.Region.IsNull() {
		if err := p.sdk.Config().SetTempConfig("region", data.Region.String()); err != nil {
			tflog.Error(ctx, "fail to set region")
		}
	}

	if !data.Env.IsNull() {
		if err := p.sdk.Config().SetTempConfig("env", data.Env.String()); err != nil {
			tflog.Error(ctx, "fail to set env")
		}
	}

	resp.DataSourceData = p.sdk
	resp.ResourceData = p.sdk
}

func (p *MgcProvider) Resources(ctx context.Context) []func() resource.Resource {
	tflog.Info(ctx, "configuring MGC provider resources")

	return []func() resource.Resource{
		NewVirtualMachineInstancesResource,
		NewVirtualMachineSnapshotsResource,
	}
}

func (p *MgcProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string, commit string, date string) func() provider.Provider {
	sdk := mgcSdk.DefaultSdk()

	return func() provider.Provider {
		return &MgcProvider{
			sdk:     sdk,
			version: version,
			commit:  commit,
			date:    date,
		}
	}
}
