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
	ID                types.String       `tfsdk:"id"`
	Name              types.String       `tfsdk:"name"`
	Region            types.String       `tfsdk:"region"`
	Env               types.String       `tfsdk:"env"`
	LastUpdated       types.String       `tfsdk:"last_updated"`
	SshKeyName        types.String       `tfsdk:"ssh_key_name"`
	DeletePublicIP    types.Bool         `tfsdk:"delete_public_ip"`
	State             types.String       `tfsdk:"state"`
	Status            types.String       `tfsdk:"status"`
	IPV6              types.String       `tfsdk:"ipv6"`
	PrivateAddress    types.String       `tfsdk:"private_address"`
	PublicIpAddress   types.String       `tfsdk:"public_address"`
	AssociatePublicIP types.Bool         `tfsdk:"associate_public_ip"`
	MachineType       genericIDNameModel `tfsdk:"machine_type"`
	Image             genericIDNameModel `tfsdk:"image"`
	VPC               genericIDNameModel `tfsdk:"vpc"`
}

type genericIDNameModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Schema defines the schema for the resource.
func (r *vm) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Required: true,
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
			"name": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Required: true,
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
			"delete_public_ip": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"associate_public_ip": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
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

			"image": schema.SingleNestedAttribute{
				Required: true,
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
			"machine_type": schema.SingleNestedAttribute{
				Required: true,
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
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *vm) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &vmResourceModel{}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// CHANGE IT TO DEFAULT - COMMON BLOCK
	config := &Config{
		region: new(string),
		env:    new(string),
	}
	*config.env = "prod"
	*config.region = "br-se1"

	if !plan.Region.IsNull() {
		*config.region = plan.Region.ValueString()
	}
	if !plan.Env.IsNull() {
		*config.env = plan.Env.ValueString()
	}

	// Get image and machine type ID
	var imageID string
	imageList, err := r.vmImages.List(vmImages.ListParameters{},
		vmImages.ListConfigs{Env: config.env, Region: config.region},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not load image list",
		)
		return
	}
	for _, x := range imageList.Images {
		if strings.Contains(x.Name, plan.Image.Name.ValueString()) {
			imageID = x.Id
			break
		}
	}

	if imageID == "" {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not found image ID with name: "+plan.Image.Name.ValueString(),
		)
		return
	}

	var machineTypeID string
	machineTypeList, err := r.vmMachineTypes.List(vmMachineTypes.ListParameters{},
		vmMachineTypes.ListConfigs{Env: config.env, Region: config.region},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not load machine-type list, unexpected error: "+err.Error(),
		)
		return
	}
	for _, x := range machineTypeList.InstanceTypes {
		if x.Name == plan.MachineType.Name.ValueString() {
			machineTypeID = x.Id
			break
		}
	}

	if machineTypeID == "" {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not found machine-type ID with name: "+plan.MachineType.Name.ValueString(),
		)
		return
	}

	createParams := vmInstances.CreateParameters{
		Name:       plan.Name.ValueString(),
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
	}

	createParams.Network.AssociatePublicIp = plan.AssociatePublicIP.ValueBoolPointer()

	if !plan.VPC.ID.IsNull() && plan.VPC.ID.ValueString() != "" {
		createParams.Network.Vpc = &vmInstances.CreateParametersNetworkVpc{
			CreateParametersImage0: vmInstances.CreateParametersImage0{
				Id: plan.VPC.ID.ValueString(),
			},
		}
	} else if !plan.VPC.Name.IsNull() && plan.VPC.Name.ValueString() != "" {
		vpcs, err := r.nwVPCs.List(nwVPCs.ListConfigs{Env: config.env, Region: config.region})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating vm",
				"Could not load vpc list",
			)
			return
		}
		for _, x := range *vpcs.Vpcs {
			if strings.Contains(*x.Name, plan.VPC.Name.ValueString()) {
				createParams.Network.Vpc = &vmInstances.CreateParametersNetworkVpc{
					CreateParametersImage0: vmInstances.CreateParametersImage0{
						Id: *x.Id,
					},
					CreateParametersImage1: vmInstances.CreateParametersImage1{
						Name: *x.Name,
					},
				}
				break
			}
		}
	}

	result, err := r.vmInstances.Create(createParams, vmInstances.CreateConfigs{Env: config.env, Region: config.region})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not create virtual-machine, unexpected error: "+err.Error(),
		)
		return
	}

	var getResult vmInstances.GetResult
	expand := &vmInstances.GetParametersExpand{"network"}

	//Save current taint state
	plan.ID = types.StringValue(result.Id)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.MachineType.ID = types.StringValue(machineTypeID)
	plan.Image.ID = types.StringValue(imageID)
	plan.State = types.StringValue("new")
	plan.Status = types.StringValue("creating")

	duration := 5 * time.Minute
	startTime := time.Now()
	for {
		elapsed := time.Since(startTime)
		remaining := duration - elapsed
		if remaining <= 0 {
			resp.Diagnostics.AddError(
				"Error Reading VM",
				"Timeout to read VM ID "+result.Id+": "+err.Error(),
			)
			return
		}

		getResult, err = r.vmInstances.Get(vmInstances.GetParameters{
			Id:     result.Id,
			Expand: expand,
		}, vmInstances.GetConfigs{Env: config.env, Region: config.region})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading VM",
				"Could not read created VM ID "+result.Id+": "+err.Error(),
			)
			return
		}
		if getResult.Status == "completed" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	plan.VPC = genericIDNameModel{
		ID:   types.StringValue(""),
		Name: types.StringValue(""),
	}

	plan.IPV6 = types.StringValue("")
	plan.PrivateAddress = types.StringValue("")
	plan.PublicIpAddress = types.StringValue("")

	if getResult.Network.GetResultNetwork1.Ports != nil && len(*getResult.Network.GetResultNetwork1.Ports) > 0 {
		ports := (*getResult.Network.GetResultNetwork1.Ports)[0]

		plan.VPC = genericIDNameModel{
			ID:   types.StringValue(getResult.Network.GetResultNetwork1.Vpc.Id),
			Name: types.StringValue(getResult.Network.GetResultNetwork1.Vpc.Name),
		}

		plan.PrivateAddress = types.StringValue(ports.IpAddresses.PrivateIpAddress)

		if ports.IpAddresses.IpV6address != nil {
			plan.IPV6 = types.StringValue(*ports.IpAddresses.IpV6address)
		}

		if ports.IpAddresses.PublicIpAddress != nil {
			plan.PublicIpAddress = types.StringValue(*ports.IpAddresses.PublicIpAddress)
		}

	}

	plan.State = types.StringValue(getResult.State)
	plan.Status = types.StringValue(getResult.Status)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *vm) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data vmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	config := &Config{
		region: new(string),
		env:    new(string),
	}

	*config.env = DefaultEnv(data.Env.ValueString())
	*config.region = data.Region.ValueString()

	var getResult vmInstances.GetResult
	expand := &vmInstances.GetParametersExpand{"network"}

	getResult, err := r.vmInstances.Get(vmInstances.GetParameters{
		Id:     data.ID.ValueString(),
		Expand: expand,
	}, vmInstances.GetConfigs{Env: config.env, Region: config.region})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not read VM ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	//Save current taint state
	data.ID = types.StringValue(getResult.Id)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	data.MachineType.ID = types.StringValue(getResult.MachineType.Id)
	data.Image.ID = types.StringValue(getResult.Image.Id)

	if getResult.Network.GetResultNetwork1.Ports != nil && len(*getResult.Network.GetResultNetwork1.Ports) > 0 {
		ports := (*getResult.Network.GetResultNetwork1.Ports)[0]

		data.VPC = genericIDNameModel{
			ID:   types.StringValue(getResult.Network.GetResultNetwork1.Vpc.Id),
			Name: types.StringValue(getResult.Network.GetResultNetwork1.Vpc.Name),
		}

		data.PrivateAddress = types.StringValue(ports.IpAddresses.PrivateIpAddress)

		if ports.IpAddresses.IpV6address != nil {
			data.IPV6 = types.StringValue(*ports.IpAddresses.IpV6address)
		}

		if ports.IpAddresses.PublicIpAddress != nil {
			data.PublicIpAddress = types.StringValue(*ports.IpAddresses.PublicIpAddress)
		}

	}

	data.State = types.StringValue(getResult.State)
	data.Status = types.StringValue(getResult.Status)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vm) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data vmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	config := &Config{
		region: new(string),
		env:    new(string),
	}

	*config.env = DefaultEnv(data.Env.ValueString())
	*config.region = data.Region.ValueString()

	var machineTypeID string
	machineTypeList, err := r.vmMachineTypes.List(vmMachineTypes.ListParameters{},
		vmMachineTypes.ListConfigs{Env: config.env, Region: config.region},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not load machine-type list, unexpected error: "+err.Error(),
		)
		return
	}
	for _, x := range machineTypeList.InstanceTypes {
		if x.Name == data.MachineType.Name.ValueString() {
			machineTypeID = x.Id
			break
		}
	}

	if machineTypeID == "" {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not found machine-type ID with name: "+data.MachineType.Name.ValueString(),
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
	}, vmInstances.RetypeConfigs{Env: config.env, Region: config.region})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not update VM machine-type "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	var getResult vmInstances.GetResult
	duration := 5 * time.Minute
	startTime := time.Now()
	for {
		elapsed := time.Since(startTime)
		remaining := duration - elapsed
		if remaining <= 0 {
			resp.Diagnostics.AddError(
				"Error Reading VM",
				"Timeout to read VM ID "+data.ID.ValueString()+": "+err.Error(),
			)
			return
		}

		getResult, err = r.vmInstances.Get(vmInstances.GetParameters{
			Id: data.ID.ValueString(),
		}, vmInstances.GetConfigs{Env: config.env, Region: config.region})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading VM",
				"Could not read created VM ID "+data.ID.ValueString()+": "+err.Error(),
			)
			return
		}
		if getResult.Status == "completed" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	data.MachineType.ID = types.StringValue(machineTypeID)
	data.State = types.StringValue(getResult.State)
	data.Status = types.StringValue(getResult.Status)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vm) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	config := &Config{
		region: new(string),
		env:    new(string),
	}

	*config.env = DefaultEnv(data.Env.ValueString())
	*config.region = data.Region.ValueString()

	err := r.vmInstances.Delete(
		vmInstances.DeleteParameters{
			DeletePublicIp: data.DeletePublicIP.ValueBoolPointer(),
			Id:             data.ID.ValueString(),
		},
		vmInstances.DeleteConfigs{Env: config.env, Region: config.region})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not read VM ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

}
