package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkRules "magalu.cloud/lib/products/network/rules"
	networkSecurityGroupsRules "magalu.cloud/lib/products/network/security_groups/rules"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkSecurityGroupRuleModel struct {
	Id              types.String `tfsdk:"id"`
	Description     types.String `tfsdk:"description"`
	Direction       types.String `tfsdk:"direction"`
	Ethertype       types.String `tfsdk:"ethertype"`
	PortRangeMax    types.Int64  `tfsdk:"port_range_max"`
	PortRangeMin    types.Int64  `tfsdk:"port_range_min"`
	Protocol        types.String `tfsdk:"protocol"`
	RemoteGroupId   types.String `tfsdk:"remote_group_id"`
	RemoteIpPrefix  types.String `tfsdk:"remote_ip_prefix"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
}

type NetworkSecurityGroupsRulesResource struct {
	sdkClient                  *mgcSdk.Client
	networkSecurityGroupsRules networkSecurityGroupsRules.Service
	networkRules               networkRules.Service
}

func NewNetworkSecurityGroupsRulesResource() resource.Resource {
	return &NetworkSecurityGroupsRulesResource{}
}

func (r *NetworkSecurityGroupsRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Security Group Rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the security group rule",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the security group rule. Example: 'Allow incoming SSH traffic'",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"direction": schema.StringAttribute{
				Description: "Direction of traffic flow. Allowed values: 'ingress' or 'egress'. Example: 'ingress'",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ingress", "egress"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ethertype": schema.StringAttribute{
				Description: "Network protocol version. Allowed values: 'IPv4' or 'IPv6'. Example: 'IPv4'",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("IPv4", "IPv6"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port_range_max": schema.Int64Attribute{
				Description: "Maximum port number in the range. Valid values: 1-65535. Example: 22",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"port_range_min": schema.Int64Attribute{
				Description: "Minimum port number in the range. Valid values: 1-65535. Example: 22",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Description: "IP protocol. Common values: tcp, udp, icmp. Example: 'tcp'",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"remote_group_id": schema.StringAttribute{
				Description: "ID of the remote security group. Example: 'sg-0123456789abcdef0'",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"remote_ip_prefix": schema.StringAttribute{
				Description: "CIDR notation of remote IP range. Example: '192.168.1.0/24'",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Description: "ID of the security group to which this rule will be added. Example: 'sg-0123456789abcdef0'",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NetworkSecurityGroupsRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	var errDetail error
	r.sdkClient, err, errDetail = client.NewSDKClient(req)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			errDetail.Error(),
		)
		return
	}

	r.networkSecurityGroupsRules = networkSecurityGroupsRules.NewService(ctx, r.sdkClient)
	r.networkRules = networkRules.NewService(ctx, r.sdkClient)
}

func (r *NetworkSecurityGroupsRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_security_groups_rules"
}

func (r *NetworkSecurityGroupsRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworkSecurityGroupRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdRequest := networkSecurityGroupsRules.CreateParameters{
		Description:     data.Description.ValueStringPointer(),
		Direction:       data.Direction.ValueStringPointer(),
		Ethertype:       data.Ethertype.ValueString(),
		PortRangeMax:    tfutil.ConvertInt64PointerToIntPointer(data.PortRangeMax.ValueInt64Pointer()),
		PortRangeMin:    tfutil.ConvertInt64PointerToIntPointer(data.PortRangeMin.ValueInt64Pointer()),
		Protocol:        data.Protocol.ValueStringPointer(),
		RemoteGroupId:   data.RemoteGroupId.ValueStringPointer(),
		RemoteIpPrefix:  data.RemoteIpPrefix.ValueStringPointer(),
		SecurityGroupId: data.SecurityGroupId.ValueString(),
	}

	created, err := r.networkSecurityGroupsRules.CreateContext(ctx, createdRequest, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSecurityGroupsRules.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to create security group rule", err.Error())
		return
	}

	data.Id = types.StringValue(created.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *NetworkSecurityGroupsRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworkSecurityGroupRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.networkRules.GetContext(ctx, networkRules.GetParameters{
		RuleId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkRules.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get security group rule", err.Error())
		return
	}

	data.Description = types.StringPointerValue(rule.Description)
	data.Direction = types.StringPointerValue(rule.Direction)
	data.Ethertype = types.StringPointerValue(rule.Ethertype)
	data.PortRangeMax = types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(rule.PortRangeMax))
	data.PortRangeMin = types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(rule.PortRangeMin))
	data.Protocol = types.StringPointerValue(rule.Protocol)
	data.RemoteGroupId = types.StringPointerValue(rule.RemoteGroupId)
	data.RemoteIpPrefix = types.StringPointerValue(rule.RemoteIpPrefix)
	data.SecurityGroupId = types.StringPointerValue(rule.SecurityGroupId)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *NetworkSecurityGroupsRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkSecurityGroupRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.networkRules.DeleteContext(ctx, networkRules.DeleteParameters{
		RuleId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkRules.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to delete security group rule", err.Error())
		return
	}
}

func (r *NetworkSecurityGroupsRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for security group rules", "")
}

func (r *NetworkSecurityGroupsRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ruleId := req.ID
	data := NetworkSecurityGroupRuleModel{}

	rule, err := r.networkRules.GetContext(ctx, networkRules.GetParameters{
		RuleId: ruleId,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkRules.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to import security group rule", err.Error())
		return
	}

	data.Id = types.StringPointerValue(rule.Id)
	data.Description = types.StringPointerValue(rule.Description)
	data.Direction = types.StringPointerValue(rule.Direction)
	data.Ethertype = types.StringPointerValue(rule.Ethertype)
	data.PortRangeMax = types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(rule.PortRangeMax))
	data.PortRangeMin = types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(rule.PortRangeMin))
	data.Protocol = types.StringPointerValue(rule.Protocol)
	data.RemoteGroupId = types.StringPointerValue(rule.RemoteGroupId)
	data.RemoteIpPrefix = types.StringPointerValue(rule.RemoteIpPrefix)
	data.SecurityGroupId = types.StringPointerValue(rule.SecurityGroupId)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
