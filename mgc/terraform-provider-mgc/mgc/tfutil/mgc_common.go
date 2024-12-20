package tfutil

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GenericIDNameModel struct {
	Name types.String `tfsdk:"name"`
	ID   types.String `tfsdk:"id"`
}

type GenericIDModel struct {
	ID types.String `tfsdk:"id"`
}

func ConvertInt64PointerToIntPointer(int64Ptr *int64) *int {
	if int64Ptr == nil {
		return nil
	}
	intVal := int(*int64Ptr)
	return &intVal
}

func ConvertIntPointerToInt64Pointer(intPtr *int) *int64 {
	if intPtr == nil {
		return nil
	}
	int64Val := int64(*intPtr)
	return &int64Val
}

type ResponseFrom interface {
	resource.ConfigureResponse | datasource.ConfigureResponse
}

func AddCLIAuthWarning[T ResponseFrom](resp *T) {
	title := "[DEPRECATED] Using CLI Auth is not recommended and will be removed in future versions."
	text := "Please note that authentication via the Command Line Interface (CLI) will be discontinued.\nGoing forward, you will need to use API Key authentication exclusively.\nAccess the documentation https://docs.magalu.cloud/docs/devops-tools/terraform/how-to/auth#autentica%C3%A7%C3%A3o-com-api-key"

	switch tp := any(resp).(type) {
	case *resource.ConfigureResponse:
		tp.Diagnostics.AddWarning(title, text)
	case *datasource.ConfigureResponse:
		tp.Diagnostics.AddWarning(title, text)
	}
}

func SdkParamValueToString(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case *string:
		if val != nil {
			return *val
		}
	case *float64:
		if val != nil {
			return fmt.Sprintf("%g", *val)
		}
	case *int:
		if val != nil {
			return fmt.Sprint(*val)
		}
	case *bool:
		if val != nil {
			return fmt.Sprint(*val)
		}
	case string:
		return val
	case float64:
		return fmt.Sprintf("%g", val)
	case int:
		return fmt.Sprint(val)
	case bool:
		return fmt.Sprint(val)
	}

	return ""
}
