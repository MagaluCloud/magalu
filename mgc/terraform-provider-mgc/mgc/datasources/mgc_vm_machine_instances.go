package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	sdkVMInstances "github.com/MagaluCloud/magalu/mgc/lib/products/virtual_machine/instances"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceVmInstances{}

type DataSourceVmInstances struct {
	sdkClient   *mgcSdk.Client
	vmInstances sdkVMInstances.Service
}

type VMInstancesItemModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	SshKeyName       types.String `tfsdk:"ssh_key_name"`
	Status           types.String `tfsdk:"status"`
	State            types.String `tfsdk:"state"`
	ImageID          types.String `tfsdk:"image_id"`
	MachineTypeID    types.String `tfsdk:"machine_type_id"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
}

type VMInstancesModel struct {
	Instances []VMInstancesItemModel `tfsdk:"instances"`
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
	var errDetail error
	r.sdkClient, err, errDetail = client.NewSDKClient(req, resp)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			errDetail.Error(),
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
						"availability_zone": schema.StringAttribute{
							Computed:    true,
							Description: "Availability zone of instance",
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
	if resp.Diagnostics.HasError() {
		return
	}

	instances, err := r.vmInstances.ListContext(ctx, sdkVMInstances.ListParameters{},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVMInstances.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get instances", err.Error())
		return
	}

	for _, instance := range instances.Instances {
		data.Instances = append(data.Instances, VMInstancesItemModel{
			ID:               types.StringValue(instance.Id),
			Name:             types.StringPointerValue(instance.Name),
			SshKeyName:       types.StringPointerValue(instance.SshKeyName),
			Status:           types.StringValue(instance.Status),
			State:            types.StringValue(instance.State),
			ImageID:          types.StringValue(instance.Image.Id),
			MachineTypeID:    types.StringValue(instance.MachineType.Id),
			AvailabilityZone: types.StringPointerValue(instance.AvailabilityZone),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
