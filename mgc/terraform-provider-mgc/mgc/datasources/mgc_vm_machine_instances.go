package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkVMInstances "magalu.cloud/lib/products/virtual_machine/instances"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DataSourceVmInstances{}

type DataSourceVmInstances struct {
	sdkClient   *mgcSdk.Client
	vmInstances sdkVMInstances.Service
}

type VMInstanceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	PublicIPV4    types.String `tfsdk:"public_ipv4"`
	PublicIPV6    types.String `tfsdk:"public_ipv6"`
	PrivateIPV4   types.String `tfsdk:"private_ipv4"`
	SshKeyName    types.String `tfsdk:"ssh_key_name"`
	Status        types.String `tfsdk:"status"`
	State         types.String `tfsdk:"state"`
	ImageID       types.String `tfsdk:"image_id"`
	MachineTypeID types.String `tfsdk:"machine_type_id"`
}

type VMInstancesModel struct {
	Instances []VMInstanceModel `tfsdk:"instances"`
}

func NewDataSourceVmInstances() datasource.DataSource {
	return &DataSourceVmInstances{}
}

func (r *DataSourceVmInstances) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_instances"
}

func (r *DataSourceVmInstances) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	r.sdkClient, err = client.NewSDKClient(req)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.vmInstances = sdkVMInstances.NewService(ctx, r.sdkClient)
}

func (r *DataSourceVmInstances) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of available VM instances.",
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
				},
			},
		},
	}
	resp.Schema.Description = "Get the available virtual-machine instances."
}

func (r *DataSourceVmInstances) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VMInstancesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.vmInstances.ListContext(ctx, sdkVMInstances.ListParameters{Expand: &sdkVMInstances.ListParametersExpand{"network"}},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVMInstances.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get instances", err.Error())
		return
	}

	for _, instance := range sdkOutput.Instances {
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

		data.Instances = append(data.Instances, VMInstanceModel{
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
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
