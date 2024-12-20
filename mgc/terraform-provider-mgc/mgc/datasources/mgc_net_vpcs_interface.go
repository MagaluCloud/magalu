package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkInterfaces "magalu.cloud/lib/products/network/ports"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkVPCInterfaceModel struct {
	CreatedAt             types.String                        `tfsdk:"created_at"`
	Description           types.String                        `tfsdk:"description"`
	Id                    types.String                        `tfsdk:"id"`
	IpAddress             []NetworkVPCInterfaceIpAddressModel `tfsdk:"ip_address"`
	IsAdminStateUp        types.Bool                          `tfsdk:"is_admin_state_up"`
	IsPortSecurityEnabled types.Bool                          `tfsdk:"is_port_security_enabled"`
	Name                  types.String                        `tfsdk:"name"`
	PublicIp              []NetworkVPCInterfacePublicIpModel  `tfsdk:"public_ip"`
	SecurityGroups        []types.String                      `tfsdk:"security_groups"`
	Updated               types.String                        `tfsdk:"updated"`
	VpcId                 types.String                        `tfsdk:"vpc_id"`
}

type NetworkVPCInterfaceIpAddressModel struct {
	Ethertype types.String `tfsdk:"ethertype"`
	IpAddress types.String `tfsdk:"ip_address"`
	SubnetId  types.String `tfsdk:"subnet_id"`
}

type NetworkVPCInterfacePublicIpModel struct {
	PublicIp   types.String `tfsdk:"public_ip"`
	PublicIpId types.String `tfsdk:"public_ip_id"`
}

type NetworkVPCInterfaceDatasource struct {
	sdkClient         *mgcSdk.Client
	networkInterfaces networkInterfaces.Service
}

func NewDataSourceNetworkVPCInterface() datasource.DataSource {
	return &NetworkVPCInterfaceDatasource{}
}

func (r *NetworkVPCInterfaceDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_vpcs_interface"
}

func (r *NetworkVPCInterfaceDatasource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network VPC Interface",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp of the VPC interface",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the VPC interface",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the VPC interface",
				Required:    true,
			},
			"ip_address": schema.ListNestedAttribute{
				Description: "The IP addresses associated with the VPC interface",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ethertype": schema.StringAttribute{
							Description: "The ethertype of the IP address",
							Computed:    true,
						},
						"ip_address": schema.StringAttribute{
							Description: "The IP address",
							Computed:    true,
						},
						"subnet_id": schema.StringAttribute{
							Description: "The subnet ID",
							Computed:    true,
						},
					},
				},
			},
			"is_admin_state_up": schema.BoolAttribute{
				Description: "Administrative state of the VPC interface",
				Computed:    true,
			},
			"is_port_security_enabled": schema.BoolAttribute{
				Description: "Port security status of the VPC interface",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the VPC interface",
				Computed:    true,
			},
			"public_ip": schema.ListNestedAttribute{
				Description: "Public IP configuration of the VPC interface",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"public_ip": schema.StringAttribute{
							Description: "The public IP address",
							Computed:    true,
						},
						"public_ip_id": schema.StringAttribute{
							Description: "The public IP ID",
							Computed:    true,
						},
					},
				},
			},
			"security_groups": schema.ListAttribute{
				Description: "List of security group IDs",
				Computed:    true,
				ElementType: types.StringType,
			},
			"updated": schema.StringAttribute{
				Description: "Last update timestamp of the VPC interface",
				Computed:    true,
			},
			"vpc_id": schema.StringAttribute{
				Description: "ID of the VPC this interface belongs to",
				Computed:    true,
			},
		},
	}
}

func (r *NetworkVPCInterfaceDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.networkInterfaces = networkInterfaces.NewService(ctx, r.sdkClient)
}

func (r *NetworkVPCInterfaceDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &NetworkVPCInterfaceModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpcInterface, err := r.networkInterfaces.GetContext(ctx, networkInterfaces.GetParameters{
		PortId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkInterfaces.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get VPC interface", err.Error())
		return
	}

	data.Name = types.StringPointerValue(vpcInterface.Name)
	data.VpcId = types.StringPointerValue(vpcInterface.VpcId)
	data.IpAddress = []NetworkVPCInterfaceIpAddressModel{}
	if vpcInterface.IpAddress == nil {
		for _, ipAddress := range *vpcInterface.IpAddress {
			data.IpAddress = append(data.IpAddress, NetworkVPCInterfaceIpAddressModel{
				Ethertype: types.StringPointerValue(ipAddress.Ethertype),
				IpAddress: types.StringValue(ipAddress.IpAddress),
				SubnetId:  types.StringValue(ipAddress.SubnetId),
			})
		}
	}
	data.IsAdminStateUp = types.BoolPointerValue(vpcInterface.IsAdminStateUp)
	data.IsPortSecurityEnabled = types.BoolPointerValue(vpcInterface.IsPortSecurityEnabled)
	data.PublicIp = []NetworkVPCInterfacePublicIpModel{}
	if vpcInterface.PublicIp == nil {
		for _, publicIp := range *vpcInterface.PublicIp {
			data.PublicIp = append(data.PublicIp, NetworkVPCInterfacePublicIpModel{
				PublicIp:   types.StringPointerValue(publicIp.PublicIp),
				PublicIpId: types.StringPointerValue(publicIp.PublicIpId),
			})
		}
	}
	data.SecurityGroups = []types.String{}
	if vpcInterface.SecurityGroups != nil {
		for _, securityGroup := range *vpcInterface.SecurityGroups {
			data.SecurityGroups = append(data.SecurityGroups, types.StringValue(securityGroup))
		}
	}
	data.CreatedAt = types.StringPointerValue(vpcInterface.CreatedAt)
	data.Updated = types.StringPointerValue(vpcInterface.Updated)
	data.Description = types.StringPointerValue(vpcInterface.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
