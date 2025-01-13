package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	sdkVmInstances "github.com/MagaluCloud/magalu/mgc/lib/products/virtual_machine/instances"
	sdkVmInstancesInterfaces "github.com/MagaluCloud/magalu/mgc/lib/products/virtual_machine/instances/network_interface"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	tfutil "github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
)

func NewVirtualMachineInterfaceAttachResource() resource.Resource {
	return &vmInterfaceAttach{}
}

type vmInterfaceAttach struct {
	sdkClient   *mgcSdk.Client
	vmInterface sdkVmInstancesInterfaces.Service
	vmInstance  sdkVmInstances.Service
}

type vmInterfaceAttachResourceModel struct {
	InterfaceID types.String `tfsdk:"interface_id"`
	InstanceID  types.String `tfsdk:"instance_id"`
}

func (r *vmInterfaceAttach) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_interface_attach"
}

func (r *vmInterfaceAttach) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.vmInterface = sdkVmInstancesInterfaces.NewService(ctx, r.sdkClient)
	r.vmInstance = sdkVmInstances.NewService(ctx, r.sdkClient)
}

func (r *vmInterfaceAttach) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attaches a network interface to a VM instance.",
		Attributes: map[string]schema.Attribute{
			"interface_id": schema.StringAttribute{
				Description: "ID of the network interface to attach.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Description: "ID of the VM instance to attach the interface to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *vmInterfaceAttach) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data vmInterfaceAttachResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := sdkVmInstancesInterfaces.AttachParameters{
		Instance: sdkVmInstancesInterfaces.AttachParametersInstance{
			Id: data.InstanceID.ValueString(),
		},
		Network: sdkVmInstancesInterfaces.AttachParametersNetwork{
			Interface: sdkVmInstancesInterfaces.AttachParametersNetworkInterface{
				Id: data.InterfaceID.ValueString(),
			},
		},
	}

	err := r.vmInterface.AttachContext(ctx, params, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstancesInterfaces.AttachConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Error attaching interface to VM", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vmInterfaceAttach) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data vmInterfaceAttachResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getInstance, err := r.vmInstance.GetContext(ctx, sdkVmInstances.GetParameters{
		Id:     data.InstanceID.ValueString(),
		Expand: &sdkVmInstances.GetParametersExpand{"network"},
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstances.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Error getting VM instance", err.Error())
		return
	}

	if getInstance.Network == nil || getInstance.Network.Interfaces == nil || len(*getInstance.Network.Interfaces) == 0 {
		resp.Diagnostics.AddError("No network interfaces found for VM instance", "")
		return
	}

	hasInterfaceId := slices.ContainsFunc(*getInstance.Network.Interfaces, func(i sdkVmInstances.GetResultNetworkInterfacesItem) bool {
		return i.Id == data.InterfaceID.ValueString()
	})

	if !hasInterfaceId {
		resp.Diagnostics.AddError("Network interface not found on VM instance", "")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vmInterfaceAttach) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Update is not supported for this resource")
}

func (r *vmInterfaceAttach) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vmInterfaceAttachResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.vmInterface.DetachContext(ctx,
		sdkVmInstancesInterfaces.DetachParameters{
			Instance: sdkVmInstancesInterfaces.DetachParametersInstance{
				Id: data.InstanceID.ValueString(),
			},
			Network: sdkVmInstancesInterfaces.DetachParametersNetwork{
				Interface: sdkVmInstancesInterfaces.DetachParametersNetworkInterface{
					Id: data.InterfaceID.ValueString(),
				},
			},
		},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVmInstancesInterfaces.DetachConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Error detaching interface from VM", err.Error())
	}
}
