package resources

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	bws "github.com/geffersonFerraz/brazilian-words-sorter"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkNetworkVPCs "magalu.cloud/lib/products/network/vpc"
	sdkVmImages "magalu.cloud/lib/products/virtual_machine/images"
	sdkVmInstances "magalu.cloud/lib/products/virtual_machine/instances"
	sdkVmMachineTypes "magalu.cloud/lib/products/virtual_machine/machine_types"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	tfutil "magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var (
	_ resource.Resource              = &vmInstances{}
	_ resource.ResourceWithConfigure = &vmInstances{}
)

func NewVirtualMachineInstancesResource() resource.Resource {
	return &vmInstances{}
}

type vmInstances struct {
	sdkClient      *mgcSdk.Client
	vmInstances    sdkVmInstances.Service
	vmImages       sdkVmImages.Service
	vmMachineTypes sdkVmMachineTypes.Service
	nwVPCs         sdkNetworkVPCs.Service
}

func (r *vmInstances) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_instances"
}

func (r *vmInstances) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	r.sdkClient, err = client.NewSDKClient(req)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.vmInstances = sdkVmInstances.NewService(ctx, r.sdkClient)
	r.vmImages = sdkVmImages.NewService(ctx, r.sdkClient)
	r.vmMachineTypes = sdkVmMachineTypes.NewService(ctx, r.sdkClient)
	r.nwVPCs = sdkNetworkVPCs.NewService(ctx, r.sdkClient)
}

type vmInstancesResourceModel struct {
	ID           types.String                `tfsdk:"id"`
	Name         types.String                `tfsdk:"name"`
	NameIsPrefix types.Bool                  `tfsdk:"name_is_prefix"`
	FinalName    types.String                `tfsdk:"final_name"`
	UpdatedAt    types.String                `tfsdk:"updated_at"`
	CreatedAt    types.String                `tfsdk:"created_at"`
	SshKeyName   types.String                `tfsdk:"ssh_key_name"`
	State        types.String                `tfsdk:"state"`
	Status       types.String                `tfsdk:"status"`
	Network      networkVmInstancesModel     `tfsdk:"network"`
	MachineType  vmInstancesMachineTypeModel `tfsdk:"machine_type"`
	Image        tfutil.GenericIDNameModel   `tfsdk:"image"`
}

type networkVmInstancesModel struct {
	IPV6              types.String                          `tfsdk:"ipv6"`
	PrivateAddress    types.String                          `tfsdk:"private_address"`
	PublicIpAddress   types.String                          `tfsdk:"public_address"`
	DeletePublicIP    types.Bool                            `tfsdk:"delete_public_ip"`
	AssociatePublicIP types.Bool                            `tfsdk:"associate_public_ip"`
	VPC               *tfutil.GenericIDNameModel            `tfsdk:"vpc"`
	Interface         *vmInstancesNetworkSecurityGroupModel `tfsdk:"interface"`
}

type vmInstancesNetworkSecurityGroupModel struct {
	SecurityGroups []tfutil.GenericIDModel `tfsdk:"security_groups"`
}

type vmInstancesMachineTypeModel struct {
	ID    types.String `tfsdk:"id"`
	Disk  types.Number `tfsdk:"disk"`
	Name  types.String `tfsdk:"name"`
	RAM   types.Number `tfsdk:"ram"`
	VCPUs types.Number `tfsdk:"vcpus"`
}

