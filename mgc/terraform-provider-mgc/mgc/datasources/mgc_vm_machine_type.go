package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkVMMachineTypes "magalu.cloud/lib/products/virtual_machine/machine_types"
	"magalu.cloud/sdk"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DataSourceVmMachineType{}

type DataSourceVmMachineType struct {
	sdkClient      *mgcSdk.Client
	vmMachineTypes sdkVMMachineTypes.Service
}

type MachineTypeModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Disk types.Int64  `tfsdk:"disk"`
	Ram  types.Int64  `tfsdk:"ram"`
	VCPU types.Int64  `tfsdk:"vcpu"`
	GPU  types.Int64  `tfsdk:"gpu"`
}

type MachineTypesModel struct {
	MachineTypes []MachineTypeModel `tfsdk:"machine_types"`
}

func NewDataSourceVmMachineType() datasource.DataSource {
	return &DataSourceVmMachineType{}
}

func (r *DataSourceVmMachineType) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_types"
}

func (r *DataSourceVmMachineType) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(tfutil.ProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	sdk := sdk.NewSdk()
	r.sdkClient = mgcSdk.NewClient(sdk)
	if config.Region.ValueString() != "" {
		_ = r.sdkClient.Sdk().Config().SetTempConfig("region", config.Region.ValueString())
	}
	if config.Env.ValueString() != "" {
		_ = r.sdkClient.Sdk().Config().SetTempConfig("env", config.Env.ValueString())
	}
	if config.ApiKey.ValueString() != "" {
		_ = r.sdkClient.Sdk().Auth().SetAPIKey(config.ApiKey.ValueString())
	}

	r.sdkClient = mgcSdk.NewClient(sdk)
	r.vmMachineTypes = sdkVMMachineTypes.NewService(ctx, r.sdkClient)
}

func (r *DataSourceVmMachineType) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"machine_types": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of available VM machine-types.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of machine-type.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of type.",
						},
						"disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Disk",
						},
						"ram": schema.Int64Attribute{
							Computed:    true,
							Description: "Ram",
						},
						"vcpu": schema.Int64Attribute{
							Computed:    true,
							Description: "VCpu",
						},
						"gpu": schema.Int64Attribute{
							Computed:    true,
							Description: "GPU",
						},
					},
				},
			},
		},
	}
	resp.Schema.Description = "Get the available virtual-machine types."
}

const typeActive string = "active"

func (r *DataSourceVmMachineType) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MachineTypesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.vmMachineTypes.ListContext(ctx, sdkVMMachineTypes.ListParameters{},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVMMachineTypes.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, typ := range sdkOutput.MachineTypes {
		if typ.Status != typeActive {
			continue
		}

		if typ.Gpu == nil {
			typ.Gpu = new(int)
			*typ.Gpu = 0
		}

		gpu := *typ.Gpu

		data.MachineTypes = append(data.MachineTypes, MachineTypeModel{
			ID:   types.StringValue(typ.Id),
			Name: types.StringValue(typ.Name),
			Disk: types.Int64Value(int64(typ.Disk)),
			Ram:  types.Int64Value(int64(typ.Ram)),
			VCPU: types.Int64Value(int64(typ.Vcpus)),
			GPU:  types.Int64Value(int64(gpu)),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
