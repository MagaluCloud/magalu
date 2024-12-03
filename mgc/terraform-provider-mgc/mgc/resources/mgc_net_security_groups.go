package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkSecurityGroups "magalu.cloud/lib/products/network/security_groups"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkSecurityGroupsModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	DisableDefaultRules types.Bool   `tfsdk:"disable_default_rules"`
}

type NetworkSecurityGroupsResource struct {
	sdkClient             *mgcSdk.Client
	networkSecurityGroups networkSecurityGroups.Service
}

func NewNetworkSecurityGroupsResource() resource.Resource {
	return &NetworkSecurityGroupsResource{}
}

func (r *NetworkSecurityGroupsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_security_groups"
}

func (r *NetworkSecurityGroupsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.networkSecurityGroups = networkSecurityGroups.NewService(ctx, r.sdkClient)
}

func (r *NetworkSecurityGroupsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Security Group",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the Security Group",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Security Group",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the Security Group",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disable_default_rules": schema.BoolAttribute{
				Description: "Disable default rules, when creating the Security Group",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NetworkSecurityGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworkSecurityGroupsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	sc, err := r.networkSecurityGroups.GetContext(ctx, networkSecurityGroups.GetParameters{
		SecurityGroupId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityGroups.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to read Security Group", err.Error())
		return
	}

	data.Name = types.StringPointerValue(sc.Name)
	if sc.Description != nil && *sc.Description == "" {
		data.Description = types.StringPointerValue(nil)
	} else {
		data.Description = types.StringPointerValue(sc.Description)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *NetworkSecurityGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworkSecurityGroupsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.networkSecurityGroups.CreateContext(ctx, networkSecurityGroupTerraformModelToSdk(data),
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityGroups.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to create Security Group", err.Error())
		return
	}

	data.Id = types.StringValue(created.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func networkSecurityGroupTerraformModelToSdk(create NetworkSecurityGroupsModel) networkSecurityGroups.CreateParameters {
	return networkSecurityGroups.CreateParameters{
		Name:             create.Name.ValueString(),
		Description:      create.Description.ValueStringPointer(),
		SkipDefaultRules: create.DisableDefaultRules.ValueBoolPointer(),
	}
}

func (r *NetworkSecurityGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkSecurityGroupsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.networkSecurityGroups.DeleteContext(ctx, networkSecurityGroups.DeleteParameters{
		SecurityGroupId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityGroups.DeleteConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Security Group", err.Error())
		return
	}
}

func (r *NetworkSecurityGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for Security Group", "")
}

func (r *NetworkSecurityGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	scID := req.ID
	data := NetworkSecurityGroupsModel{}
	sc, err := r.networkSecurityGroups.GetContext(ctx, networkSecurityGroups.GetParameters{
		SecurityGroupId: scID,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityGroups.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to import Security Group", err.Error())
		return
	}

	data.Id = types.StringPointerValue(sc.Id)
	data.Name = types.StringPointerValue(sc.Name)
	data.Description = types.StringPointerValue(sc.Description)
	data.DisableDefaultRules = types.BoolValue(false)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
