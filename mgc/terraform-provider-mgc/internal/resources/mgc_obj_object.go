package resources

import "github.com/hashicorp/terraform-plugin-framework/types"

type ObjectStorageObject struct {
	Bucket       types.String `tfsdk:"bucket"`
	Source       types.String `tfsdk:"source"`
	StorageClass types.String `tfsdk:"storage_class"`
}
