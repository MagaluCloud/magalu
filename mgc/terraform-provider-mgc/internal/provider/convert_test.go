package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"magalu.cloud/core"
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

var attrInfo = map[string]*attribute{
	"value": {
		name: "value",
		attributes: map[string]*attribute{
			"value_nested": {
				name: "value_nested",
			},
		},
	},
}
var attrInfoList = map[string]*attribute{
	"0": {
		name: "0",
		attributes: map[string]*attribute{
			"value": {
				name: "value",
			},
		},
	},
}
var attrInfoTFNameObjectNested = map[string]*attribute{
	"value": {
		name: "value",
		attributes: map[string]*attribute{
			"value_nested": {
				name: "tf_value_nested",
			},
		},
	},
}
var attrInfoTFNameObjectInList = map[string]*attribute{
	"0": {
		name: "0",
		attributes: map[string]*attribute{
			"value": {
				name: "tf_value",
			},
		},
	},
}
var attrInfoTFInstanceCreate = map[string]*attribute{
	"allocate_fip":      {name: "allocate_fip", isOptional: true, isComputed: false},
	"availability_zone": {name: "availability_zone"},
	"created_at":        {name: "created_at"},
	"image":             {name: "desired_image"},
	"status":            {name: "desired_status"},
	"error":             {name: "error"},
	"id":                {name: "id"},
	"instance_id":       {name: "instance_id"},
	"key_name":          {name: "key_name"},
	"memory":            {name: "memory"},
	"name":              {name: "name"},
	"power_state":       {name: "power_state"},
	"power_state_label": {name: "power_state_label"},
	"root_storage":      {name: "root_storage"},
	"type":              {name: "type"},
	"updated_at":        {name: "updated_at"},
	"user_data":         {name: "user_data", isOptional: true, isComputed: false},
	"vcpus":             {name: "vcpus"},
}

var attrInfos = []map[string]*attribute{
	{},
	{},
	{},
	{},
	{},
	{},
	{"0": {}},
	attrInfoList,
	attrInfo,
	attrInfoTFNameObjectNested,
	attrInfoTFNameObjectInList,
	attrInfoTFInstanceCreate,
}

func TestConvertTFToValue(t *testing.T) {
	conv := converter{
		ctx:  context.Background(),
		diag: diag.Diagnostics{},
	}

	for i := 0; i < len(states); i++ {
		atinfo := attribute{
			name:       "schema",
			attributes: attrInfos[i],
		}
		result := conv.convertTFToValue(schemas[i], &atinfo, states[i], true)
		if !reflect.DeepEqual(result, results[i]) {
			t.Fatalf("result differs from expected: %T:%+v %T:%+v %+v", result, result, results[i], results[i], conv.diag)
		}
	}
}
