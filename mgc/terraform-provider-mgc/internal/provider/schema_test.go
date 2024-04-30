package provider

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/sdk"
)

type testCase struct {
	res            *MgcResource
	expectedInput  resAttrInfoMap
	expectedOutput resAttrInfoMap
	expectedFinal  map[tfName]schema.Attribute
}

var create = core.NewSimpleExecutor(
	core.ExecutorSpec{
		DescriptorSpec: core.DescriptorSpec{
			Name:        "mock create",
			Version:     "v1",
			Description: "mock create",
		},
		ParametersSchema: mgcSchemaPkg.NewObjectSchema(
			map[string]*mgcSchemaPkg.Schema{
				"image": mgcSchemaPkg.NewStringSchema(),
				"name":  mgcSchemaPkg.NewStringSchema(),
				"count": {
					Type:        "number",
					Description: "count description",
				},
			},
			[]string{"name", "image"},
		),
		ConfigsSchema: mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{}, []string{}),
		ResultSchema: mgcSchemaPkg.NewObjectSchema(
			map[string]*mgcSchemaPkg.Schema{"id": mgcSchemaPkg.NewStringSchema()},
			[]string{"id"},
		),
		Execute: func(e core.Executor, context context.Context, parameters core.Parameters, configs core.Configs) (result core.Result, err error) {
			return nil, nil
		},
	},
)

var read = core.NewSimpleExecutor(
	core.ExecutorSpec{
		DescriptorSpec: core.DescriptorSpec{
			Name:        "mock read",
			Version:     "v1",
			Description: "mock read",
		},
		ParametersSchema: mgcSchemaPkg.NewObjectSchema(
			map[string]*mgcSchemaPkg.Schema{
				"id": mgcSchemaPkg.NewStringSchema(),
			},
			[]string{"id"},
		),
		ConfigsSchema: mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{}, []string{}),
		ResultSchema: mgcSchemaPkg.NewObjectSchema(
			map[string]*mgcSchemaPkg.Schema{
				"id":        mgcSchemaPkg.NewStringSchema(),
				"image":     mgcSchemaPkg.NewStringSchema(),
				"name":      mgcSchemaPkg.NewStringSchema(),
				"count":     mgcSchemaPkg.NewIntegerSchema(),
				"createdAt": mgcSchemaPkg.NewIntegerSchema(),
				"extra_field": mgcSchemaPkg.NewArraySchema(
					mgcSchemaPkg.NewObjectSchema(
						map[string]*mgcSchemaPkg.Schema{
							"value": mgcSchemaPkg.NewBooleanSchema(),
						},
						[]string{},
					),
				),
			},
			[]string{"id", "image", "name", "count", "createdAt", "extra_field"},
		),
		Execute: func(e core.Executor, context context.Context, parameters core.Parameters, configs core.Configs) (result core.Result, err error) {
			return nil, nil
		},
	},
)

var update = core.NewSimpleExecutor(
	core.ExecutorSpec{
		DescriptorSpec: core.DescriptorSpec{
			Name:        "mock update",
			Version:     "v1",
			Description: "mock update",
		},
		ParametersSchema: mgcSchemaPkg.NewObjectSchema(
			map[string]*mgcSchemaPkg.Schema{
				"id":    mgcSchemaPkg.NewStringSchema(),
				"name":  mgcSchemaPkg.NewStringSchema(),
				"count": mgcSchemaPkg.NewNumberSchema(),
				"extra_field": mgcSchemaPkg.NewArraySchema(
					mgcSchemaPkg.NewObjectSchema(
						map[string]*mgcSchemaPkg.Schema{
							"value": mgcSchemaPkg.NewBooleanSchema(),
						},
						[]string{},
					),
				),
			},
			[]string{"id"},
		),
		ConfigsSchema: mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{}, []string{}),
		ResultSchema: mgcSchemaPkg.NewObjectSchema(
			map[string]*mgcSchemaPkg.Schema{
				"id": mgcSchemaPkg.NewStringSchema(),
			},
			[]string{},
		),
		Execute: func(e core.Executor, context context.Context, parameters core.Parameters, configs core.Configs) (result core.Result, err error) {
			return nil, nil
		},
	},
)

