package tfutil

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type ReplaceIfChangeAndNotIsNotSetOnPlan struct{}

func (p ReplaceIfChangeAndNotIsNotSetOnPlan) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	if request.PlanValue.IsUnknown() {
		response.RequiresReplace = false
		return
	}
	if !request.StateValue.Equal(request.PlanValue) {
		response.RequiresReplace = true
		return
	}
	response.RequiresReplace = false
}

func (p ReplaceIfChangeAndNotIsNotSetOnPlan) Description(context.Context) string {
	return "Requires replace if the value is different from the state and the plan value is not unknown."
}

func (p ReplaceIfChangeAndNotIsNotSetOnPlan) MarkdownDescription(context.Context) string {
	return "Requires replace if the value is different from the state and the plan value is not unknown."
}

type ReplaceIfChangeAndNotIsNotSetOnPlanBool struct{}

func (p ReplaceIfChangeAndNotIsNotSetOnPlanBool) PlanModifyBool(ctx context.Context, request planmodifier.BoolRequest, response *planmodifier.BoolResponse) {
	if request.PlanValue.IsUnknown() {
		response.RequiresReplace = false
		return
	}
	if !request.StateValue.Equal(request.PlanValue) {
		response.RequiresReplace = true
		return
	}
	response.RequiresReplace = false
}

func (p ReplaceIfChangeAndNotIsNotSetOnPlanBool) Description(context.Context) string {
	return "Requires replace if the value is different from the state and the plan value is not unknown."
}

func (p ReplaceIfChangeAndNotIsNotSetOnPlanBool) MarkdownDescription(context.Context) string {
	return "Requires replace if the value is different from the state and the plan value is not unknown."
}
