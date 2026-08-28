package cmd

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/spf13/cobra"
)

const (
	paginateFlag         = "cli.paginate"
	paginateLimitProp    = "_limit"
	paginateOffsetProp   = "_offset"
	paginateResultsField = "results"
	paginateDefaultLimit = 50
)

func addPaginateFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().Bool(
		paginateFlag,
		false,
		`Automatically loop through all pages of a paginated list operation, aggregating the results.
Only works on operations that expose '--control.limit' and '--control.offset'.`,
	)
}

func getPaginateFlag(cmd *cobra.Command) bool {
	v, err := cmd.Root().PersistentFlags().GetBool(paginateFlag)
	if err != nil {
		return false
	}
	return v
}

// supportsPagination reports whether the given parameters schema exposes the
// pagination knobs (`_limit` and `_offset`) the CLI loops over.
func supportsPagination(schema *mgcSchemaPkg.Schema) bool {
	if schema == nil || schema.Properties == nil {
		return false
	}
	_, hasLimit := schema.Properties[paginateLimitProp]
	_, hasOffset := schema.Properties[paginateOffsetProp]
	return hasLimit && hasOffset
}

// resolveLimit picks the page size to use during the auto-pagination loop.
// User-provided value wins; otherwise we use the schema's maximum, falling
// back to a conservative default.
func resolveLimit(schema *mgcSchemaPkg.Schema, parameters core.Parameters) int {
	if v, ok := parameters[paginateLimitProp]; ok {
		if n, ok := toInt(v); ok && n > 0 {
			return n
		}
	}
	if ref, ok := schema.Properties[paginateLimitProp]; ok && ref != nil && ref.Value != nil {
		if ref.Value.Max != nil && *ref.Value.Max > 0 {
			return int(*ref.Value.Max)
		}
	}
	return paginateDefaultLimit
}

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	case float32:
		return int(n), true
	}
	return 0, false
}

// extractPageItems pulls the array of items out of a single-page result value,
// supporting both the common `{"results": [...]}` shape and top-level arrays.
func extractPageItems(value any) (items []any, fieldName string, ok bool) {
	switch v := value.(type) {
	case []any:
		return v, "", true
	case map[string]any:
		if r, has := v[paginateResultsField]; has {
			if arr, isArr := r.([]any); isArr {
				return arr, paginateResultsField, true
			}
		}
	}
	return nil, "", false
}

// runPaginated repeatedly invokes exec, advancing `_offset` by `_limit` each
// iteration, until a page comes back with fewer items than the limit. The
// aggregated value mirrors the shape of the first page (top-level array, or
// object with a `results` field), so existing formatters and jsonpath queries
// keep working.
func runPaginated(
	ctx context.Context,
	exec core.Executor,
	parameters core.Parameters,
	configs core.Configs,
) (core.Result, error) {
	schema := exec.ParametersSchema()
	if !supportsPagination(schema) {
		return nil, core.UsageError{Err: fmt.Errorf(
			"--%s is only supported on operations that expose '--control.limit' and '--control.offset'",
			paginateFlag,
		)}
	}

	limit := resolveLimit(schema, parameters)

	loopParams := make(core.Parameters, len(parameters)+2)
	for k, v := range parameters {
		loopParams[k] = v
	}
	loopParams[paginateLimitProp] = limit
	loopParams[paginateOffsetProp] = 0

	var (
		aggregated []any
		lastResult core.Result
		fieldName  string
	)

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		result, err := exec.Execute(ctx, loopParams, configs)
		if err != nil {
			return nil, err
		}
		lastResult = result

		rv, ok := core.ResultAs[core.ResultWithValue](result)
		if !ok {
			// Non-value result (e.g. raw reader): nothing to aggregate.
			return result, nil
		}

		items, fname, ok := extractPageItems(rv.Value())
		if !ok {
			return nil, fmt.Errorf(
				"--%s: unexpected response shape from %q, expected an array or an object with a %q field",
				paginateFlag, exec.Name(), paginateResultsField,
			)
		}
		fieldName = fname
		aggregated = append(aggregated, items...)

		if len(items) < limit {
			break
		}

		offset, _ := toInt(loopParams[paginateOffsetProp])
		loopParams[paginateOffsetProp] = offset + limit
	}

	rv, _ := core.ResultAs[core.ResultWithValue](lastResult)
	merged := buildPaginatedValue(rv.Value(), aggregated, fieldName)

	source := lastResult.Source()
	source.Parameters = parameters
	return core.NewSimpleResult(source, rv.Schema(), merged), nil
}

// buildPaginatedValue rebuilds the result value preserving the original
// shape: a top-level array stays an array, an object keeps its other fields
// (e.g. `meta`) but with the merged `results` array.
func buildPaginatedValue(lastValue any, aggregated []any, fieldName string) any {
	if fieldName == "" {
		return aggregated
	}
	out := map[string]any{}
	if m, ok := lastValue.(map[string]any); ok {
		for k, v := range m {
			out[k] = v
		}
	}
	out[fieldName] = aggregated
	return out
}