var deleteOp = core.NewSimpleExecutor(
	core.ExecutorSpec{
		DescriptorSpec: core.DescriptorSpec{
			Name:        "mock delete",
			Version:     "v1",
			Description: "mock delete",
		},
		ParametersSchema: mgcSchemaPkg.NewObjectSchema(
			map[string]*mgcSchemaPkg.Schema{
				"id": mgcSchemaPkg.NewStringSchema(),
			},
			[]string{"id"},
		),
		ConfigsSchema: mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{}, []string{}),
		ResultSchema:  mgcSchemaPkg.NewNullSchema(),
		Execute: func(e core.Executor, context context.Context, parameters core.Parameters, configs core.Configs) (result core.Result, err error) {
			return nil, nil
		},
	},
)

var testCases = []testCase{
	{
		res: &MgcResource{create: create, read: read, update: update, delete: deleteOp},
		expectedInput: resAttrInfoMap{
			"count": {
				mgcName:   "count",
				tfName:    "count",
				mgcSchema: (*mgcSchemaPkg.Schema)(create.ParametersSchema().Properties["count"].Value),
				tfSchema: schema.NumberAttribute{
					Description:   "count description",
					Optional:      true,
					Computed:      false, // False because read result attr has different schema
					PlanModifiers: []planmodifier.Number{},
				},
			},
			"extra_field": {
				tfName:    "extra_field",
				mgcName:   "extra_field",
				mgcSchema: (*mgcSchemaPkg.Schema)(update.ParametersSchema().Properties["extra_field"].Value),
				tfSchema: schema.ListNestedAttribute{
					NestedObject: schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"value": schema.BoolAttribute{
								Optional: true,
								Computed: false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional: true,
						Computed: false,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					}.GetNestedObject().(schema.NestedAttributeObject),
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
				},
				childAttributes: resAttrInfoMap{
					"0": {
						tfName:    "0",
						mgcName:   "0",
						mgcSchema: (*mgcSchemaPkg.Schema)(update.ParametersSchema().Properties["extra_field"].Value.Items.Value),
						tfSchema: schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"value": schema.BoolAttribute{
									Optional: true,
									Computed: false,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							Optional: true,
							Computed: false,
							PlanModifiers: []planmodifier.Object{
								objectplanmodifier.UseStateForUnknown(),
							},
						},
						childAttributes: resAttrInfoMap{
							"value": {
								tfName:    "value",
								mgcName:   "value",
								mgcSchema: (*mgcSchemaPkg.Schema)(update.ParametersSchema().Properties["extra_field"].Value.Items.Value.Properties["value"].Value),
								tfSchema: schema.BoolAttribute{
									Optional: true,
									Computed: false,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"image": {
				mgcName:   "image",
				tfName:    "image",
				mgcSchema: (*mgcSchemaPkg.Schema)(create.ParametersSchema().Properties["image"].Value),
				tfSchema: schema.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"name": {
				mgcName:   "name",
				tfName:    "name",
				mgcSchema: (*mgcSchemaPkg.Schema)(create.ParametersSchema().Properties["name"].Value),
				tfSchema: schema.StringAttribute{
					Required:      true,
					PlanModifiers: []planmodifier.String{},
				},
			},
			"id": {
				mgcName:   "id",
				tfName:    "id",
				mgcSchema: (*mgcSchemaPkg.Schema)(create.ParametersSchema().Properties["name"].Value),
				tfSchema: schema.StringAttribute{
					Computed:      true,
					PlanModifiers: []planmodifier.String{},
				},
			},
		},
		expectedOutput: resAttrInfoMap{
			"count": {
				mgcName:   "count",
				tfName:    "count", // will be renamed to 'current_count' for final
				mgcSchema: (*mgcSchemaPkg.Schema)(read.ResultSchema().Properties["count"].Value),
				tfSchema: schema.Int64Attribute{
					Computed: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
				},
			},
			"createdAt": {
				mgcName:   "createdAt",
				tfName:    "created_at",
				mgcSchema: (*mgcSchemaPkg.Schema)(read.ResultSchema().Properties["createdAt"].Value),
				tfSchema: schema.Int64Attribute{
					Computed: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
				},
			},
			"extra_field": {
				mgcName:   "extra_field",
				tfName:    "extra_field",
				mgcSchema: (*mgcSchemaPkg.Schema)(read.ResultSchema().Properties["extra_field"].Value),
				tfSchema: schema.ListNestedAttribute{
					NestedObject: schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"value": schema.BoolAttribute{
								Computed: false,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Computed: false,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					}.GetNestedObject().(schema.NestedAttributeObject),
					Computed: false,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
				},
				childAttributes: resAttrInfoMap{
					"0": {
						tfName:    "0",
						mgcName:   "0",
						mgcSchema: (*mgcSchemaPkg.Schema)(read.ResultSchema().Properties["extra_field"].Value.Items.Value),
						tfSchema: schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"value": schema.BoolAttribute{
									Computed: true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
							Computed: true,
							PlanModifiers: []planmodifier.Object{
								objectplanmodifier.UseStateForUnknown(),
							},
						},
						childAttributes: resAttrInfoMap{
							"value": {
								tfName:    "value",
								mgcName:   "value",
								mgcSchema: (*mgcSchemaPkg.Schema)(update.ParametersSchema().Properties["extra_field"].Value.Items.Value.Properties["value"].Value),
								tfSchema: schema.BoolAttribute{
									Computed: true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"id": {
				mgcName:   "id",
				tfName:    "id",
				mgcSchema: (*mgcSchemaPkg.Schema)(create.ResultSchema().Properties["id"].Value),
				tfSchema: schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"image": {
				mgcName:   "image",
				tfName:    "image",
				mgcSchema: (*mgcSchemaPkg.Schema)(read.ResultSchema().Properties["image"].Value),
				tfSchema: schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": {
				mgcName:   "name",
				tfName:    "name",
				mgcSchema: (*mgcSchemaPkg.Schema)(read.ResultSchema().Properties["name"].Value),
				tfSchema: schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		expectedFinal: map[tfName]schema.Attribute{
			"current_count": schema.Int64Attribute{
				Required:      false,
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{},
			},
			"count": schema.NumberAttribute{
				Optional:    true,
				Computed:    false,
				Description: "count description",
				PlanModifiers: []planmodifier.Number{
					numberplanmodifier.UseStateForUnknown(),
				},
			},
			"extra_field": schema.ListNestedAttribute{
				NestedObject: schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"value": schema.BoolAttribute{
							Optional: true,
							Computed: false,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
					Optional: true,
					Computed: false,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
				}.GetNestedObject().(schema.NestedAttributeObject),
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	},
}

func TestGenerateTFAttributes(t *testing.T) {
	ctx := context.Background()
	// Use deep.Equal because we might compare two functions, and reflect.DeepEqual always fails on function
	// comparisons (unless both functions are nil)
	// TODO: investigate lib reporting false negatives/positives on deep comparisons...
	for _, testCase := range testCases {
		res, err := newMgcResource(
			ctx,
			testCase.res.sdk,
			"resource", "resource",
			"",
			testCase.res.create, testCase.res.read, testCase.res.update, testCase.res.delete,
			nil,
		)
		if err != nil {
			t.Errorf("unable to instantiate resource for test case: %v", err)
		}
		if diff := deep.Equal(res.attrTree.input, testCase.expectedInput); diff != nil {
			t.Errorf("MgcResource.readInputAttributes failed. Diff list: %v", diff)
		}
		if diff := deep.Equal(res.attrTree.output, testCase.expectedOutput); diff != nil {
			t.Errorf("MgcResource.readOutputAttributes failed. Diff list: %v", diff)
		}
		finalAttr := generateTFAttributes(ctx, res.attrTree)
		if diff := deep.Equal(finalAttr, testCase.expectedFinal); diff != nil {
			t.Errorf("MgcResource.generateTFAttributes failed. Diff list: %v", diff)
		}
	}
}

func TestMgcToTfSchemaDefaultValues(t *testing.T) {
	t.Run("required attributes", func(t *testing.T) {
		s := mgcSchemaPkg.NewStringSchema()
		s.Default = "default"
		m := attributeModifiers{isRequired: true}
		ctx := context.Background()

		var expected defaults.String = nil

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)

		found := sAtt.(schema.StringAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("no default", func(t *testing.T) {
		s := mgcSchemaPkg.NewStringSchema()
		m := attributeModifiers{}
		ctx := context.Background()

		var expected defaults.String = nil

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)

		found := sAtt.(schema.StringAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("string", func(t *testing.T) {
		def := "foo"

		s := mgcSchemaPkg.NewStringSchema()
		s.Default = def
		m := attributeModifiers{isComputed: true}
		ctx := context.Background()

		expected := stringdefault.StaticString(def)

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)

		found := sAtt.(schema.StringAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("number", func(t *testing.T) {
		def := float64(3.14)

		s := mgcSchemaPkg.NewNumberSchema()
		s.Default = def
		m := attributeModifiers{isComputed: true}
		ctx := context.Background()

		expected := numberdefault.StaticBigFloat(big.NewFloat(def))

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)

		found := sAtt.(schema.NumberAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("integer", func(t *testing.T) {
		def := int64(0)

		s := mgcSchemaPkg.NewIntegerSchema()
		s.Default = def
		m := attributeModifiers{isComputed: true}
		ctx := context.Background()

		expected := int64default.StaticInt64(def)

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)

		found := sAtt.(schema.Int64Attribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("boolean", func(t *testing.T) {
		def := false

		s := mgcSchemaPkg.NewBooleanSchema()
		s.Default = def
		m := attributeModifiers{isComputed: true}
		ctx := context.Background()

		expected := booldefault.StaticBool(def)

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)

		found := sAtt.(schema.BoolAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("object", func(t *testing.T) {
		nameSchema := mgcSchemaPkg.NewStringSchema()
		nameSchema.Default = "pedro"

		p := map[string]*mgcSchemaPkg.Schema{
			"name": nameSchema,
			"age":  mgcSchemaPkg.NewIntegerSchema(),
		}

		s := mgcSchemaPkg.NewObjectSchema(p, []string{})
		s.Default = map[string]any{
			"name": "pedro",
			"age":  int64(10),
		}

		m := attributeModifiers{
			isComputed: true,
			getChildModifiers: func(ctx context.Context, mgcSchema *sdk.Schema, mgcName mgcName) attributeModifiers {
				return attributeModifiers{isComputed: true}
			},
		}

		ctx := context.Background()

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)
		found := sAtt.(schema.SingleNestedAttribute).Default

		obj, _ := types.ObjectValue(
			map[string]attr.Type{
				"name": types.StringType,
				"age":  types.Int64Type,
			},
			map[string]attr.Value{
				"name": types.StringValue("pedro"),
				"age":  types.Int64Value(10),
			},
		)
		expected := objectdefault.StaticValue(obj)

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("list", func(t *testing.T) {
		s := mgcSchemaPkg.NewArraySchema(mgcSchemaPkg.NewStringSchema())
		s.Default = []any{"hello", "world"}

		ctx := context.Background()

		m := attributeModifiers{
			isComputed: true,
			getChildModifiers: func(ctx context.Context, mgcSchema *sdk.Schema, mgcName mgcName) attributeModifiers {
				return attributeModifiers{}
			},
		}

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)
		found := sAtt.(schema.ListAttribute).Default

		lst, _ := types.ListValue(
			types.StringType,
			[]attr.Value{types.StringValue("hello"), types.StringValue("world")},
		)

		expected := listdefault.StaticValue(lst)

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("list empty", func(t *testing.T) {
		s := mgcSchemaPkg.NewArraySchema(mgcSchemaPkg.NewStringSchema())
		s.Default = []any{}

		ctx := context.Background()

		m := attributeModifiers{
			isComputed: true,
			getChildModifiers: func(ctx context.Context, mgcSchema *sdk.Schema, mgcName mgcName) attributeModifiers {
				return attributeModifiers{isComputed: true}
			},
		}

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)
		found := sAtt.(schema.ListAttribute).Default

		lst, _ := types.ListValue(
			types.StringType,
			[]attr.Value{},
		)

		expected := listdefault.StaticValue(lst)

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("list nested", func(t *testing.T) {
		s := mgcSchemaPkg.NewArraySchema(
			mgcSchemaPkg.NewObjectSchema(
				map[string]*mgcSchemaPkg.Schema{
					"key": mgcSchemaPkg.NewStringSchema(),
				},
				[]string{},
			),
		)

		s.Default = []any{
			map[string]any{
				"key": "hello",
			},
			map[string]any{
				"key": "world",
			},
		}

		m := attributeModifiers{
			isComputed: true,
			getChildModifiers: func(ctx context.Context, mgcSchema *sdk.Schema, n mgcName) attributeModifiers {
				return attributeModifiers{
					isComputed: true,
					getChildModifiers: func(ctx context.Context, mgcSchema *sdk.Schema, m mgcName) attributeModifiers {
						return attributeModifiers{isComputed: true}
					},
				}
			},
		}

		ctx := context.Background()

		sAtt, _, _ := mgcSchemaToTFAttribute(s, m, ctx)
		found := sAtt.(schema.ListNestedAttribute).Default

		hello, _ := types.ObjectValue(
			map[string]attr.Type{"key": types.StringType},
			map[string]attr.Value{"key": types.StringValue("hello")},
		)
		world, _ := types.ObjectValue(
			map[string]attr.Type{"key": types.StringType},
			map[string]attr.Value{"key": types.StringValue("world")},
		)
		lst, _ := types.ListValue(
			types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType}},
			[]attr.Value{hello, world},
		)

		expected := listdefault.StaticValue(lst)

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})
}
