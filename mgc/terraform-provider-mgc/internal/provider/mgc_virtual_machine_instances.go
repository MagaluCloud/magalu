package provider

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	bws "github.com/geffersonFerraz/brazilian-words-sorter"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkNetworkVPCs "magalu.cloud/lib/products/network/vpcs"
	sdkVmImages "magalu.cloud/lib/products/virtual_machine/images"
	sdkVmInstances "magalu.cloud/lib/products/virtual_machine/instances"
	sdkVmMachineTypes "magalu.cloud/lib/products/virtual_machine/machine_types"

	"magalu.cloud/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vmInstances{}
	_ resource.ResourceWithConfigure = &vmInstances{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewVirtualMachineInstancesResource() resource.Resource {
	return &vmInstances{}
}

// orderResource is the resource implementation.
type vmInstances struct {
	sdkClient      *mgcSdk.Client
	vmInstances    sdkVmInstances.Service
	vmImages       sdkVmImages.Service
	vmMachineTypes sdkVmMachineTypes.Service
	nwVPCs         sdkNetworkVPCs.Service
}

// Metadata returns the resource type name.
func (r *vmInstances) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_instances"
}

// Configure adds the provider configured client to the resource.
func (r *vmInstances) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.vmInstances = sdkVmInstances.NewService(ctx, r.sdkClient)
	r.vmImages = sdkVmImages.NewService(ctx, r.sdkClient)
	r.vmMachineTypes = sdkVmMachineTypes.NewService(ctx, r.sdkClient)
	r.nwVPCs = sdkNetworkVPCs.NewService(ctx, r.sdkClient)
}

// vmInstancesResourceModel maps de resource schema data.
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
	Image        genericIDNameModel          `tfsdk:"image"`
}

type networkVmInstancesModel struct {
	IPV6              types.String       `tfsdk:"ipv6"`
	PrivateAddress    types.String       `tfsdk:"private_address"`
	PublicIpAddress   types.String       `tfsdk:"public_address"`
	DeletePublicIP    types.Bool         `tfsdk:"delete_public_ip"`
	AssociatePublicIP types.Bool         `tfsdk:"associate_public_ip"`
	VPC               genericIDNameModel `tfsdk:"vpc"`
}

type vmInstancesMachineTypeModel struct {
	ID    types.String `tfsdk:"id"`
	Disk  types.Number `tfsdk:"disk"`
	Name  types.String `tfsdk:"name"`
	RAM   types.Number `tfsdk:"ram"`
	VCPUs types.Number `tfsdk:"vcpus"`
}

// Schema defines the schema for the resource.
func (r *vmInstances) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "Operations with instances, including create, delete, start, stop, reboot and other actions."
	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
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
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
					"disk": schema.NumberAttribute{
						Computed: true,
					},
					"ram": schema.NumberAttribute{
						Computed: true,
					},
					"vcpus": schema.NumberAttribute{
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

func (r *vmInstances) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	resp.RequiresReplace = []path.Path{path.Root("network").AtName("associate_public_ip")}
}
func (r *vmInstances) ModifyPlanResponse(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// nothing
}

