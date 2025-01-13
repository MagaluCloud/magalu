package datasources

import (
	"context"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	dbaas "github.com/MagaluCloud/magalu/mgc/lib/products/dbaas/instances"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceDbInstance{}

type DataSourceDbInstance struct {
	sdkClient *mgcSdk.Client
	instances dbaas.Service
}

type dbInstanceDataSourceModel struct {
	ID                  types.String            `tfsdk:"id"`
	Addresses           []InstanceAddress       `tfsdk:"addresses"`
	BackupRetentionDays types.Int64             `tfsdk:"backup_retention_days"`
	CreatedAt           types.String            `tfsdk:"created_at"`
	EngineID            types.String            `tfsdk:"engine_id"`
	InstanceTypeID      types.String            `tfsdk:"instance_type_id"`
	Name                types.String            `tfsdk:"name"`
	Parameters          map[string]types.String `tfsdk:"parameters"`
	Status              types.String            `tfsdk:"status"`
	VolumeSize          types.Int64             `tfsdk:"volume_size"`
	VolumeType          types.String            `tfsdk:"volume_type"`
}

func NewDataSourceDbaasInstance() datasource.DataSource {
	return &DataSourceDbInstance{}
}

func (r *DataSourceDbInstance) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instance"
}

func (r *DataSourceDbInstance) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *DataSourceDbInstance) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a database instance by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the instance to fetch",
				Required:    true,
			},
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
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp of the instance",
				Computed:    true,
			},
			"engine_id": schema.StringAttribute{
				Description: "ID of the engine",
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
		},
	}
}

func (r *DataSourceDbInstance) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dbInstanceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := dbaas.GetParameters{
		InstanceId: data.ID.ValueString(),
	}
	instance, err := r.instances.GetContext(ctx, params, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaas.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to get db instance", err.Error())
		return
	}

	var addresses []InstanceAddress
	for _, address := range instance.Addresses {
		addresses = append(addresses, InstanceAddress{
			Access:  types.StringValue(address.Access),
			Address: types.StringPointerValue(address.Address),
			Type:    types.StringPointerValue(address.Type),
		})
	}

	data.Addresses = addresses
	data.BackupRetentionDays = types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(&instance.BackupRetentionDays))
	data.CreatedAt = types.StringValue(instance.CreatedAt)
	data.EngineID = types.StringValue(instance.EngineId)
	data.InstanceTypeID = types.StringValue(instance.InstanceTypeId)
	data.Name = types.StringValue(instance.Name)
	data.Parameters = convertToStringMapInstance(instance.Parameters)
	data.Status = types.StringValue(instance.Status)
	data.VolumeSize = types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(&instance.Volume.Size))
	data.VolumeType = types.StringValue(instance.Volume.Type)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertToStringMapInstance(params dbaas.GetResultParameters) map[string]types.String {
	result := make(map[string]types.String, len(params))
	for _, value := range params {
		result[value.Name] = types.StringValue(tfutil.SdkParamValueToString(value.Value))
	}
	return result
}
