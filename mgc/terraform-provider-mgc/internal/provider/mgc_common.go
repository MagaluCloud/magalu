package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type genericIDNameModel struct {
	genericIDModel
	Name types.String `tfsdk:"name"`
}

type genericIDModel struct {
	ID types.String `tfsdk:"id"`
}
