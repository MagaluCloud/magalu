package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	nwVPCs "magalu.cloud/lib/products/network/vpcs"
	vmImages "magalu.cloud/lib/products/virtual_machine/images"
	vmInstances "magalu.cloud/lib/products/virtual_machine/instances"
	vmMachineTypes "magalu.cloud/lib/products/virtual_machine/machine_types"

	"magalu.cloud/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vm{}
	_ resource.ResourceWithConfigure = &vm{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewVirtualMachineResource() resource.Resource {
	return &vm{}
}

// orderResource is the resource implementation.
type vm struct {
	sdkClient      *mgcSdk.Client
	vmInstances    vmInstances.Service
	vmImages       vmImages.Service
	vmMachineTypes vmMachineTypes.Service
	nwVPCs         nwVPCs.Service
}

// Configure adds the provider configured client to the resource.
func (r *vm) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.sdkClient = mgcSdk.NewClient(sdk)
	r.vmInstances = vmInstances.NewService(ctx, r.sdkClient)
	r.vmImages = vmImages.NewService(ctx, r.sdkClient)
	r.vmMachineTypes = vmMachineTypes.NewService(ctx, r.sdkClient)
	r.nwVPCs = nwVPCs.NewService(ctx, r.sdkClient)
}

// Metadata returns the resource type name.
func (r *vm) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine"
}

// vmResourceModel maps de resource schema data.
type vmResourceModel struct {
	ID           types.String       `tfsdk:"id"`
	Name         types.String       `tfsdk:"name"`
	NameIsPrefix types.Bool         `tfsdk:"name_is_prefix"`
	FinalName    types.String       `tfsdk:"final_name"`
	Region       types.String       `tfsdk:"region"`
	Env          types.String       `tfsdk:"env"`
	LastUpdated  types.String       `tfsdk:"last_updated"`
	SshKeyName   types.String       `tfsdk:"ssh_key_name"`
	State        types.String       `tfsdk:"state"`
	Status       types.String       `tfsdk:"status"`
	Network      networkVmModel     `tfsdk:"network"`
	MachineType  genericIDNameModel `tfsdk:"machine_type"`
	Image        genericIDNameModel `tfsdk:"image"`
}

type networkVmModel struct {
	IPV6              types.String       `tfsdk:"ipv6"`
	PrivateAddress    types.String       `tfsdk:"private_address"`
	PublicIpAddress   types.String       `tfsdk:"public_address"`
	DeletePublicIP    types.Bool         `tfsdk:"delete_public_ip"`
	AssociatePublicIP types.Bool         `tfsdk:"associate_public_ip"`
	VPC               genericIDNameModel `tfsdk:"vpc"`
}

// Schema defines the schema for the resource.
func (r *vm) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Optional: true,
			},
			"env": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Optional: true,
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"name_is_prefix": schema.BoolAttribute{

				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"name": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Required: true,
			},
			"final_name": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"ssh_key_name": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Required: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"image": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"name": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Optional: true,
						Computed: true,
					},
				},
			},
			"machine_type": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"name": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Optional: true,
						Computed: true,
					},
				},
			},
			"network": schema.SingleNestedAttribute{
				Required: false,
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"delete_public_ip": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"associate_public_ip": schema.BoolAttribute{
						Required: true,
					},
					"ipv6": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Computed: true,
					},
					"private_address": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Computed: true,
					},
					"public_address": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Computed: true,
					},
					"vpc": schema.SingleNestedAttribute{
						Required: false,
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								Optional: true,
								Computed: true,
							},
							"name": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
								Optional: true,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *vm) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	//do nothing
}

