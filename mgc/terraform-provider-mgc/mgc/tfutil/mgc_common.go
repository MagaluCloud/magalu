package tfutil

import (
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
