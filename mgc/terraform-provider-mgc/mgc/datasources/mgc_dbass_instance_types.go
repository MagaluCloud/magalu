package datasources

import (
	"context"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	dbaas "github.com/MagaluCloud/magalu/mgc/lib/products/dbaas/instance_types"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceDbInstanceTypes{}

type DataSourceDbInstanceTypes struct {
	sdkClient     *mgcSdk.Client
	instanceTypes dbaas.Service
}

type dbInstanceTypeModel struct {
	Status        types.String     `tfsdk:"status"`
	InstanceTypes []DbInstanceType `tfsdk:"instance_types"`
}

type DbInstanceType struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Ram  types.String `tfsdk:"ram"`
	Size types.String `tfsdk:"size"`
	Vcpu types.String `tfsdk:"vcpu"`
}

func NewDataSourceDbaasInstanceTypes() datasource.DataSource {
	return &DataSourceDbInstanceTypes{}
}

func (r *DataSourceDbInstanceTypes) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instance_types"
}

func (r *DataSourceDbInstanceTypes) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.instanceTypes = dbaas.NewService(ctx, r.sdkClient)
}

func (r *DataSourceDbInstanceTypes) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A list of available database instance types.",
		Attributes: map[string]schema.Attribute{
			"instance_types": schema.ListNestedAttribute{
				Description: "List of available database instance types",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the instance type",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the instance type",
							Computed:    true,
						},
						"ram": schema.StringAttribute{
							Description: "The RAM size of the instance type",
							Computed:    true,
						},
						"size": schema.StringAttribute{
							Description: "The size category of the instance type",
							Computed:    true,
						},
						"vcpu": schema.StringAttribute{
							Description: "The number of virtual CPUs",
							Computed:    true,
						},
					},
				},
			},
			"status": schema.StringAttribute{
				Description: "Filter to return instance types based on their status, by default it returns all active instance types. If set to DEPRECATED, returns only deprecated instance types.",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{stringvalidator.OneOf("DEPRECATED", "ACTIVE")},
			},
		},
	}
}

func (r *DataSourceDbInstanceTypes) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := dbInstanceTypeModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	limit := 50
	params := dbaas.ListParameters{
		Status: data.Status.ValueStringPointer(),
		Limit:  &limit,
	}
	instanceTypes, err := r.instanceTypes.ListContext(ctx, params, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaas.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to list db instance types", err.Error())
		return
	}

	instanceTypeModels := make([]DbInstanceType, 0, len(instanceTypes.Results))
	for _, instanceType := range instanceTypes.Results {
		instanceTypeModels = append(instanceTypeModels, DbInstanceType{
			ID:   types.StringValue(instanceType.Id),
			Name: types.StringValue(instanceType.Label),
			Ram:  types.StringValue(instanceType.Ram),
			Size: types.StringValue(instanceType.Size),
			Vcpu: types.StringValue(instanceType.Vcpu),
		})
	}
	data.InstanceTypes = instanceTypeModels
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
