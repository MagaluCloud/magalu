package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &objectStorage{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewObjectStorageResource() resource.Resource {
	return &objectStorage{}
}

// orderResource is the resource implementation.
type objectStorage struct{}

// Metadata returns the resource type name.
func (r *objectStorage) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage"
}

// Schema defines the schema for the resource.
func (r *objectStorage) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Create creates the resource and sets the initial Terraform state.
func (r *objectStorage) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read refreshes the Terraform state with the latest data.
func (r *objectStorage) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *objectStorage) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *objectStorage) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
