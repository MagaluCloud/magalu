package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	core "magalu.cloud/core/schema"
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

var schemas = []*core.Schema{
	core.NewStringSchema(),
	core.NewBooleanSchema(),
	core.NewIntegerSchema(),
	core.NewNumberSchema(),
	core.NewNumberSchema(),
	core.NewNumberSchema(),
	core.NewArraySchema(core.NewStringSchema()),
	core.NewArraySchema(
		core.NewObjectSchema(map[string]*core.Schema{
			"value": core.NewStringSchema(),
		}, []string{"value"}),
	),
	core.NewObjectSchema(map[string]*core.Schema{
		"value": core.NewObjectSchema(map[string]*core.Schema{
			"value_nested": core.NewStringSchema(),
		}, []string{"value_nested"}),
	}, []string{"value"}),
	core.NewObjectSchema(map[string]*core.Schema{
		"value": core.NewObjectSchema(map[string]*core.Schema{
			"value_nested": core.NewStringSchema(),
		}, []string{"value_nested"}),
	}, []string{"value"}),
	core.NewArraySchema(
		core.NewObjectSchema(map[string]*core.Schema{
			"value": core.NewStringSchema(),
		}, []string{"value"}),
	),
	core.NewObjectSchema(
		map[string]*core.Schema{
			"allocate_fip":      core.NewBooleanSchema(),
			"availability_zone": core.NewStringSchema(),
			"created_at":        core.NewStringSchema(),
			"status":            core.NewStringSchema(),
			"image":             core.NewStringSchema(),
			"error":             core.NewStringSchema(),
			"id":                core.NewStringSchema(),
			"instance_id":       core.NewStringSchema(),
			"key_name":          core.NewStringSchema(),
			"memory":            core.NewNumberSchema(),
			"name":              core.NewStringSchema(),
			"power_state":       core.NewNumberSchema(),
			"power_state_label": core.NewStringSchema(),
			"root_storage":      core.NewNumberSchema(),
			"type":              core.NewStringSchema(),
			"updated_at":        core.NewStringSchema(),
			"user_data":         core.NewStringSchema(),
			"vcpus":             core.NewNumberSchema(),
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

var attrInfo = mgcAttributes{
	"value": {
		tfName: "value",
		mgcSchema: core.NewObjectSchema(map[string]*core.Schema{
			"value_nested": core.NewStringSchema(),
		}, []string{}),
		attributes: mgcAttributes{
			"value_nested": {
				tfName:    "value_nested",
				mgcSchema: core.NewStringSchema(),
			},
		},
	},
}
var attrInfoList = mgcAttributes{
	"0": {
		tfName: "0",
		mgcSchema: core.NewObjectSchema(map[string]*core.Schema{
			"value": core.NewStringSchema(),
		}, []string{"value"}),
		attributes: mgcAttributes{
			"value": {
				tfName:    "value",
				mgcSchema: core.NewStringSchema(),
			},
		},
	},
}
var attrInfoTFNameObjectNested = mgcAttributes{
	"value": {
		tfName: "value",
		mgcSchema: core.NewObjectSchema(map[string]*core.Schema{
			"value_nested": core.NewStringSchema(),
		}, []string{}),
		attributes: mgcAttributes{
			"value_nested": {
				tfName:    "tf_value_nested",
				mgcSchema: core.NewStringSchema(),
			},
		},
	},
}
var attrInfoTFNameObjectInList = mgcAttributes{
	"0": {
		tfName: "0",
		mgcSchema: core.NewObjectSchema(map[string]*core.Schema{
			"value": core.NewStringSchema(),
		}, []string{}),
		attributes: mgcAttributes{
			"value": {
				tfName:    "tf_value",
				mgcSchema: core.NewStringSchema(),
			},
		},
	},
}

func tstCreateAttribute(mgcName mgcName, tfName tfName, mgcSchema *core.Schema, isRequired bool, isOptional bool, isComputed bool, useStateForUnknown bool, requiresReplaceWhenChanged bool) *attribute {
	tfSchema, childAttrs, err := mgcToTFSchema(mgcSchema, attributeModifiers{isRequired, isOptional, isComputed, useStateForUnknown, requiresReplaceWhenChanged, getInputChildModifiers}, context.Background())
	if err != nil {
		panic("Could not create test TF Schema")
	}
	return &attribute{
		mgcName:    mgcName,
		tfName:     tfName,
		mgcSchema:  mgcSchema,
		tfSchema:   tfSchema,
		attributes: childAttrs,
	}
}

var attrInfoTFInstanceCreate = mgcAttributes{
	"allocate_fip":      tstCreateAttribute("allocate_fip", "allocate_fip", core.NewBooleanSchema(), false, true, false, false, false),
	"availability_zone": tstCreateAttribute("availability_zone", "availability_zone", core.NewStringSchema(), true, false, false, false, false),
	"created_at":        tstCreateAttribute("created_at", "created_at", core.NewStringSchema(), true, false, false, false, false),
	"image":             tstCreateAttribute("image", "desired_image", core.NewStringSchema(), true, false, false, false, false),
	"status":            tstCreateAttribute("status", "desired_status", core.NewStringSchema(), true, false, false, false, false),
	"error":             tstCreateAttribute("error", "error", core.NewStringSchema(), true, false, false, false, false),
	"id":                tstCreateAttribute("id", "id", core.NewStringSchema(), true, false, false, false, false),
	"instance_id":       tstCreateAttribute("instance_id", "instance_id", core.NewStringSchema(), true, false, false, false, false),
	"key_name":          tstCreateAttribute("key_name", "key_name", core.NewStringSchema(), true, false, false, false, false),
	"memory":            tstCreateAttribute("memory", "memory", core.NewNumberSchema(), true, false, false, false, false),
	"name":              tstCreateAttribute("name", "name", core.NewStringSchema(), true, false, false, false, false),
	"power_state":       tstCreateAttribute("power_state", "power_state", core.NewNumberSchema(), true, false, false, false, false),
	"power_state_label": tstCreateAttribute("power_state_label", "power_state_label", core.NewStringSchema(), true, false, false, false, false),
	"root_storage":      tstCreateAttribute("root_storage", "root_storage", core.NewNumberSchema(), true, false, false, false, false),
	"type":              tstCreateAttribute("type", "type", core.NewStringSchema(), true, false, false, false, false),
	"updated_at":        tstCreateAttribute("updated_at", "updated_at", core.NewStringSchema(), true, false, false, false, false),
	"user_data":         tstCreateAttribute("user_data", "user_data", core.NewStringSchema(), false, true, false, false, false),
	"vcpus":             tstCreateAttribute("vcpus", "vcpus", core.NewNumberSchema(), true, false, false, false, false),
}

var attrInfos = []mgcAttributes{
	{},
	{},
	{},
	{},
	{},
	{},
	{"0": {mgcSchema: core.NewStringSchema()}},
	attrInfoList,
	attrInfo,
	attrInfoTFNameObjectNested,
	attrInfoTFNameObjectInList,
	attrInfoTFInstanceCreate,
}

func TestToMgcSchemaValue(t *testing.T) {
	conv := tfStateConverter{
		ctx:  context.Background(),
		diag: &diag.Diagnostics{},
	}

	for i := 0; i < len(states); i++ {
		atinfo := attribute{
			tfName:     "schema",
			mgcSchema:  schemas[i],
			attributes: attrInfos[i],
		}
		result, _ := conv.toMgcSchemaValue(&atinfo, states[i], true, true)
		if !reflect.DeepEqual(result, results[i]) {
			t.Fatalf("result %d differs from expected: %T -> %T:\nRECEIVED: %+v\nEXPECTED: %+v\nDIAG: %+v\n", i, result, results[i], result, results[i], conv.diag)
		}
	}
}
