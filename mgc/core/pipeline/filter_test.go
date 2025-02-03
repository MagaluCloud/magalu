package pipeline_test

import (
	"context"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
	"golang.org/x/exp/constraints"
)

type EvenFilter[T constraints.Integer] struct{}

func (f EvenFilter[T]) Filter(ctx context.Context, entry T) pipeline.FilterStatus {
	if int64(entry)%2 == 0 {
		return pipeline.FilterInclude
	} else {
		return pipeline.FilterExclude
	}
}

type ThreeFilter[T constraints.Integer] struct{}

func (f ThreeFilter[T]) Filter(ctx context.Context, entry T) pipeline.FilterStatus {
	if int64(entry)%3 == 0 {
		return pipeline.FilterInclude
	} else {
		return pipeline.FilterExclude
	}
}

var OddFilterInstance = pipeline.FilterRuleNot[int]{Not: EvenFilter[int]{}}

func TestFilter(t *testing.T) {
	ctx := context.Background()
	genChan := pipeline.RangeGenerator(ctx, 10)

	filteredChan := pipeline.Filter[int](ctx, genChan, EvenFilter[int]{})

	for num := range filteredChan {
		if num%2 != 0 {
			t.Error("Found odd number after EvenFilter: ", num)
		}
	}
}

func TestNotFilter(t *testing.T) {
	ctx := context.Background()
	genChan := pipeline.RangeGenerator(ctx, 10)

	filteredChan := pipeline.Filter[int](ctx, genChan, OddFilterInstance)

	for num := range filteredChan {
		if num%2 == 0 {
			t.Error("Found even number after OddFilter: ", num)
		}
	}
}

func TestAllFilter(t *testing.T) {
	// TODO this test has been started but is not working properly
	t.Skip()
	ctx := context.Background()
	genChan := pipeline.RangeGenerator(ctx, 15)

	filter := pipeline.FilterRuleAll[int]{
		All: []pipeline.FilterRule[int]{
			EvenFilter[int]{},
			ThreeFilter[int]{},
		},
	}

	filteredChan := pipeline.Filter[int](ctx, genChan, filter)

	for num := range filteredChan {
		if num%2 != 0 || num%3 != 0 {
			t.Error("Found non multiple of 6 after AllFilter: ", num)
		}
	}
}
