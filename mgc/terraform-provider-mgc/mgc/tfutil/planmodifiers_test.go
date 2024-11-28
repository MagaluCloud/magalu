package tfutil

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestReplaceIfChangeAndNotIsNotSetOnPlan_PlanModifyString(t *testing.T) {
	tests := []struct {
		name            string
		stateValue      types.String
		planValue       types.String
		expectedReplace bool
	}{
		{
			name:            "plan_value_unknown",
			stateValue:      types.StringValue("zone-1"),
			planValue:       types.StringUnknown(),
			expectedReplace: false,
		},
		{
			name:            "values_different",
			stateValue:      types.StringValue("zone-1"),
			planValue:       types.StringValue("zone-2"),
			expectedReplace: true,
		},
		{
			name:            "values_equal",
			stateValue:      types.StringValue("zone-1"),
			planValue:       types.StringValue("zone-1"),
			expectedReplace: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := ReplaceIfChangeAndNotIsNotSetOnPlan{}
			request := planmodifier.StringRequest{
				StateValue: tt.stateValue,
				PlanValue:  tt.planValue,
			}
			response := &planmodifier.StringResponse{}

			modifier.PlanModifyString(context.Background(), request, response)

			assert.Equal(t, tt.expectedReplace, response.RequiresReplace)
		})
	}
}

func TestReplaceIfChangeAndNotIsNotSetOnPlan_Description(t *testing.T) {
	modifier := ReplaceIfChangeAndNotIsNotSetOnPlan{}
	expected := "Requires replace if the value is different from the state and the plan value is not unknown."

	assert.Equal(t, expected, modifier.Description(context.Background()))
}

func TestReplaceIfChangeAndNotIsNotSetOnPlan_MarkdownDescription(t *testing.T) {
	modifier := ReplaceIfChangeAndNotIsNotSetOnPlan{}
	expected := "Requires replace if the value is different from the state and the plan value is not unknown."

	assert.Equal(t, expected, modifier.MarkdownDescription(context.Background()))
}
