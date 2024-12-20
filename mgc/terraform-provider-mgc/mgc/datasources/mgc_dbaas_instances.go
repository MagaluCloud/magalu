package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	dbaas "magalu.cloud/lib/products/dbaas/instances"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DataSourceDbInstances{}

type DataSourceDbInstances struct {
	sdkClient *mgcSdk.Client
	instances dbaas.Service
}

type dbInstanceModel struct {
	Instances []dbInstance `tfsdk:"instances"`
	Status    types.String `tfsdk:"status"`
}

type dbInstance struct {
	Addresses           []InstanceAddress       `tfsdk:"addresses"`
	BackupRetentionDays types.Int64             `tfsdk:"backup_retention_days"`
	CreatedAt           types.String            `tfsdk:"created_at"`
	EngineID            types.String            `tfsdk:"engine_id"`
	ID                  types.String            `tfsdk:"id"`
	InstanceTypeID      types.String            `tfsdk:"instance_type_id"`
	Name                types.String            `tfsdk:"name"`
	Parameters          map[string]types.String `tfsdk:"parameters"`
	Status              types.String            `tfsdk:"status"`
	VolumeSize          types.Int64             `tfsdk:"volume_size"`
	VolumeType          types.String            `tfsdk:"volume_type"`
}

type InstanceAddress struct {
	Access  types.String `tfsdk:"access"`
	Address types.String `tfsdk:"address"`
	Type    types.String `tfsdk:"type"`
}

func NewDataSourceDbaasInstances() datasource.DataSource {
	return &DataSourceDbInstances{}
}

func (r *DataSourceDbInstances) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instances"
}

func (r *DataSourceDbInstances) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.instances = dbaas.NewService(ctx, r.sdkClient)
}

func (r *DataSourceDbInstances) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A list of database instances.",
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Description: "List of database instances",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"addresses": schema.ListNestedAttribute{
							Description: "List of instance addresses",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"access": schema.StringAttribute{
										Description: "Access type of the address",
										Computed:    true,
									},
									"address": schema.StringAttribute{
										Description: "IP address",
										Computed:    true,
									},
									"type": schema.StringAttribute{
										Description: "Type of the address",
										Computed:    true,
									},
								},
							},
						},
						"backup_retention_days": schema.Int64Attribute{
							Description: "Number of days to retain backups",
							Computed:    true,
						},
						"engine_id": schema.StringAttribute{
							Description: "ID of the engine",
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "ID of the instance",
							Computed:    true,
						},
						"instance_type_id": schema.StringAttribute{
							Description: "ID of the instance type",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the instance",
							Computed:    true,
						},
						"parameters": schema.MapAttribute{
							Description: "Map of parameters",
							Computed:    true,
							ElementType: types.StringType,
						},
						"status": schema.StringAttribute{
							Description: "Status of the instance",
							Computed:    true,
						},
						"volume_size": schema.Int64Attribute{
							Description: "Size of the volume",
							Computed:    true,
						},
						"volume_type": schema.StringAttribute{
							Description: "Type of the volume",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Creation timestamp of the instance",
							Computed:    true,
						},
					},
				},
			},
			"status": schema.StringAttribute{
				Description: "Status of the instances",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "BACKING_UP", "CREATING", "DELETED", "DELETING", "ERROR", "ERROR_DELETING", "MAINTENANCE", "PENDING", "REBOOT", "RESIZING", "RESTORING", "STARTING", "STOPPED", "STOPPING"),
				},
			},
		},
	}
}

func (r *DataSourceDbInstances) getAllInstances(ctx context.Context, params dbaas.ListParameters, configs dbaas.ListConfigs) ([]dbaas.ListResultResultsItem, error) {
	var allResults []dbaas.ListResultResultsItem
	params.Offset = new(int)
	for {
		instances, err := r.instances.ListContext(ctx, params, configs)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, instances.Results...)
		currentCount := *params.Offset + instances.Meta.Page.Count
		if currentCount >= instances.Meta.Page.Total {
			break
		}
		*params.Offset = currentCount
	}
	return allResults, nil
}

func (r *DataSourceDbInstances) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := dbInstanceModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	params := dbaas.ListParameters{
		Limit:  new(int),
		Status: data.Status.ValueStringPointer(),
	}
	*params.Limit = 25
	instances, err := r.getAllInstances(ctx, params, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaas.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to list db instances", err.Error())
		return
	}

	var instanceModels []dbInstance
	for _, instance := range instances {
		var addresses []InstanceAddress
		for _, address := range instance.Addresses {
			addresses = append(addresses, InstanceAddress{
				Access:  types.StringValue(address.Access),
				Address: types.StringPointerValue(address.Address),
				Type:    types.StringPointerValue(address.Type),
			})
		}
		instanceModels = append(instanceModels, dbInstance{
			Addresses:           addresses,
			BackupRetentionDays: types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(&instance.BackupRetentionDays)),
			CreatedAt:           types.StringValue(instance.CreatedAt),
			EngineID:            types.StringValue(instance.EngineId),
			ID:                  types.StringValue(instance.Id),
			InstanceTypeID:      types.StringValue(instance.InstanceTypeId),
			Name:                types.StringValue(instance.Name),
			Parameters:          convertToStringMapInstances(instance.Parameters),
			Status:              types.StringValue(instance.Status),
			VolumeSize:          types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(&instance.Volume.Size)),
			VolumeType:          types.StringValue(instance.Volume.Type),
		})
	}
	data.Instances = instanceModels
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertToStringMapInstances(params dbaas.ListResultResultsItemParameters) map[string]types.String {
	result := make(map[string]types.String, len(params))
	for _, value := range params {
		result[value.Name] = types.StringValue(tfutil.SdkParamValueToString(value.Value))
	}
	return result
}
