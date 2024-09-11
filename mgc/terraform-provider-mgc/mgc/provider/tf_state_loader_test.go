package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	mgcSchemaPkg "magalu.cloud/core/schema"
)

var states = []tftypes.Value{
	tftypes.NewValue(tftypes.String, "test_string"),
	tftypes.NewValue(tftypes.Bool, true),
	tftypes.NewValue(tftypes.Number, 10),
	tftypes.NewValue(tftypes.Number, 10),
	tftypes.NewValue(tftypes.Number, 10.0),
	tftypes.NewValue(tftypes.Number, 0.000000000000000000000000001),
	tftypes.NewValue(
		tftypes.List{ElementType: tftypes.String},
		[]tftypes.Value{tftypes.NewValue(tftypes.String, "zero"), tftypes.NewValue(tftypes.String, "one")},
	),

	tftypes.NewValue(
		tftypes.List{
			ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"value": tftypes.String,
				},
			},
		},
		[]tftypes.Value{
			tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"value": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"value": tftypes.NewValue(tftypes.String, "myvalueresult"),
				},
			),
		},
	),

	tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"value": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"value_nested": tftypes.String,
					},
				},
			},
		},
		map[string]tftypes.Value{
			"value": tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"value_nested": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"value_nested": tftypes.NewValue(tftypes.String, "myvalueresult"),
				},
			),
		},
	),

	tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"value": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"tf_value_nested": tftypes.String,
					},
				},
			},
		},
		map[string]tftypes.Value{
			"value": tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"tf_value_nested": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"tf_value_nested": tftypes.NewValue(tftypes.String, "myvalueresult"),
				},
			),
		},
	),

	tftypes.NewValue(
		tftypes.List{
			ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"tf_value": tftypes.String,
				},
			},
		},
		[]tftypes.Value{
			tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"tf_value": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"tf_value": tftypes.NewValue(tftypes.String, "myvalueresult"),
				},
			),
		},
	),

	tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"allocate_fip":      tftypes.Bool,
				"availability_zone": tftypes.String,
				"created_at":        tftypes.String,
				"desired_image":     tftypes.String,
				"desired_status":    tftypes.String,
				"error":             tftypes.String,
				"id":                tftypes.String,
				"instance_id":       tftypes.String,
				"key_name":          tftypes.String,
				"memory":            tftypes.Number,
				"name":              tftypes.String,
				"power_state":       tftypes.Number,
				"power_state_label": tftypes.String,
				"root_storage":      tftypes.Number,
				"type":              tftypes.String,
				"updated_at":        tftypes.String,
				"user_data":         tftypes.String,
				"vcpus":             tftypes.Number,
			},
		},
		map[string]tftypes.Value{
			"allocate_fip":      tftypes.NewValue(tftypes.Bool, nil),
			"availability_zone": tftypes.NewValue(tftypes.String, "br-ne-1c"),
			"created_at":        tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"desired_image":     tftypes.NewValue(tftypes.String, "cloud-ubuntu-22.04 LTS"),
			"desired_status":    tftypes.NewValue(tftypes.String, "active"),
			"error":             tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"id":                tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"instance_id":       tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"key_name":          tftypes.NewValue(tftypes.String, "luizalabs-key"),
			"memory":            tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			"name":              tftypes.NewValue(tftypes.String, "my-tf-vm"),
			"power_state":       tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			"power_state_label": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"root_storage":      tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			"type":              tftypes.NewValue(tftypes.String, "cloud-bs1.xsmall"),
			"updated_at":        tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"user_data":         tftypes.NewValue(tftypes.String, nil),
			"vcpus":             tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
		},
	),
}

