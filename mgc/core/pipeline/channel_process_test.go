package pipeline_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/pipeline"
)

func double(ctx context.Context, value int) (int, pipeline.ProcessStatus) {
	return value * 2, pipeline.ProcessOutput
}

func TestSimplePipeline(t *testing.T) {
	inputSize := 4
	ctx := context.Background()
	genChan := pipeline.RangeGenerator(ctx, inputSize)

	pipe := pipeline.Process(ctx, genChan, double, nil)
	pipe = pipeline.Process(ctx, pipe, double, nil)

	checks := map[int]bool{
		0:  false,
		4:  false,
		8:  false,
		12: false,
	}
	for result := range pipe {
		checks[result] = true
	}
	for k, v := range checks {
		if !v {
			t.Error("Not all entries numbers were found: ", k)
		}
	}
}

func TestParallelPipeline(t *testing.T) {
	inputSize := 4
	ctx := context.Background()
	genChan := pipeline.RangeGenerator(ctx, inputSize)

	pipe := pipeline.ParallelProcess(ctx, 4, genChan, double, nil)
	pipe = pipeline.ParallelProcess(ctx, 4, pipe, double, nil)

	checks := map[int]bool{
		0:  false,
		4:  false,
		8:  false,
		12: false,
	}
	for result := range pipe {
		checks[result] = true
	}
	for k, v := range checks {
		if !v {
			t.Error("Not all entries numbers were found: ", k)
		}
	}
}

func TestSimplePipelineWithSlice(t *testing.T) {
	ctx := context.Background()
	genChan := pipeline.SliceItemGenerator(ctx, []int{1, 2, 3})

	pipe := pipeline.Process(ctx, genChan, double, nil)
	pipe = pipeline.Process(ctx, pipe, double, nil)

	checks := map[int]bool{
		4:  false,
		8:  false,
		12: false,
	}
	for result := range pipe {
		checks[result] = true
	}
	for k, v := range checks {
		if !v {
			t.Error("Not all entries numbers were found: ", k)
		}
	}
}

func TestBatchPipeline(t *testing.T) {
	ctx := context.Background()
	genChan := pipeline.RangeGenerator(ctx, 12)
	batchSize := 3

	batches := pipeline.Batch(ctx, genChan, batchSize)
	checks := map[string]bool{
		"0,1,2":   false,
		"3,4,5":   false,
		"6,7,8":   false,
		"9,10,11": false,
	}
	for batch := range batches {
		batchStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(batch)), ","), "[]")
		checks[batchStr] = true
	}
	for k, v := range checks {
		if !v {
			t.Error("Not all entries numbers were found: ", k)
		}
	}
}
