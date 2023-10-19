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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"magalu.cloud/core"
	mgc "magalu.cloud/core/schema"
	"magalu.cloud/sdk"
)

type testCase struct {
	res            *MgcResource
	expectedInput  mgcAttributes
	expectedOutput mgcAttributes
	expectedFinal  map[tfName]schema.Attribute
}

var create = core.NewRawStaticExecute(
	"mock create",
	"v1",
	"",
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{
			"image": mgc.NewStringSchema(),
			"name":  mgc.NewStringSchema(),
			"count": {
				Type:        "number",
				Description: "count description",
			},
		},
		[]string{"name", "image"},
	),
	nil,
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{"id": mgc.NewStringSchema()},
		[]string{"id"},
	),
	nil,
	nil,
	func(context context.Context, parameters core.Parameters, configs core.Configs) (result any, err error) {
		return nil, nil
	},
)

var read = core.NewRawStaticExecute(
	"mock read",
	"v1",
	"",
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{
			"id": mgc.NewStringSchema(),
		},
		[]string{"id"},
	),
	nil,
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{
			"id":        mgc.NewStringSchema(),
			"image":     mgc.NewStringSchema(),
			"name":      mgc.NewStringSchema(),
			"count":     mgc.NewIntegerSchema(),
			"createdAt": mgc.NewIntegerSchema(),
			"extra_field": mgc.NewArraySchema(
				mgc.NewObjectSchema(
					map[string]*mgc.Schema{
						"value": mgc.NewBooleanSchema(),
					},
					[]string{},
				),
			),
		},
		[]string{"id", "image", "name", "count", "createdAt", "extra_field"},
	),
	nil,
	nil,
	func(context context.Context, parameters core.Parameters, configs core.Configs) (result any, err error) {
		return nil, nil
	},
)

var update = core.NewRawStaticExecute(
	"mock update",
	"v1",
	"",
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{
			"id":    mgc.NewStringSchema(),
			"name":  mgc.NewStringSchema(),
			"count": mgc.NewNumberSchema(),
			"extra_field": mgc.NewArraySchema(
				mgc.NewObjectSchema(
					map[string]*mgc.Schema{
						"value": mgc.NewBooleanSchema(),
					},
					[]string{},
				),
			),
		},
		[]string{"id"},
	),
	nil,
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{
			"id": mgc.NewStringSchema(),
		},
		[]string{},
	),
	nil,
	nil,
	func(context context.Context, parameters core.Parameters, configs core.Configs) (result any, err error) {
		return nil, nil
	},
)

var delete = core.NewRawStaticExecute(
	"mock delete",
	"v1",
	"",
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{
			"id": mgc.NewStringSchema(),
		},
		[]string{"id"},
	),
	nil,
	mgc.NewNullSchema(),
	nil,
	nil,
	func(context context.Context, parameters core.Parameters, configs core.Configs) (result any, err error) {
		return nil, nil
	},
)

