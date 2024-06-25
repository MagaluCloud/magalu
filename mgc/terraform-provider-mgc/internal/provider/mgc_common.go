package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type genericIDNameModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
