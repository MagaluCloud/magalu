package openapi

import "context"

type RawOutput bool

type RawOutputKey struct{}

func GetRawOutputFlag(ctx context.Context) bool {
	raw, found := ctx.Value(RawOutputKey{}).(RawOutput)
	if !found {
		return true
	}
	return bool(raw)
}

func WithRawOutputFlag(ctx context.Context, raw bool) context.Context {
	return context.WithValue(ctx, RawOutputKey{}, RawOutput(raw))
}