var schemas = []*mgcSchemaPkg.Schema{
	mgcSchemaPkg.NewStringSchema(),
	mgcSchemaPkg.NewBooleanSchema(),
	mgcSchemaPkg.NewIntegerSchema(),
	mgcSchemaPkg.NewNumberSchema(),
	mgcSchemaPkg.NewNumberSchema(),
	mgcSchemaPkg.NewNumberSchema(),
	mgcSchemaPkg.NewArraySchema(mgcSchemaPkg.NewStringSchema()),
	mgcSchemaPkg.NewArraySchema(
		mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"value": mgcSchemaPkg.NewStringSchema(),
		}, []string{"value"}),
	),
	mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
		"value": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"value_nested": mgcSchemaPkg.NewStringSchema(),
		}, []string{"value_nested"}),
	}, []string{"value"}),
	mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
		"value": mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"value_nested": mgcSchemaPkg.NewStringSchema(),
		}, []string{"value_nested"}),
	}, []string{"value"}),
	mgcSchemaPkg.NewArraySchema(
		mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"value": mgcSchemaPkg.NewStringSchema(),
		}, []string{"value"}),
	),
	mgcSchemaPkg.NewObjectSchema(
		map[string]*mgcSchemaPkg.Schema{
			"allocate_fip":      mgcSchemaPkg.NewBooleanSchema(),
			"availability_zone": mgcSchemaPkg.NewStringSchema(),
			"created_at":        mgcSchemaPkg.NewStringSchema(),
			"status":            mgcSchemaPkg.NewStringSchema(),
			"image":             mgcSchemaPkg.NewStringSchema(),
			"error":             mgcSchemaPkg.NewStringSchema(),
			"id":                mgcSchemaPkg.NewStringSchema(),
			"instance_id":       mgcSchemaPkg.NewStringSchema(),
			"key_name":          mgcSchemaPkg.NewStringSchema(),
			"memory":            mgcSchemaPkg.NewNumberSchema(),
			"name":              mgcSchemaPkg.NewStringSchema(),
			"power_state":       mgcSchemaPkg.NewNumberSchema(),
			"power_state_label": mgcSchemaPkg.NewStringSchema(),
			"root_storage":      mgcSchemaPkg.NewNumberSchema(),
			"type":              mgcSchemaPkg.NewStringSchema(),
			"updated_at":        mgcSchemaPkg.NewStringSchema(),
			"user_data":         mgcSchemaPkg.NewStringSchema(),
			"vcpus":             mgcSchemaPkg.NewNumberSchema(),
		}, []string{"name", "type", "key_name", "status", "image"},
	),
}

var results = []any{
	"test_string",
	true,
	(int64)(10),
	(float64)(10),
	(float64)(10.0),
	(float64)(0.000000000000000000000000001),
	[]any{"zero", "one"},
	[]any{map[string]any{"value": "myvalueresult"}},
	map[string]any{"value": map[string]any{"value_nested": "myvalueresult"}},
	map[string]any{"value": map[string]any{"value_nested": "myvalueresult"}},
	[]any{map[string]any{"value": "myvalueresult"}},
	map[string]any{
		"availability_zone": "br-ne-1c",
		"image":             "cloud-ubuntu-22.04 LTS",
		"status":            "active",
		"key_name":          "luizalabs-key",
		"name":              "my-tf-vm",
		"type":              "cloud-bs1.xsmall",
	},
}

var attrInfo = resAttrInfoMap{
	"value": {
		tfName: "value",
		mgcSchema: mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"value_nested": mgcSchemaPkg.NewStringSchema(),
		}, []string{}),
		childAttributes: resAttrInfoMap{
			"value_nested": {
				tfName:    "value_nested",
				mgcSchema: mgcSchemaPkg.NewStringSchema(),
			},
		},
	},
}
var attrInfoList = resAttrInfoMap{
	"0": {
		tfName: "0",
		mgcSchema: mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"value": mgcSchemaPkg.NewStringSchema(),
		}, []string{"value"}),
		childAttributes: resAttrInfoMap{
			"value": {
				tfName:    "value",
				mgcSchema: mgcSchemaPkg.NewStringSchema(),
			},
		},
	},
}
var attrInfoTFNameObjectNested = resAttrInfoMap{
	"value": {
		tfName: "value",
		mgcSchema: mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"value_nested": mgcSchemaPkg.NewStringSchema(),
		}, []string{}),
		childAttributes: resAttrInfoMap{
			"value_nested": {
				tfName:    "tf_value_nested",
				mgcSchema: mgcSchemaPkg.NewStringSchema(),
			},
		},
	},
}
var attrInfoTFNameObjectInList = resAttrInfoMap{
	"0": {
		tfName: "0",
		mgcSchema: mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			"value": mgcSchemaPkg.NewStringSchema(),
		}, []string{}),
		childAttributes: resAttrInfoMap{
			"value": {
				tfName:    "tf_value",
				mgcSchema: mgcSchemaPkg.NewStringSchema(),
			},
		},
	},
}

