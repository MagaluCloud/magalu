package common

import (
	"context"

	"magalu.cloud/core/pipeline"
)

type Filters struct {
	FilterParams []FilterParams `json:"filter,omitempty" jsonschema:"description=File name pattern to include or exclude"`
}

type FilterParams struct {
	Include string `json:"include,omitempty" jsonschema:"description=Filename pattern to include"`
	Exclude string `json:"exclude,omitempty" jsonschema:"description=Filename pattern to exclude"`
}

func ApplyFilters(ctx context.Context, entries <-chan pipeline.WalkDirEntry, params []FilterParams, cancel context.CancelCauseFunc) <-chan pipeline.WalkDirEntry {
	filters := []pipeline.FilterRule[pipeline.WalkDirEntry]{}
	for _, filter := range params {
		if filter.Include != "" {
			filters = append(filters, pipeline.FilterWalkDirEntryIncludeGlobMatch{
				Pattern: filter.Include, CancelOnError: cancel,
			})
		}
		if filter.Exclude != "" {
			filters = append(filters, pipeline.FilterRuleNot[pipeline.WalkDirEntry]{
				Not: pipeline.FilterWalkDirEntryIncludeGlobMatch{Pattern: filter.Exclude, CancelOnError: cancel},
			})
		}
	}

	if len(filters) < 1 {
		return entries
	}

	filterRule := pipeline.FilterRuleFirst[pipeline.WalkDirEntry]{Filters: filters}
	return pipeline.Filter[pipeline.WalkDirEntry](ctx, entries, filterRule)
}