func (r *vmInstances) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "Operations with instances, including create, delete, start, stop, reboot and other actions."
	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the virtual machine instance.",
				Computed:    true,
			},
			"name_is_prefix": schema.BoolAttribute{
				Description: "Indicates whether the provided name is a prefix or the exact name of the virtual machine instance.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the virtual machine instance.",
				Required:    true,
			},
			"final_name": schema.StringAttribute{
				Description: "The final name of the virtual machine instance after applying any naming conventions or modifications.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The timestamp when the virtual machine instance was last updated.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the virtual machine instance was created.",
				Computed:    true,
			},
			"ssh_key_name": schema.StringAttribute{
				Description: "The name of the SSH key associated with the virtual machine instance. If the image is Windows, this field is not used.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Optional: true,
			},
			"state": schema.StringAttribute{
				Description: "The current state of the virtual machine instance.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the virtual machine instance.",
				Computed:    true,
			},
			"image": schema.SingleNestedAttribute{
				Description: "The image used to create the virtual machine instance.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The unique identifier of the image.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the image.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Required: true,
					},
				},
			},
			"machine_type": schema.SingleNestedAttribute{
				Description: "The machine type of the virtual machine instance.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The unique identifier of the machine type.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the machine type.",
						Required:    true,
					},
					"disk": schema.NumberAttribute{
						Description: "The disk size of the machine type.",
						Computed:    true,
					},
					"ram": schema.NumberAttribute{
						Description: "The RAM size of the machine type.",
						Computed:    true,
					},
					"vcpus": schema.NumberAttribute{
						Description: "The number of virtual CPUs of the machine type.",
						Computed:    true,
					},
				},
				Validators: []validator.Object{
					objectvalidator.IsRequired(),
				},
			},
			"network": schema.SingleNestedAttribute{
				Description: "The network configuration of the virtual machine instance.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"delete_public_ip": schema.BoolAttribute{
						Description: "Indicates whether to delete the public IP address associated with the virtual machine instance.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
						Default:  booldefault.StaticBool(true),
						Optional: true,
						Computed: true,
					},
					"associate_public_ip": schema.BoolAttribute{
						Description: "Indicates whether to associate a public IP address with the virtual machine instance.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"ipv6": schema.StringAttribute{
						Description: "The IPv6 address of the virtual machine instance.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Computed: true,
					},
					"private_address": schema.StringAttribute{
						Description: "The private IP address of the virtual machine instance.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Computed: true,
					},
					"public_address": schema.StringAttribute{
						Description: "The public IP address of the virtual machine instance.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Computed: true,
					},
					"vpc": schema.SingleNestedAttribute{
						Description: "The VPC (Virtual Private Cloud) associated with the virtual machine instance.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: "The unique identifier of the VPC.",
								Computed:    true,
								Optional:    true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"name": schema.StringAttribute{
								Description: "The name of the VPC.",
								Computed:    true,
							},
						},
					},
					"interface": schema.SingleNestedAttribute{
						Description: "The network interface configuration of the virtual machine instance.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"security_groups": schema.ListNestedAttribute{
								Description: "The security groups associated with the network interface.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "The unique identifier of the security group.",
											Optional:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.RequiresReplace(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *vmInstances) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := vmInstancesResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	getResult, err := r.getVmStatus(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading VM", err.Error())
		return
	}
	outputModel := addVMLocalDataToModel(*toVirtualMachineState(getResult), state.Network.AssociatePublicIP.ValueBool(), state.NameIsPrefix.ValueBool(), state.Network.DeletePublicIP.ValueBool(), state.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, outputModel)...)
}

func (r *vmInstances) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state := vmInstancesResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.FinalName = types.StringValue(state.Name.ValueString())
	if state.NameIsPrefix.ValueBool() {
		bwords := bws.BrazilianWords(3, "-")
		state.FinalName = types.StringValue(state.Name.ValueString() + "-" + bwords.Sort())
	}

	createParams := sdkVmInstances.CreateParameters{
		Name:       state.FinalName.ValueString(),
		SshKeyName: state.SshKeyName.ValueStringPointer(),
		Image: sdkVmInstances.CreateParametersImage{
			Name: state.Image.Name.ValueStringPointer(),
		},
		MachineType: sdkVmInstances.CreateParametersMachineType{
			Name: state.MachineType.Name.ValueStringPointer(),
		},
		Network: &sdkVmInstances.CreateParametersNetwork{},
	}

	if state.Network.VPC != nil && state.Network.VPC.ID.ValueString() != "" {
		createParams.Network.Vpc = &sdkVmInstances.CreateParametersNetworkVpc{
			Id: state.Network.VPC.ID.ValueString(),
		}
	}
	if state.Network.Interface != nil && len(state.Network.Interface.SecurityGroups) > 0 {
		network := sdkVmInstances.CreateParametersNetwork{}
		network.Interface = &sdkVmInstances.CreateParametersNetworkInterface{}
		network.Interface.SecurityGroups = &sdkVmInstances.CreateParametersNetworkInterfaceSecurityGroups{}
		items := []sdkVmInstances.CreateParametersNetworkInterfaceSecurityGroupsItem{}
		for _, sg := range state.Network.Interface.SecurityGroups {
			items = append(items, sdkVmInstances.CreateParametersNetworkInterfaceSecurityGroupsItem{
				Id: sg.ID.ValueString(),
			})
		}
		vmInstancesNetworkInterfaceSecurityGroups := sdkVmInstances.CreateParametersNetworkInterfaceSecurityGroups(items)
		createParams.Network.Interface = &sdkVmInstances.CreateParametersNetworkInterface{}
		createParams.Network.Interface.SecurityGroups = &vmInstancesNetworkInterfaceSecurityGroups
	}

	createParams.Network.AssociatePublicIp = state.Network.AssociatePublicIP.ValueBoolPointer()
	createdVM, err := r.vmInstances.CreateContext(ctx, createParams, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstances.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Error creating vm", err.Error())
		return
	}

	getResult, err := r.getVmStatus(ctx, createdVM.Id)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading VM", err.Error())
		return
	}
	outputModel := addVMLocalDataToModel(*toVirtualMachineState(getResult), state.Network.AssociatePublicIP.ValueBool(), state.NameIsPrefix.ValueBool(), state.Network.DeletePublicIP.ValueBool(), state.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, outputModel)...)
}