func tstCreateAttribute(mgcName mgcName, attrTfName tfName, mgcSchema *mgcSchemaPkg.Schema, isRequired bool, isOptional bool, isComputed bool, useStateForUnknown bool, requiresReplaceWhenChanged bool) *resAttrInfo {
	tfSchema, childAttrs, err := mgcSchemaToTFAttribute(
		mgcSchema,
		attributeModifiers{
			isRequired,
			isOptional,
			isComputed,
			useStateForUnknown,
			requiresReplaceWhenChanged,
			tfName(""),
			false,
			getInputChildModifiers,
		},
		context.Background(),
	)
	if err != nil {
		panic("Could not create test TF Schema")
	}
	return &resAttrInfo{
		mgcName:         mgcName,
		tfName:          attrTfName,
		mgcSchema:       mgcSchema,
		tfSchema:        tfSchema,
		childAttributes: childAttrs,
	}
}

var attrInfoTFInstanceCreate = resAttrInfoMap{
	"allocate_fip":      tstCreateAttribute("allocate_fip", "allocate_fip", mgcSchemaPkg.NewBooleanSchema(), false, true, false, false, false),
	"availability_zone": tstCreateAttribute("availability_zone", "availability_zone", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"created_at":        tstCreateAttribute("created_at", "created_at", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"image":             tstCreateAttribute("image", "desired_image", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"status":            tstCreateAttribute("status", "desired_status", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"error":             tstCreateAttribute("error", "error", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"id":                tstCreateAttribute("id", "id", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"instance_id":       tstCreateAttribute("instance_id", "instance_id", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"key_name":          tstCreateAttribute("key_name", "key_name", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"memory":            tstCreateAttribute("memory", "memory", mgcSchemaPkg.NewNumberSchema(), true, false, false, false, false),
	"name":              tstCreateAttribute("name", "name", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"power_state":       tstCreateAttribute("power_state", "power_state", mgcSchemaPkg.NewNumberSchema(), true, false, false, false, false),
	"power_state_label": tstCreateAttribute("power_state_label", "power_state_label", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"root_storage":      tstCreateAttribute("root_storage", "root_storage", mgcSchemaPkg.NewNumberSchema(), true, false, false, false, false),
	"type":              tstCreateAttribute("type", "type", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"updated_at":        tstCreateAttribute("updated_at", "updated_at", mgcSchemaPkg.NewStringSchema(), true, false, false, false, false),
	"user_data":         tstCreateAttribute("user_data", "user_data", mgcSchemaPkg.NewStringSchema(), false, true, false, false, false),
	"vcpus":             tstCreateAttribute("vcpus", "vcpus", mgcSchemaPkg.NewNumberSchema(), true, false, false, false, false),
}

var attrInfos = []resAttrInfoMap{
	{},
	{},
	{},
	{},
	{},
	{},
	{"0": {mgcSchema: mgcSchemaPkg.NewStringSchema()}},
	attrInfoList,
	attrInfo,
	attrInfoTFNameObjectNested,
	attrInfoTFNameObjectInList,
	attrInfoTFInstanceCreate,
}

func TestToMgcSchemaValue(t *testing.T) {
	for i := 0; i < len(states); i++ {
		atinfo := resAttrInfo{
			tfName:          "schema",
			mgcSchema:       schemas[i],
			childAttributes: attrInfos[i],
		}
		result, _, d := loadMgcSchemaValue(context.Background(), &atinfo, states[i], true, true)
		if !reflect.DeepEqual(result, results[i]) {
			t.Fatalf("result %d differs from expected: %T -> %T:\nRECEIVED: %+v\nEXPECTED: %+v\nDIAG: %+v\n", i, result, results[i], result, results[i], d)
		}
	}
}
