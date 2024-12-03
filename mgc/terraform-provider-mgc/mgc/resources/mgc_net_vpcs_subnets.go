package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	mgcSdk "magalu.cloud/lib"
	networkSubnets "magalu.cloud/lib/products/network/subnets"
	networkVpcsSubnets "magalu.cloud/lib/products/network/vpcs/subnets"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type mgcNetworkVpcsSubnetsModel struct {
	ID             types.String   `tfsdk:"id"`
	CidrBlock      types.String   `tfsdk:"cidr_block"`
	Description    types.String   `tfsdk:"description"`
	DnsNameservers []types.String `tfsdk:"dns_nameservers"`
	IpVersion      types.String   `tfsdk:"ip_version"`
	Name           types.String   `tfsdk:"name"`
	SubnetpoolId   types.String   `tfsdk:"subnetpool_id"`
	VpcId          types.String   `tfsdk:"vpc_id"`
}

type mgcNetworkVpcsSubnetsResource struct {
	sdkClient          *mgcSdk.Client
	networkVpcsSubnets networkVpcsSubnets.Service
	networkSubnets     networkSubnets.Service
}

func NewNetworkVpcsSubnetsResource() resource.Resource {
	return &mgcNetworkVpcsSubnetsResource{}
}

func (r *mgcNetworkVpcsSubnetsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_vpcs_subnets"
}

func (r *mgcNetworkVpcsSubnetsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network VPC Subnet",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the VPC subnet",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cidr_block": schema.StringAttribute{
				Description: "The CIDR block of the VPC subnet. Example: '192.168.1.0/24', '0.0.0.0/0', '::/0' or '2001:db8::/32'",
				Required:    true,
				Validators: []validator.String{
					tfutil.CidrValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the VPC subnet",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dns_nameservers": schema.ListAttribute{
				Description: "The DNS nameservers of the VPC subnet",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"ip_version": schema.StringAttribute{
				Description: "Network protocol version. Allowed values: 'IPv4' or 'IPv6'. Example: 'IPv4'",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("IPv4", "IPv6"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the VPC subnet",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnetpool_id": schema.StringAttribute{
				Description: "The ID of the subnet pool",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Description: "The ID of the VPC",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *mgcNetworkVpcsSubnetsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.networkVpcsSubnets = networkVpcsSubnets.NewService(ctx, r.sdkClient)
	r.networkSubnets = networkSubnets.NewService(ctx, r.sdkClient)
}

func (r *mgcNetworkVpcsSubnetsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := mgcNetworkVpcsSubnetsModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsCreateParam := networkVpcsSubnets.CreateParametersDnsNameservers{}
	for _, dns := range data.DnsNameservers {
		dnsCreateParam = append(dnsCreateParam, dns.ValueString())
	}

	createParam := networkVpcsSubnets.CreateParameters{
		CidrBlock:      data.CidrBlock.ValueString(),
		Description:    data.Description.ValueStringPointer(),
		DnsNameservers: &dnsCreateParam,
		IpVersion:      convertIPStringToIPVersion(data.IpVersion.ValueString()),
		Name:           data.Name.ValueString(),
		SubnetpoolId:   data.SubnetpoolId.ValueStringPointer(),
		VpcId:          data.VpcId.ValueString(),
	}

	subnet, err := r.networkVpcsSubnets.CreateContext(ctx, createParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVpcsSubnets.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to create VPC subnet", err.Error())
		return
	}

	data.ID = types.StringPointerValue(&subnet.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mgcNetworkVpcsSubnetsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := &mgcNetworkVpcsSubnetsModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, err := r.networkSubnets.GetContext(ctx, networkSubnets.GetParameters{
		SubnetId: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubnets.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get VPC subnet", err.Error())
		return
	}

	dnsNameServers := []types.String{}
	for _, dns := range subnet.DnsNameservers {
		dnsNameServers = append(dnsNameServers, types.StringPointerValue(&dns))
	}

	data.DnsNameservers = dnsNameServers
	data.CidrBlock = types.StringPointerValue(&subnet.CidrBlock)
	data.Description = types.StringPointerValue(subnet.Description)
	data.IpVersion = types.StringValue(subnet.IpVersion)
	data.Name = types.StringPointerValue(subnet.Name)
	data.SubnetpoolId = types.StringPointerValue(&subnet.SubnetpoolId)
	data.VpcId = types.StringValue(subnet.VpcId)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *mgcNetworkVpcsSubnetsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := mgcNetworkVpcsSubnetsModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsServers := networkSubnets.UpdateParametersDnsNameservers{}
	for _, dns := range data.DnsNameservers {
		dnsServers = append(dnsServers, dns.ValueString())
	}
	subnetUpdateParams := networkSubnets.UpdateParameters{
		SubnetId:       data.ID.ValueString(),
		DnsNameservers: &dnsServers,
	}

	_, err := r.networkSubnets.UpdateContext(ctx, subnetUpdateParams,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubnets.UpdateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("unable to update VPC subnet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mgcNetworkVpcsSubnetsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := mgcNetworkVpcsSubnetsModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.networkSubnets.DeleteContext(ctx, networkSubnets.DeleteParameters{
		SubnetId: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubnets.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to delete VPC subnet", err.Error())
		return
	}
}

func convertIPStringToIPVersion(ipVersion string) int {
	if ipVersion == "IPv4" {
		return 4
	}
	return 6
}

func (r *mgcNetworkVpcsSubnetsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	subnetId := req.ID
	data := mgcNetworkVpcsSubnetsModel{}

	subnet, err := r.networkSubnets.GetContext(ctx, networkSubnets.GetParameters{
		SubnetId: subnetId,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubnets.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("unable to import VPC subnet", err.Error())
		return
	}

	dnsNameServers := []types.String{}
	for _, dns := range subnet.DnsNameservers {
		dnsNameServers = append(dnsNameServers, types.StringPointerValue(&dns))
	}

	data.DnsNameservers = dnsNameServers
	data.CidrBlock = types.StringPointerValue(&subnet.CidrBlock)
	data.Description = types.StringPointerValue(subnet.Description)
	data.IpVersion = types.StringValue(subnet.IpVersion)
	data.Name = types.StringPointerValue(subnet.Name)
	// data.SubnetpoolId = types.StringPointerValue(subnet.subnetPoolId)
	data.VpcId = types.StringValue(subnet.VpcId)
	data.ID = types.StringPointerValue(&subnet.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
