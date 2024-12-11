package resources

import (
	"context"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkSecurityAttach "magalu.cloud/lib/products/network/ports"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkSecurityGroupsAttachModel struct {
	SecurityGroupID types.String `tfsdk:"security_group_id"`
	InterfaceID     types.String `tfsdk:"interface_id"`
}

type NetworkSecurityGroupsAttachResource struct {
	sdkClient                   *mgcSdk.Client
	networkSecurityGroupsAttach networkSecurityAttach.Service
}

func NewNetworkSecurityGroupsAttachResource() resource.Resource {
	return &NetworkSecurityGroupsAttachResource{}
}

func (r *NetworkSecurityGroupsAttachResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_security_groups_attach"
}

func (r *NetworkSecurityGroupsAttachResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.networkSecurityGroupsAttach = networkSecurityAttach.NewService(ctx, r.sdkClient)
}

func (r *NetworkSecurityGroupsAttachResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Security Group Attach",
		Attributes: map[string]schema.Attribute{
			"interface_id": schema.StringAttribute{
				Description: "The ID of the Network Interface",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Description: "The ID of the Security Group",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NetworkSecurityGroupsAttachResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NetworkSecurityGroupsAttachModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getParam := networkSecurityAttach.GetParameters{
		PortId: data.InterfaceID.ValueString(),
	}
	interfaceResponse, err := r.networkSecurityGroupsAttach.GetContext(ctx, getParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityAttach.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get Network Interface", err.Error())
		return
	}

	if interfaceResponse.SecurityGroups == nil {
		resp.Diagnostics.AddError("Interface without security groups", "Security Group is nil")
		return
	}

	if !slices.Contains(*interfaceResponse.SecurityGroups, data.SecurityGroupID.ValueString()) {
		resp.Diagnostics.AddError("Security Group not attach to interface", "Security Group not found in interface")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkSecurityGroupsAttachResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworkSecurityGroupsAttachModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attachParam := networkSecurityAttach.AttachParameters{
		PortId:          data.InterfaceID.ValueString(),
		SecurityGroupId: data.SecurityGroupID.ValueString(),
	}
	err := r.networkSecurityGroupsAttach.AttachContext(ctx, attachParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityAttach.AttachConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to attach Security Group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *NetworkSecurityGroupsAttachResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for Network Security Groups Attach", "")
}

func (r *NetworkSecurityGroupsAttachResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkSecurityGroupsAttachModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	detachParam := networkSecurityAttach.DetachParameters{
		PortId:          data.InterfaceID.ValueString(),
		SecurityGroupId: data.SecurityGroupID.ValueString(),
	}
	err := r.networkSecurityGroupsAttach.DetachContext(ctx, detachParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityAttach.DetachConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to detach Security Group", err.Error())
		return
	}
}

func (r *NetworkSecurityGroupsAttachResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var securityGroupID, interfaceID string
	input := strings.Split(req.ID, ",")
	if len(input) != 2 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be in the format security_group_id,interface_id")
		return
	}
	securityGroupID = input[0]
	interfaceID = input[1]

	data := NetworkSecurityGroupsAttachModel{
		SecurityGroupID: types.StringValue(securityGroupID),
		InterfaceID:     types.StringValue(interfaceID),
	}

	getParam := networkSecurityAttach.GetParameters{
		PortId: interfaceID,
	}
	getResponse, err := r.networkSecurityGroupsAttach.GetContext(ctx, getParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityAttach.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to import Security Group Attach", err.Error())
		return
	}

	if getResponse.SecurityGroups == nil {
		resp.Diagnostics.AddError("Interface without security groups", "Security Group is nil")
		return
	}

	if !slices.Contains(*getResponse.SecurityGroups, securityGroupID) {
		resp.Diagnostics.AddError("Security Group not attach to interface", "Security Group not found in interface")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
