package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type genericIDNameModel struct {
	Name types.String `tfsdk:"name"`
	ID   types.String `tfsdk:"id"`
}

type genericIDModel struct {
	ID types.String `tfsdk:"id"`
}
