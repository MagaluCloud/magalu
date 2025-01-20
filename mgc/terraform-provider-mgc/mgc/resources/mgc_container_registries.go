package resources

import (
	"context"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	registries "github.com/MagaluCloud/magalu/mgc/lib/products/container_registry/registries"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ContainerRegistryModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type ContainerRegistryResource struct {
	sdkClient       *mgcSdk.Client
	registryService registries.Service
}

func NewContainerRegistryRegistriesResource() resource.Resource {
	return &ContainerRegistryResource{}
}

func (r *ContainerRegistryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_registries"
}

func (r *ContainerRegistryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.registryService = registries.NewService(ctx, r.sdkClient)
}

func (r *ContainerRegistryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Container Registry",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the registry",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the registry",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ContainerRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ContainerRegistryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.registryService.CreateContext(ctx, registries.CreateParameters{
		Name: data.Name.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, registries.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create registry", err.Error())
		return
	}

	data.Id = types.StringValue(created.Id)
	data.Name = types.StringValue(created.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContainerRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ContainerRegistryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	registry, err := r.registryService.GetContext(ctx, registries.GetParameters{
		RegistryId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, registries.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read registry", err.Error())
		return
	}

	data.Name = types.StringValue(registry.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContainerRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Container Registry cannot be updated after creation",
	)
}

func (r *ContainerRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ContainerRegistryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.registryService.DeleteContext(ctx, registries.DeleteParameters{
		RegistryId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, registries.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete registry", err.Error())
		return
	}
}

func (r *ContainerRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data := &ContainerRegistryModel{Id: types.StringValue(req.ID)}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