func (r *vmInstances) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	newState := vmInstancesResourceModel{}
	currState := &vmInstancesResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &newState)...)
	resp.Diagnostics.Append(req.State.Get(ctx, currState)...)

	if !currState.Name.Equal(newState.Name) {
		newState.FinalName = types.StringValue(newState.Name.ValueString())
		if newState.NameIsPrefix.ValueBool() {
			bwords := bws.BrazilianWords(3, "-")
			newState.FinalName = types.StringValue(newState.Name.ValueString() + "-" + bwords.Sort())
		}
		err := r.vmInstances.RenameContext(ctx, sdkVmInstances.RenameParameters{
			Id:   currState.ID.ValueString(),
			Name: newState.FinalName.ValueString(),
		}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstances.RenameConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Error to rename vm", err.Error())
			return
		}
		_, _ = r.getVmStatus(ctx, currState.ID.ValueString())
	}

	machineType, err := r.getMachineTypeID(ctx, newState.MachineType.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error to retype vm", err.Error())
		return
	}
	if !currState.MachineType.ID.Equal(machineType.ID) {
		err = r.vmInstances.RetypeContext(ctx, sdkVmInstances.RetypeParameters{
			Id: currState.ID.ValueString(),
			MachineType: sdkVmInstances.RetypeParametersMachineType{
				Id: machineType.ID.ValueString(),
			},
		}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstances.RetypeConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Error on Update VM", err.Error())
			return
		}
	}

	getResult, err := r.getVmStatus(ctx, currState.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading VM", err.Error())
		return
	}
	outputModel := addVMLocalDataToModel(*toVirtualMachineState(getResult), newState.Network.AssociatePublicIP.ValueBool(), newState.NameIsPrefix.ValueBool(), newState.Network.DeletePublicIP.ValueBool(), newState.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, outputModel)...)
}

func (r *vmInstances) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vmInstancesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	err := r.vmInstances.DeleteContext(ctx,
		sdkVmInstances.DeleteParameters{
			DeletePublicIp: data.Network.DeletePublicIP.ValueBoolPointer(),
			Id:             data.ID.ValueString(),
		},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstances.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting VM", err.Error())
		return
	}
	r.checkVmIsNotFound(ctx, data.ID.ValueString())
}

func toVirtualMachineState(server *sdkVmInstances.GetResult) *vmInstancesResourceModel {
	data := vmInstancesResourceModel{}
	data.ID = types.StringValue(server.Id)
	data.FinalName = types.StringValue(*server.Name)
	data.State = types.StringValue(server.State)
	data.Status = types.StringValue(server.Status)
	data.MachineType.ID = types.StringValue(server.MachineType.Id)
	data.MachineType.Name = types.StringValue(*server.MachineType.Name)
	data.MachineType.Disk = types.NumberValue(new(big.Float).SetInt64(int64(*server.MachineType.Disk)))
	data.MachineType.RAM = types.NumberValue(new(big.Float).SetInt64(int64(*server.MachineType.Ram)))
	data.MachineType.VCPUs = types.NumberValue(new(big.Float).SetInt64(int64(*server.MachineType.Vcpus)))
	data.Image.Name = types.StringValue(*server.Image.Name)
	data.Image.ID = types.StringValue(server.Image.Id)
	data.SshKeyName = types.StringValue(*server.SshKeyName)
	data.CreatedAt = types.StringValue(server.CreatedAt)
	data.UpdatedAt = types.StringValue(*server.UpdatedAt)

	if server.Network.Vpc != nil {
		vpc := tfutil.GenericIDNameModel{}
		vpc.ID = types.StringValue(server.Network.Vpc.Id)
		vpc.Name = types.StringValue(server.Network.Vpc.Name)
		data.Network.VPC = &vpc
	}

	data.Network.IPV6 = types.StringValue("")
	data.Network.PrivateAddress = types.StringValue("")
	data.Network.PublicIpAddress = types.StringValue("")

	if server.Network.Ports != nil && len(*server.Network.Ports) > 0 {
		ports := (*server.Network.Ports)[0]
		data.Network.PrivateAddress = types.StringValue(ports.IpAddresses.PrivateIpAddress)
		if ports.IpAddresses.IpV6address != nil {
			data.Network.IPV6 = types.StringValue(*ports.IpAddresses.IpV6address)
		}
		if ports.IpAddresses.PublicIpAddress != nil {
			data.Network.PublicIpAddress = types.StringValue(*ports.IpAddresses.PublicIpAddress)
		}
	}

	return &data
}