var testCases = []testCase{
	{
		res: &MgcResource{create: create, read: read, update: update, delete: delete},
		expectedInput: mgcAttributes{
			"count": {
				mgcName:   "count",
				tfName:    "count", // will be renamed to 'desired_count' for final
				mgcSchema: (*mgc.Schema)(create.ParametersSchema().Properties["count"].Value),
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
				mgcSchema: (*mgc.Schema)(update.ParametersSchema().Properties["extra_field"].Value),
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
				attributes: mgcAttributes{
					"0": {
						tfName:    "0",
						mgcName:   "0",
						mgcSchema: (*mgc.Schema)(update.ParametersSchema().Properties["extra_field"].Value.Items.Value),
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
						attributes: mgcAttributes{
							"value": {
								tfName:    "value",
								mgcName:   "value",
								mgcSchema: (*mgc.Schema)(update.ParametersSchema().Properties["extra_field"].Value.Items.Value.Properties["value"].Value),
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
				mgcSchema: (*mgc.Schema)(create.ParametersSchema().Properties["image"].Value),
				tfSchema: schema.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"name": {
				mgcName:   "name",
				tfName:    "name",
				mgcSchema: (*mgc.Schema)(create.ParametersSchema().Properties["name"].Value),
				tfSchema: schema.StringAttribute{
					Required:      true,
					PlanModifiers: []planmodifier.String{},
				},
			},
			"id": {
				mgcName:   "id",
				tfName:    "id",
				mgcSchema: (*mgc.Schema)(create.ParametersSchema().Properties["name"].Value),
				tfSchema: schema.StringAttribute{
					Computed:      true,
					PlanModifiers: []planmodifier.String{},
				},
			},
		},
		expectedOutput: mgcAttributes{
			"count": {
				mgcName:   "count",
				tfName:    "count", // will be renamed to 'current_count' for final
				mgcSchema: (*mgc.Schema)(read.ResultSchema().Properties["count"].Value),
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
				mgcSchema: (*mgc.Schema)(read.ResultSchema().Properties["createdAt"].Value),
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
				mgcSchema: (*mgc.Schema)(read.ResultSchema().Properties["extra_field"].Value),
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
				attributes: mgcAttributes{
					"0": {
						tfName:    "0",
						mgcName:   "0",
						mgcSchema: (*mgc.Schema)(read.ResultSchema().Properties["extra_field"].Value.Items.Value),
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
						attributes: mgcAttributes{
							"value": {
								tfName:    "value",
								mgcName:   "value",
								mgcSchema: (*mgc.Schema)(update.ParametersSchema().Properties["extra_field"].Value.Items.Value.Properties["value"].Value),
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
				mgcSchema: (*mgc.Schema)(create.ResultSchema().Properties["id"].Value),
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
				mgcSchema: (*mgc.Schema)(read.ResultSchema().Properties["image"].Value),
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
				mgcSchema: (*mgc.Schema)(read.ResultSchema().Properties["name"].Value),
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
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{},
			},
			"desired_count": schema.NumberAttribute{
				Optional:      true,
				Computed:      false,
				Description:   "count description",
				PlanModifiers: []planmodifier.Number{},
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
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{},
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
		testCase.res.ReadInputAttributes(ctx)
		if diff := deep.Equal(testCase.res.inputAttr, testCase.expectedInput); diff != nil {
			t.Errorf("MgcResource.readInputAttributes failed. Diff list: %v", diff)
		}
		testCase.res.ReadOutputAttributes(ctx)
		if diff := deep.Equal(testCase.res.outputAttr, testCase.expectedOutput); diff != nil {
			t.Errorf("MgcResource.readOutputAttributes failed. Diff list: %v", diff)
		}
		finalAttr, _ := generateTFAttributes(testCase.res, ctx)
		if diff := deep.Equal(finalAttr, testCase.expectedFinal); diff != nil {
			t.Errorf("MgcResource.generateTFAttributes failed. Diff list: %v", diff)
		}
	}
}

func TestMgcToTfSchemaDefaultValues(t *testing.T) {
	t.Run("non computed attriubte", func(t *testing.T) {
		s := mgc.NewStringSchema()
		s.Default = "default"
		m := attributeModifiers{}
		ctx := context.Background()

		var expected defaults.String = nil

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)

		found := sAtt.(schema.StringAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("no default", func(t *testing.T) {
		s := mgc.NewStringSchema()
		m := attributeModifiers{}
		ctx := context.Background()

		var expected defaults.String = nil

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)

		found := sAtt.(schema.StringAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("string", func(t *testing.T) {
		def := "foo"

		s := mgc.NewStringSchema()
		s.Default = def
		m := attributeModifiers{isComputed: true}
		ctx := context.Background()

		expected := stringdefault.StaticString(def)

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)

		found := sAtt.(schema.StringAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("number", func(t *testing.T) {
		def := float64(3.14)

		s := mgc.NewNumberSchema()
		s.Default = def
		m := attributeModifiers{isComputed: true}
		ctx := context.Background()

		expected := numberdefault.StaticBigFloat(big.NewFloat(def))

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)

		found := sAtt.(schema.NumberAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("integer", func(t *testing.T) {
		def := int64(0)

		s := mgc.NewIntegerSchema()
		s.Default = def
		m := attributeModifiers{isComputed: true}
		ctx := context.Background()

		expected := int64default.StaticInt64(def)

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)

		found := sAtt.(schema.Int64Attribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("boolean", func(t *testing.T) {
		def := false

		s := mgc.NewBooleanSchema()
		s.Default = def
		m := attributeModifiers{isComputed: true}
		ctx := context.Background()

		expected := booldefault.StaticBool(def)

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)

		found := sAtt.(schema.BoolAttribute).Default

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected default == %+v, found: %+v", expected, found)
		}
	})

	t.Run("object", func(t *testing.T) {
		nameSchema := mgc.NewStringSchema()
		nameSchema.Default = "pedro"

		p := map[string]*mgc.Schema{
			"name": nameSchema,
			"age":  mgc.NewIntegerSchema(),
		}

		s := mgc.NewObjectSchema(p, []string{})
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

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)
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
		s := mgc.NewArraySchema(mgc.NewStringSchema())
		s.Default = []any{"hello", "world"}

		ctx := context.Background()

		m := attributeModifiers{
			isComputed: true,
			getChildModifiers: func(ctx context.Context, mgcSchema *sdk.Schema, mgcName mgcName) attributeModifiers {
				return attributeModifiers{}
			},
		}

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)
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
		s := mgc.NewArraySchema(mgc.NewStringSchema())
		s.Default = []any{}

		ctx := context.Background()

		m := attributeModifiers{
			isComputed: true,
			getChildModifiers: func(ctx context.Context, mgcSchema *sdk.Schema, mgcName mgcName) attributeModifiers {
				return attributeModifiers{isComputed: true}
			},
		}

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)
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
		s := mgc.NewArraySchema(
			mgc.NewObjectSchema(
				map[string]*mgc.Schema{
					"key": mgc.NewStringSchema(),
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

		sAtt, _, _ := mgcToTFSchema(s, m, ctx)
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
