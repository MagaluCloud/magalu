package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	tfutil "magalu.cloud/terraform-provider-mgc/mgc/tfutil"

	sdkBucketsAcl "magalu.cloud/lib/products/object_storage/buckets/acl"
	sdkBucketsVersioning "magalu.cloud/lib/products/object_storage/buckets/versioning"
	"magalu.cloud/sdk"
)

var _ datasource.DataSource = &DatasourceBucket{}

type DatasourceBucket struct {
	sdkClient  *mgcSdk.Client
	versioning sdkBucketsVersioning.Service
	acl        sdkBucketsAcl.Service
}
type BucketDetailModel struct {
	Name       types.String         `tfsdk:"name"`
	Versioning types.String         `tfsdk:"versioning"`
	MFADelete  types.String         `tfsdk:"mfadelete"`
	Owner      GenAccessModel       `tfsdk:"owner"`
	Grants     []GranteeAccessModel `tfsdk:"grantee"`
}

type GenAccessModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
}

type GranteeAccessModel struct {
	GenAccessModel
	Permission types.String `tfsdk:"permission"`
}

func NewDatasourceBucket() datasource.DataSource {
	return &DatasourceBucket{}
}

func (r *DatasourceBucket) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_bucket"
}

func (r *DatasourceBucket) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	if config.ObjectStorage != nil && config.ObjectStorage.ObjectKeyPair != nil {
		sdk.Config().AddTempKeyPair("apikey", config.ObjectStorage.ObjectKeyPair.KeyID.ValueString(), config.ObjectStorage.ObjectKeyPair.KeySecret.ValueString())
	}

	r.versioning = sdkBucketsVersioning.NewService(ctx, r.sdkClient)
	r.acl = sdkBucketsAcl.NewService(ctx, r.sdkClient)
}

func (r *DatasourceBucket) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"buckets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of ssh-keys.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Bucket name",
						},
						"versioning": schema.StringAttribute{
							Computed:    true,
							Description: "Versioning status",
						},
						"mfadelete": schema.StringAttribute{
							Computed:    true,
							Description: "MFA Delete",
						},
						"owner": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Bucket owner",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "Owner ID",
									},
									"display_name": schema.StringAttribute{
										Computed:    true,
										Description: "Owner Name",
									},
								},
							},
						},
						"grantee": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Bucket grantee",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "Grantee ID",
									},
									"display_name": schema.StringAttribute{
										Computed:    true,
										Description: "Grantee Name",
									},
									"permission": schema.StringAttribute{
										Computed:    true,
										Description: "Grantee permission",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	resp.Schema.Description = "Get details of bucket."
}

func (r *DatasourceBucket) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BucketDetailModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	versioning, err := r.versioning.GetContext(ctx, sdkBucketsVersioning.GetParameters{Bucket: data.Name.ValueString()},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBucketsVersioning.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versioning", err.Error())
		return
	}
	acl, err := r.acl.GetContext(ctx, sdkBucketsAcl.GetParameters{Dst: data.Name.ValueString()},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBucketsAcl.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get acl", err.Error())
		return
	}

	for _, aclDetail := range acl.AccessControlList.Grant {
		data.Grants = append(data.Grants, GranteeAccessModel{
			Permission: types.StringValue(aclDetail.Permission),
			GenAccessModel: GenAccessModel{
				DisplayName: types.StringValue(aclDetail.Grantee.DisplayName),
				ID:          types.StringValue(aclDetail.Grantee.Id),
			},
		})
	}

	data.MFADelete = types.StringValue(versioning.MfaDelete)
	data.Versioning = types.StringValue(versioning.Status)

	data.Owner.DisplayName = types.StringValue(acl.Owner.DisplayName)
	data.Owner.ID = types.StringValue(acl.Owner.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