func addVMLocalDataToModel(model vmInstancesResourceModel, assPublicIP bool, nmIsPrefix bool, dellPublicIP bool, originalName string) *vmInstancesResourceModel {
	model.Name = types.StringValue(originalName)
	model.Network.AssociatePublicIP = types.BoolValue(assPublicIP)
	model.NameIsPrefix = types.BoolValue(nmIsPrefix)
	model.Network.DeletePublicIP = types.BoolValue(dellPublicIP)
	return &model
}

func (r *vmInstances) getMachineTypeID(ctx context.Context, name string) (*vmInstancesMachineTypeModel, error) {
	machineType := vmInstancesMachineTypeModel{}
	machineTypeList, err := r.vmMachineTypes.ListContext(ctx, sdkVmMachineTypes.ListParameters{}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmMachineTypes.ListConfigs{}))
	if err != nil {
		return nil, fmt.Errorf("could not load machine-type list, unexpected error: %w", err)
	}

	for _, x := range machineTypeList.MachineTypes {
		if x.Name == name {
			machineType.Disk = types.NumberValue(new(big.Float).SetInt64(int64(x.Disk)))
			machineType.ID = types.StringValue(x.Id)
			machineType.Name = types.StringValue(x.Name)
			machineType.RAM = types.NumberValue(new(big.Float).SetInt64(int64(x.Ram)))
			machineType.VCPUs = types.NumberValue(new(big.Float).SetInt64(int64(x.Vcpus)))
			break
		}
	}

	if machineType.ID.ValueString() == "" {
		return nil, fmt.Errorf("could not found machine-type ID with name: %s", name)
	}
	return &machineType, nil
}

func (r *vmInstances) getVmStatus(ctx context.Context, id string) (*sdkVmInstances.GetResult, error) {
	getResult := &sdkVmInstances.GetResult{}
	expand := &sdkVmInstances.GetParametersExpand{"network", "machine-types", "image"}

	duration := 5 * time.Minute
	startTime := time.Now()
	getParam := sdkVmInstances.GetParameters{Id: id, Expand: expand}
	var err error
	for {
		elapsed := time.Since(startTime)
		remaining := duration - elapsed
		if remaining <= 0 {
			if getResult.Status != "" {
				return getResult, nil
			}
			return getResult, fmt.Errorf("timeout to read VM ID")
		}

		*getResult, err = r.vmInstances.GetContext(ctx, getParam, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstances.GetConfigs{}))
		if err != nil {
			return getResult, err
		}
		if getResult.Status == "completed" {
			return getResult, nil
		}
		time.Sleep(3 * time.Second)
	}
}

func (r *vmInstances) checkVmIsNotFound(ctx context.Context, id string) {
	getResult := &sdkVmInstances.GetResult{}
	duration := 5 * time.Minute
	startTime := time.Now()
	getParam := sdkVmInstances.GetParameters{Id: id}
	var err error
	for {
		elapsed := time.Since(startTime)
		remaining := duration - elapsed
		if remaining <= 0 {
			if getResult.Status != "" {
				return
			}
			return
		}
		*getResult, err = r.vmInstances.GetContext(ctx, getParam, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstances.GetConfigs{}))
		if err != nil {
			return
		}
		time.Sleep(3 * time.Second)
	}
}

func (r *vmInstances) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Invalid import ID", "The ID must be provided")
		return
	}
	getResult, err := r.getVmStatus(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading VM", err.Error())
		return
	}
	data := toVirtualMachineState(getResult)
	data.Name = data.FinalName
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
