package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
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
			fmt.Sprintf("Expected *apiKey, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.sdkClient = mgcSdk.NewClient(sdk)
	r.vmInstances = vmInstances.NewService(ctx, r.sdkClient)
	r.vmImages = vmImages.NewService(ctx, r.sdkClient)
	r.vmMachineTypes = vmMachineTypes.NewService(ctx, r.sdkClient)

}

// Metadata returns the resource type name.
func (r *vm) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine"
}

// vmResourceModel maps de resource schema data.
type vmResourceModel struct {
	ID             types.String       `tfsdk:"id"`
	Name           types.String       `tfsdk:"name"`
	LastUpdated    types.String       `tfsdk:"last_updated"`
	SshKeyName     types.String       `tfsdk:"ssh_key_name"`
	DeletePublicIP types.Bool         `tfsdk:"delete_public_ip"`
	MachineType    genericIDNameModel `tfsdk:"machine_type"`
	Image          genericIDNameModel `tfsdk:"image"`
	Netowrk        networkModel       `tfsdk:"network"`
}

type genericIDNameModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type networkModel struct {
	AssociatePublicIP types.Bool          `tfsdk:"associate_public_ip"`
	VPC               *genericIDNameModel `tfsdk:"vpc"`
}

// Schema defines the schema for the resource.
func (r *vm) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"ssh_key_name": schema.StringAttribute{
				Required: true,
			},
			"delete_public_ip": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"machine_type": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"image": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"network": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"associate_public_ip": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"vpc": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Required: true,
							},
							"name": schema.StringAttribute{
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *vm) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Default formula: PLAN > GEN REQUEST > RUN > RESPONSE

	var plan vmResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := &Config{
		region: new(string),
		env:    new(string),
	}

	*config.env = "prod"
	*config.region = "br-se1"

	// // TODO - VALIDATE EMPTY STRING ""
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
		}
	}

	if imageID == "" {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not found image ID",
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
		}
	}

	if machineTypeID == "" {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not found machine-type ID, unexpected error: "+err.Error(),
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
		Network: vmInstances.CreateParametersNetwork{
			AssociatePublicIp: plan.Netowrk.AssociatePublicIP.ValueBoolPointer(),
		},
	}

	if plan.Netowrk.VPC != nil && !plan.Netowrk.VPC.ID.IsNull() {
		createParams.Network.Vpc = &vmInstances.CreateParametersNetworkVpc{
			CreateParametersImage1: vmInstances.CreateParametersImage1{
				Name: plan.Netowrk.VPC.Name.ValueString(),
			},
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

	plan.ID = types.StringValue(result.Id)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.MachineType.ID = types.StringValue(machineTypeID)
	plan.Image.ID = types.StringValue(imageID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *vm) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vm) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vm) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
