package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkVMInstances "magalu.cloud/lib/products/virtual_machine/instances"
	"magalu.cloud/sdk"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DataSourceVmInstance{}

type DataSourceVmInstance struct {
	sdkClient   *mgcSdk.Client
	vmInstances sdkVMInstances.Service
}

func NewDataSourceVmInstance() datasource.DataSource {
	return &DataSourceVmInstance{}
}

func (r *DataSourceVmInstance) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_instance"
}

func (r *DataSourceVmInstance) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected provider config, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	r.sdkClient = mgcSdk.NewClient(sdk)
	r.vmInstances = sdkVMInstances.NewService(ctx, r.sdkClient)
}

func (r *DataSourceVmInstance) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of machine-type.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of type.",
			},
			"public_ipv4": schema.StringAttribute{
				Computed:    true,
				Description: "Public IPV4.",
			},
			"public_ipv6": schema.StringAttribute{
				Computed:    true,
				Description: "Public IPV6.",
			},
			"private_ipv4": schema.StringAttribute{
				Computed:    true,
				Description: "Private IPV4",
			},
			"ssh_key_name": schema.StringAttribute{
				Computed:    true,
				Description: "SSH Key name",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of instance.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "State of instance",
			},
			"image_id": schema.StringAttribute{
				Computed:    true,
				Description: "Image ID of instance",
			},
			"machine_type_id": schema.StringAttribute{
				Computed:    true,
				Description: "Machine type ID of instance",
			},
		},
	}
	resp.Schema.Description = "Get the available virtual-machine instance details"
}

func (r *DataSourceVmInstance) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VMInstanceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.vmInstances.Get(sdkVMInstances.GetParameters{Id: data.ID.ValueString(), Expand: &sdkVMInstances.GetParametersExpand{"network"}},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVMInstances.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to get instance", err.Error())
		return
	}

	instance := sdkOutput

	privateIpAddress := ""
	publicIpv6Adress := ""
	publicIpv4Adress := ""

	for _, port := range *instance.Network.Ports {
		privateIpAddress = port.IpAddresses.PrivateIpAddress
		if port.IpAddresses.PublicIpAddress != nil {
			publicIpv4Adress = *port.IpAddresses.PublicIpAddress
		}
		if port.IpAddresses.IpV6address != nil {
			publicIpv6Adress = *port.IpAddresses.IpV6address
		}
	}

	data = VMInstanceModel{
		ID:            types.StringValue(instance.Id),
		Name:          types.StringValue(*instance.Name),
		PublicIPV4:    types.StringValue(publicIpv4Adress),
		PublicIPV6:    types.StringValue(publicIpv6Adress),
		PrivateIPV4:   types.StringValue(privateIpAddress),
		SshKeyName:    types.StringValue(*instance.SshKeyName),
		Status:        types.StringValue(instance.Status),
		State:         types.StringValue(instance.State),
		ImageID:       types.StringValue(instance.Image.Id),
		MachineTypeID: types.StringValue(instance.MachineType.Id),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
