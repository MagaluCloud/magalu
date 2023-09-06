package provider

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	mgc "magalu.cloud/core"
)

type testCase struct {
	res            *MgcResource
	expectedInput  mgcAttributes
	expectedOutput mgcAttributes
	expectedFinal  map[tfName]schema.Attribute
}

var create = mgc.NewRawStaticExecute(
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
	func(context context.Context, parameters map[string]any, configs map[string]any) (result any, err error) {
		return nil, nil
	},
)

var read = mgc.NewRawStaticExecute(
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
	func(context context.Context, parameters map[string]any, configs map[string]any) (result any, err error) {
		return nil, nil
	},
)

var update = mgc.NewRawStaticExecute(
	"mock update",
	"v1",
	"",
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{
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
		[]string{},
	),
	nil,
	mgc.NewObjectSchema(
		map[string]*mgc.Schema{
			"id": mgc.NewStringSchema(),
		},
		[]string{},
	),
	func(context context.Context, parameters map[string]any, configs map[string]any) (result any, err error) {
		return nil, nil
	},
)

var testCases = []testCase{
	{
		res: &MgcResource{create: create, read: read, update: update},
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
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
		testCase.res.readInputAttributes(ctx)
		if diff := deep.Equal(testCase.res.inputAttr, testCase.expectedInput); diff != nil {
			t.Errorf("MgcResource.readInputAttributes failed. Diff list: %v", diff)
		}
		testCase.res.readOutputAttributes(ctx)
		if diff := deep.Equal(testCase.res.outputAttr, testCase.expectedOutput); diff != nil {
			t.Errorf("MgcResource.readOutputAttributes failed. Diff list: %v", diff)
		}
		finalAttr, _ := testCase.res.generateTFAttributes(ctx)
		if diff := deep.Equal(finalAttr, testCase.expectedFinal); diff != nil {
			t.Errorf("MgcResource.generateTFAttributes failed. Diff list: %v", diff)
		}
	}
}
