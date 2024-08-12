package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func convertStringSliceToTypesStringSlice(input []string) []types.String {
	result := make([]types.String, len(input))
	for i, v := range input {
		result[i] = types.StringValue(v)
	}
	return result
}