// Read refreshes the Terraform state with the latest data.
func (r *vmInstances) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := &vmInstancesResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	getResult, err := r.getVmStatus(data.ID.ValueString())
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
func (r *vmInstances) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &vmInstancesResourceModel{}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	machineType, err := r.getMachineTypeID(plan.MachineType.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not load machine-type list, unexpected error: "+err.Error(),
		)
		return
	}

	bwords := bws.BrazilianWords(3, "-")
	if plan.NameIsPrefix.ValueBool() && plan.FinalName.ValueString() == "" {
		plan.FinalName = types.StringValue(plan.Name.ValueString() + "-" + bwords.Sort())
	}

	createParams := sdkVmInstances.CreateParameters{
		Name:       plan.FinalName.ValueString(),
		SshKeyName: plan.SshKeyName.ValueString(),
		Image: sdkVmInstances.CreateParametersImage{
			CreateParametersImage1: sdkVmInstances.CreateParametersImage1{
				Name: plan.Image.Name.ValueString(),
			},
		},
		MachineType: sdkVmInstances.CreateParametersMachineType{
			CreateParametersImage1: sdkVmInstances.CreateParametersImage1{
				Name: plan.MachineType.Name.ValueString(),
			},
		},
		Network: &sdkVmInstances.CreateParametersNetwork{
			AssociatePublicIp: plan.Network.AssociatePublicIP.ValueBoolPointer(),
		},
	}

	if !plan.Network.VPC.ID.IsNull() && plan.Network.VPC.ID.ValueString() != "" {
		createParams.Network.Vpc = &sdkVmInstances.CreateParametersNetworkVpc{
			CreateParametersImage1: sdkVmInstances.CreateParametersImage1{
				Name: plan.Network.VPC.Name.ValueString(),
			},
		}
	} else if !plan.Network.VPC.Name.IsNull() && plan.Network.VPC.Name.ValueString() != "" {
		// vpcId, err := r.getVpcID(plan.Network.VPC.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating vm",
				"Could not load vpc list",
			)
			return
		}
		if strings.Contains(plan.Network.VPC.Name.ValueString(), plan.Network.VPC.Name.ValueString()) {
			createParams.Network.Vpc = &sdkVmInstances.CreateParametersNetworkVpc{
				CreateParametersImage1: sdkVmInstances.CreateParametersImage1{
					Name: plan.Network.VPC.Name.ValueString(),
				},
			}
		}
	}

	result, err := r.vmInstances.Create(createParams, sdkVmInstances.CreateConfigs{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not create virtual-machine, unexpected error: "+err.Error(),
		)
		return
	}

	getResult, err := r.getVmStatus(result.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not read VM ID "+result.Id+": "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(result.Id)

	r.setValuesFromServer(plan, getResult)
	r.setValuesFromMachineType(plan, machineType)

	plan.CreatedAt = types.StringValue(time.Now().Format(time.RFC850))
	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vmInstances) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := &vmInstancesResourceModel{}
	currState := &vmInstancesResourceModel{}
	req.State.Get(ctx, currState)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	machineType, err := r.getMachineTypeID(data.MachineType.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not found machine-type or load machine-type list, unexpected error: "+err.Error(),
		)
		return
	}

	if !currState.MachineType.ID.Equal(machineType.ID) {
		err = r.vmInstances.Retype(sdkVmInstances.RetypeParameters{
			Id: data.ID.ValueString(),
			MachineType: sdkVmInstances.RetypeParametersMachineType{
				RetypeParametersMachineType1: sdkVmInstances.RetypeParametersMachineType1{
					Name: data.MachineType.Name.ValueString(),
				},
			},
		}, sdkVmInstances.RetypeConfigs{})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error on Update VM",
				"Could not update VM machine-type "+data.ID.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	data.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))
	getResult, err := r.getVmStatus(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Error when get new vm status "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	r.setValuesFromServer(data, getResult)
	r.setValuesFromMachineType(data, machineType)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vmInstances) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vmInstancesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	err := r.vmInstances.Delete(
		sdkVmInstances.DeleteParameters{
			DeletePublicIp: data.Network.DeletePublicIP.ValueBoolPointer(),
			Id:             data.ID.ValueString(),
		},
		sdkVmInstances.DeleteConfigs{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting VM",
			"Could not delete VM ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *vmInstances) setValuesFromServer(data *vmInstancesResourceModel, server *sdkVmInstances.GetResult) {
	data.State = types.StringValue(server.State)
	data.Status = types.StringValue(server.Status)
	data.MachineType.ID = types.StringValue(server.MachineType.Id)
	data.MachineType.Name = types.StringValue(server.MachineType.Name)
	data.MachineType.Disk = types.NumberValue(new(big.Float).SetInt64(int64(server.MachineType.Disk)))
	data.MachineType.RAM = types.NumberValue(new(big.Float).SetInt64(int64(server.MachineType.Ram)))

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
}

func (r *vmInstances) getMachineTypeID(name string) (*vmInstancesMachineTypeModel, error) {
	machineType := vmInstancesMachineTypeModel{}
	machineTypeList, err := r.vmMachineTypes.List(sdkVmMachineTypes.ListParameters{}, sdkVmMachineTypes.ListConfigs{})
	if err != nil {
		return nil, fmt.Errorf("could not load machine-type list, unexpected error: " + err.Error())
	}

	for _, x := range machineTypeList.InstanceTypes {
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
		return nil, fmt.Errorf("could not found machine-type ID with name: " + name)
	}
	return &machineType, nil
}

func (r *vmInstances) setValuesFromMachineType(data *vmInstancesResourceModel, server *vmInstancesMachineTypeModel) {
	data.MachineType.Disk = server.Disk
	data.MachineType.ID = server.ID
	data.MachineType.Name = server.Name
	data.MachineType.RAM = server.RAM
	data.MachineType.VCPUs = server.VCPUs
}

func (r *vmInstances) getVmStatus(id string) (*sdkVmInstances.GetResult, error) {
	getResult := &sdkVmInstances.GetResult{}
	expand := &sdkVmInstances.GetParametersExpand{"network"}

	duration := 5 * time.Minute
	startTime := time.Now()
	getParam := sdkVmInstances.GetParameters{Id: id, Expand: expand}
	getConfigParam := sdkVmInstances.GetConfigs{}
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
