package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	tfutil "magalu.cloud/terraform-provider-mgc/mgc/tfutil"

	sdkBuckets "magalu.cloud/lib/products/object_storage/buckets"
	"magalu.cloud/sdk"
)

type ObjectStorageBucket struct {
	Bucket            types.String `tfsdk:"bucket"`
	BucketIsPrefix    types.Bool   `tfsdk:"bucket_is_prefix"`
	FinalName         types.String `tfsdk:"final_name"`
	AuthenticatedRead types.Bool   `tfsdk:"authenticated_read"`
	AwsExecRead       types.Bool   `tfsdk:"aws_exec_read"`
	EnableVersioning  types.Bool   `tfsdk:"enable_versioning"`
	GrantFullControl  []Grant      `tfsdk:"grant_full_control"`
	GrantRead         []Grant      `tfsdk:"grant_read"`
	GrantReadACP      []Grant      `tfsdk:"grant_read_acp"`
	GrantWrite        []Grant      `tfsdk:"grant_write"`
	GrantWriteACP     []Grant      `tfsdk:"grant_write_acp"`
	Private           types.Bool   `tfsdk:"private"`
	PublicRead        types.Bool   `tfsdk:"public_read"`
	PublicReadWrite   types.Bool   `tfsdk:"public_read_write"`
	Recursive         types.Bool   `tfsdk:"recursive"`
}

type Grant struct {
	ID types.String `tfsdk:"id"`
}

func NewObjectStorageBucketsResource() resource.Resource {
	return &objectStorageBuckets{}
}

type objectStorageBuckets struct {
	sdkClient *mgcSdk.Client
	buckets   sdkBuckets.Service
}

func (r *objectStorageBuckets) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_buckets"
}

func (r *objectStorageBuckets) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if config.ObjectStorage != nil && config.ObjectStorage.ObjectKeyPair != nil {
		sdk.Config().AddTempKeyPair("apikey", config.ObjectStorage.ObjectKeyPair.KeyID.ValueString(), config.ObjectStorage.ObjectKeyPair.KeySecret.ValueString())
	}

	r.sdkClient = mgcSdk.NewClient(sdk)

	r.buckets = sdkBuckets.NewService(ctx, r.sdkClient)
}

func (r *objectStorageBuckets) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An object storage bucket.",
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "Name of the bucket to be created.",
			},
			"bucket_is_prefix": schema.BoolAttribute{
				Required:    true,
				Description: "Use bucket name as prefix value to generate a unique bucket name.",
			},
			"final_name": schema.StringAttribute{
				Computed:    true,
				Description: "Final name of the bucket, including the prefix and auto-generated suffix.",
			},
			"authenticated_read": schema.BoolAttribute{
				Optional:    true,
				Description: "Owner gets FULL_CONTROL. Authenticated users have READ rights.",
			},
			"aws_exec_read": schema.BoolAttribute{
				Optional: true,
			},
			"enable_versioning": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable versioning for this bucket.",
			},
			"grant_full_control": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Allows grantees FULL_CONTROL.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Either a Tenant ID or a User Project ID.",
					},
				},
			},
			"grant_read": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Allows grantees to list the objects in the bucket.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Either a Tenant ID or a User Project ID.",
					},
				},
			},
			"grant_read_acp": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Allows grantees to read the bucket ACL.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Either a Tenant ID or a User Project ID.",
					},
				},
			},
			"grant_write": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Allows grantees to create objects in the bucket.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Either a Tenant ID or a User Project ID.",
					},
				},
			},
			"grant_write_acp": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Allows grantees to write the ACL for the applicable bucket.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "Either a Tenant ID or a User Project ID.",
					},
				},
			},
			"private": schema.BoolAttribute{
				Optional:    true,
				Description: "Owner gets FULL_CONTROL. Delegated users have access. No one else has access rights.",
			},
			"public_read": schema.BoolAttribute{
				Optional:    true,
				Description: "Owner gets FULL_CONTROL. Everyone else has READ rights.",
			},
			"public_read_write": schema.BoolAttribute{
				Optional:    true,
				Description: "Owner gets FULL_CONTROL. Everyone else has READ and WRITE rights.",
			},
			"recursive": schema.BoolAttribute{
				Optional:    true,
				Description: "Delete bucket including objects inside.",
			},
		},
	}
}

func (r *objectStorageBuckets) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model ObjectStorageBucket
	diags := req.Plan.Get(ctx, &model)

	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	grantFullControl := sdkBuckets.CreateParametersGrantFullControl(convertGrants(model.GrantFullControl))
	grantRead := sdkBuckets.CreateParametersGrantRead(convertGrants(model.GrantRead))
	grantReadACP := sdkBuckets.CreateParametersGrantReadAcp(convertGrants(model.GrantReadACP))
	grantWrite := sdkBuckets.CreateParametersGrantWrite(convertGrants(model.GrantWrite))
	grantWriteACP := sdkBuckets.CreateParametersGrantWriteAcp(convertGrants(model.GrantWriteACP))

	result, err := r.buckets.CreateContext(ctx, sdkBuckets.CreateParameters{
		Bucket:            model.Bucket.ValueString(),
		BucketIsPrefix:    model.BucketIsPrefix.ValueBool(),
		AuthenticatedRead: model.AuthenticatedRead.ValueBoolPointer(),
		AwsExecRead:       model.AwsExecRead.ValueBoolPointer(),
		EnableVersioning:  model.EnableVersioning.ValueBoolPointer(),
		GrantFullControl:  &grantFullControl,
		GrantRead:         &grantRead,
		GrantReadAcp:      &grantReadACP,
		GrantWrite:        &grantWrite,
		GrantWriteAcp:     &grantWriteACP,
		Private:           model.Private.ValueBoolPointer(),
		PublicRead:        model.PublicRead.ValueBoolPointer(),
		PublicReadWrite:   model.PublicReadWrite.ValueBoolPointer(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBuckets.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to create bucket", err.Error())
		return
	}

	model.FinalName = types.StringValue(result.Bucket)

	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func convertGrants(grants []Grant) []sdkBuckets.CreateParametersGrantFullControlItem {
	var result []sdkBuckets.CreateParametersGrantFullControlItem

	for _, grant := range grants {
		result = append(result, sdkBuckets.CreateParametersGrantFullControlItem{Id: grant.ID.ValueString()})
	}

	return result
}

func (r *objectStorageBuckets) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model ObjectStorageBucket
	diags := req.State.Get(ctx, &model)

	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	resp.State = req.State
}

func (r *objectStorageBuckets) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for Buckets creation", "Update is not supported")
}

func (r *objectStorageBuckets) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model ObjectStorageBucket
	diags := req.State.Get(ctx, &model)

	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	name := model.FinalName.ValueString()
	if model.FinalName.IsNull() {
		resp.Diagnostics.AddWarning("Bucket full name is required", "Bucket full name is required")
		name = model.Bucket.ValueString()
	}

	if model.Recursive.IsNull() {
		model.Recursive = types.BoolValue(false)
	}

	_, err := r.buckets.DeleteContext(ctx, sdkBuckets.DeleteParameters{
		Bucket:    name,
		Recursive: model.Recursive.ValueBool(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBuckets.DeleteConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to delete bucket", err.Error())
		return
	}

	resp.State.Set(ctx, &model)
}