// Read refreshes the Terraform state with the latest data.
func (r *vm) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := &vmResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	config := NewConfig(data.Region.ValueStringPointer(), data.Env.ValueStringPointer())

	getResult, err := r.getVmStatus(data.ID.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not read VM ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	data.ID = types.StringValue(getResult.Id)
	r.setValuesFromServer(data, getResult)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *vm) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &vmResourceModel{}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config := NewConfig(plan.Region.ValueStringPointer(), plan.Env.ValueStringPointer())

	// Get image and machine type ID
	imageID, err := r.getImageID(plan.Image.Name.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not load image list",
		)
		return
	}

	machineTypeID, err := r.getMachineTypeID(plan.MachineType.Name.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not load machine-type list, unexpected error: "+err.Error(),
		)
		return
	}

	if plan.NameIsPrefix.ValueBool() && plan.FinalName.ValueString() == "" {
		plan.FinalName = types.StringValue(plan.Name.ValueString() + "-" + getRandomWords(3, "-"))
	}

	createParams := vmInstances.CreateParameters{
		Name:       plan.FinalName.ValueString(),
		SshKeyName: plan.SshKeyName.ValueString(),
		Image: vmInstances.CreateParametersImage{
			CreateParametersImage1: vmInstances.CreateParametersImage1{
				Name: plan.Image.Name.ValueString(),
			},
			CreateParametersImage0: vmInstances.CreateParametersImage0{
				Id: imageID,
			},
		},
		MachineType: vmInstances.CreateParametersMachineType{
			CreateParametersImage1: vmInstances.CreateParametersImage1{
				Name: plan.MachineType.Name.ValueString(),
			},
			CreateParametersImage0: vmInstances.CreateParametersImage0{
				Id: machineTypeID,
			},
		},
		Network: vmInstances.CreateParametersNetwork{
			AssociatePublicIp: plan.Network.AssociatePublicIP.ValueBoolPointer(),
		},
	}

	if !plan.Network.VPC.ID.IsNull() && plan.Network.VPC.ID.ValueString() != "" {
		createParams.Network.Vpc = &vmInstances.CreateParametersNetworkVpc{
			CreateParametersImage0: vmInstances.CreateParametersImage0{
				Id: plan.Network.VPC.ID.ValueString(),
			},
		}
	} else if !plan.Network.VPC.Name.IsNull() && plan.Network.VPC.Name.ValueString() != "" {
		vpcId, err := r.getVpcID(plan.Network.VPC.Name.ValueString(), config)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating vm",
				"Could not load vpc list",
			)
			return
		}
		if strings.Contains(plan.Network.VPC.Name.ValueString(), plan.Network.VPC.Name.ValueString()) {
			createParams.Network.Vpc = &vmInstances.CreateParametersNetworkVpc{
				CreateParametersImage0: vmInstances.CreateParametersImage0{
					Id: vpcId,
				},
				CreateParametersImage1: vmInstances.CreateParametersImage1{
					Name: plan.Network.VPC.Name.ValueString(),
				},
			}
		}
	}

	result, err := r.vmInstances.Create(createParams, vmInstances.CreateConfigs{Env: config.Env(), Region: config.Region()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not create virtual-machine, unexpected error: "+err.Error(),
		)
		return
	}

	getResult, err := r.getVmStatus(result.Id, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not read VM ID "+result.Id+": "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(result.Id)
	r.setValuesFromServer(plan, getResult)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vm) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := &vmResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	config := NewConfig(data.Region.ValueStringPointer(), data.Env.ValueStringPointer())

	machineTypeID, err := r.getMachineTypeID(data.MachineType.Name.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not load machine-type list, unexpected error: "+err.Error(),
		)
		return
	}

	err = r.vmInstances.Retype(vmInstances.RetypeParameters{
		Id: data.ID.ValueString(),
		MachineType: vmInstances.RetypeParametersMachineType{
			RetypeParametersMachineType0: vmInstances.RetypeParametersMachineType0{
				Id: machineTypeID,
			},
			RetypeParametersMachineType1: vmInstances.RetypeParametersMachineType1{
				Name: data.MachineType.Name.ValueString(),
			},
		},
	}, vmInstances.RetypeConfigs{Env: config.Env(), Region: config.Region()})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not update VM machine-type "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	getResult, err := r.getVmStatus(data.ID.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Error when get new vm status "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	r.setValuesFromServer(data, getResult)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vm) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	config := NewConfig(data.Region.ValueStringPointer(), data.Env.ValueStringPointer())

	err := r.vmInstances.Delete(
		vmInstances.DeleteParameters{
			DeletePublicIp: data.Network.DeletePublicIP.ValueBoolPointer(),
			Id:             data.ID.ValueString(),
		},
		vmInstances.DeleteConfigs{Env: config.Env(), Region: config.Region()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not read VM ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

}

func (r *vm) setValuesFromServer(data *vmResourceModel, server *vmInstances.GetResult) error {
	data.State = types.StringValue(server.State)
	data.Status = types.StringValue(server.Status)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	data.MachineType.ID = types.StringValue(server.MachineType.Id)
	data.Image.ID = types.StringValue(server.Image.Id)

	data.Network.VPC = genericIDNameModel{
		ID:   types.StringValue(""),
		Name: types.StringValue(""),
	}

	data.Network.IPV6 = types.StringValue("")
	data.Network.PrivateAddress = types.StringValue("")
	data.Network.PublicIpAddress = types.StringValue("")

	if server.Network.GetResultNetwork1.Ports != nil && len(*server.Network.GetResultNetwork1.Ports) > 0 {
		ports := (*server.Network.GetResultNetwork1.Ports)[0]

		data.Network.VPC = genericIDNameModel{
			ID:   types.StringValue(server.Network.GetResultNetwork1.Vpc.Id),
			Name: types.StringValue(server.Network.GetResultNetwork1.Vpc.Name),
		}

		data.Network.PrivateAddress = types.StringValue(ports.IpAddresses.PrivateIpAddress)

		if ports.IpAddresses.IpV6address != nil {
			data.Network.IPV6 = types.StringValue(*ports.IpAddresses.IpV6address)
		}

		if ports.IpAddresses.PublicIpAddress != nil {
			data.Network.PublicIpAddress = types.StringValue(*ports.IpAddresses.PublicIpAddress)
		}

	}
	return nil
}

func (r *vm) getImageID(name string, config Config) (string, error) {
	var imageID string
	imageList, err := r.vmImages.List(vmImages.ListParameters{},
		vmImages.ListConfigs{Env: config.Env(), Region: config.Region()},
	)
	if err != nil {
		return "", fmt.Errorf("could not load image list")
	}
	for _, x := range imageList.Images {
		if strings.Contains(x.Name, name) {
			imageID = x.Id
			break
		}
	}

	if imageID == "" {
		return "", fmt.Errorf("could not found image ID with name: " + name)
	}
	return imageID, nil
}

func (r *vm) getMachineTypeID(name string, config Config) (string, error) {
	var machineTypeID string
	machineTypeList, err := r.vmMachineTypes.List(vmMachineTypes.ListParameters{},
		vmMachineTypes.ListConfigs{Env: config.Env(), Region: config.Region()},
	)
	if err != nil {
		return "", fmt.Errorf("could not load machine-type list, unexpected error: " + err.Error())
	}
	for _, x := range machineTypeList.InstanceTypes {
		if x.Name == name {
			machineTypeID = x.Id
			break
		}
	}

	if machineTypeID == "" {
		return "", fmt.Errorf("could not found machine-type ID with name: " + name)
	}
	return machineTypeID, nil
}

func (r *vm) getVpcID(name string, config Config) (string, error) {
	var vpcID string
	vpcs, err := r.nwVPCs.List(nwVPCs.ListConfigs{Env: config.Env(), Region: config.Region()})
	if err != nil {
		return "", fmt.Errorf("could not load vpc list")
	}
	for _, x := range *vpcs.Vpcs {
		if strings.Contains(*x.Name, name) {
			vpcID = *x.Id
			break
		}
	}

	if vpcID == "" {
		return "", fmt.Errorf("could not found vpc ID with name: " + name)
	}
	return vpcID, nil
}

func (r *vm) getVmStatus(id string, config Config) (*vmInstances.GetResult, error) {
	getResult := &vmInstances.GetResult{}
	expand := &vmInstances.GetParametersExpand{"network"}

	duration := 5 * time.Minute
	startTime := time.Now()
	getParam := vmInstances.GetParameters{Id: id, Expand: expand}
	getConfigParam := vmInstances.GetConfigs{}
	if config != nil {
		getConfigParam = vmInstances.GetConfigs{Env: config.Env(), Region: config.Region()}
	}
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

		*getResult, err = r.vmInstances.Get(getParam, getConfigParam)
		if err != nil {
			return getResult, err
		}

		if getResult.Status == "completed" {
			return getResult, nil
		}
		time.Sleep(1 * time.Second)
	}
}
