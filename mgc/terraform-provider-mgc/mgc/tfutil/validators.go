package tfutil

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type CidrValidator struct{}

func (v CidrValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	cidr := req.ConfigValue.ValueString()

	if cidr == "0.0.0.0/0" || cidr == "::/0" {
		return
	}

	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR",
			fmt.Sprintf("Value %q is not a valid CIDR notation: %s", cidr, err),
		)
	}
}

func (v CidrValidator) Description(ctx context.Context) string {
	return "value must be a valid CIDR notation"
}

func (v CidrValidator) MarkdownDescription(ctx context.Context) string {
	return "value must be a valid CIDR notation"
}
