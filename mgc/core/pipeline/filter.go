package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
)

type FilterStatus int

const (
	// keep going until an explicit FilterInclude or FilterExclude
	FilterUnknown FilterStatus = iota
	// The entry should be included in the output
	FilterInclude
	// The entry should NOT be included in the output
	FilterExclude
)

// Use an interface and then structs so debugging is simpler, we can get the tree
// and actual values instead of plain function names (or anonymous)

// Given an entry, say if it should be included or excluded from the output.
//
// If FilterUnknown is used, the filtering will keep going.
// Ultimately FilterUnknown values will be included.
type FilterRule[T any] interface {
	Filter(ctx context.Context, entry T) FilterStatus
}

type FilterRuleRecursiveWrapper[T any] interface {
	// Wrap any children with the given wrapper and return a new instance of the rule
	RecursiveWrap(wrapper filterRuleWrapper[T]) FilterRule[T]
}

type filterRuleWrapper[T any] func(FilterRule[T]) FilterRule[T]

func wrapFilterRuleLogArray[T any](filterRules []FilterRule[T], wrapper filterRuleWrapper[T]) []FilterRule[T] {
	wrapped := make([]FilterRule[T], len(filterRules))
	for i, f := range filterRules {
		wrapped[i] = wrapper(f)
	}
	return wrapped
}

// Inverts FilterInclude<->FilterExclude, keeps FilterUnknown
type FilterRuleNot[T any] struct {
	Not FilterRule[T]
}

func (r FilterRuleNot[T]) Filter(ctx context.Context, entry T) FilterStatus {
	switch r.Not.Filter(ctx, entry) {
	default:
		return FilterUnknown
	case FilterInclude:
		return FilterExclude
	case FilterExclude:
		return FilterInclude
	}
}

func (r FilterRuleNot[T]) RecursiveWrap(wrapper filterRuleWrapper[T]) FilterRule[T] {
	wrapped := wrapper(r.Not)
	return &FilterRuleNot[T]{wrapped}
}

var _ FilterRule[any] = (*FilterRuleNot[any])(nil)
var _ FilterRuleRecursiveWrapper[any] = (*FilterRuleNot[any])(nil)

// All must match (either FilterInclude or FilterExclude, FilterUnknown is skipped)
type FilterRuleAll[T any] struct {
	All []FilterRule[T]
}

func (r FilterRuleAll[T]) Filter(ctx context.Context, entry T) FilterStatus {
	status := FilterUnknown
	for _, f := range r.All {
		current := f.Filter(ctx, entry)
		if status == FilterUnknown {
			status = current
		} else if status != current {
			return FilterUnknown
		}
	}
	return status
}

func (r FilterRuleAll[T]) RecursiveWrap(wrapper filterRuleWrapper[T]) FilterRule[T] {
	return &FilterRuleAll[T]{wrapFilterRuleLogArray[T](r.All, wrapper)}
}

var _ FilterRule[any] = (*FilterRuleAll[any])(nil)
var _ FilterRuleRecursiveWrapper[any] = (*FilterRuleAll[any])(nil)

// Any known value is used (either FilterInclude or FilterExclude, FilterUnknown is skipped)
type FilterRuleAny[T any] struct {
	Any []FilterRule[T]
}

func (r FilterRuleAny[T]) Filter(ctx context.Context, entry T) FilterStatus {
	for _, f := range r.Any {
		current := f.Filter(ctx, entry)
		if current == FilterUnknown {
			continue
		}
		return current
	}
	return FilterUnknown
}

func (r FilterRuleAny[T]) RecursiveWrap(wrapper filterRuleWrapper[T]) FilterRule[T] {
	return &FilterRuleAny[T]{wrapFilterRuleLogArray[T](r.Any, wrapper)}
}

var _ FilterRule[any] = (*FilterRuleAny[any])(nil)
var _ FilterRuleRecursiveWrapper[any] = (*FilterRuleAny[any])(nil)

// Log the given rule
//
// ContextLogger() is used to get the logger
type FilterRuleLog[T any] struct {
	Rule FilterRule[T]
}

func (r FilterRuleLog[T]) Filter(ctx context.Context, entry T) (status FilterStatus) {
	status = r.Rule.Filter(ctx, entry)
	logger := FromContext(ctx)
	var statusStr string
	switch status {
	case FilterExclude:
		statusStr = "exclude"
	case FilterInclude:
		statusStr = "include"
	case FilterUnknown:
		statusStr = "unknown"
	default:
		statusStr = fmt.Sprintf("unknown status code: %d", status)
	}
	logger.Debugw("filter", "rule", r.Rule, "entry", entry, "status", status, "statusStr", statusStr)
	return
}

func (r FilterRuleLog[T]) MarshalJSON() ([]byte, error) {
	// Omit itself from the marshaling, helps when is the rule tree is given as log context values
	return json.Marshal(r.Rule)
}

var _ FilterRule[any] = (*FilterRuleLog[any])(nil)
var _ json.Marshaler = (*FilterRuleLog[any])(nil)

// Recursively wraps all required filter elements with logging
//
// Rules may implement FilterRuleRecursiveWrapper in order to keep the recursion going.
func RecursiveFilterRuleLog[T any](filterRule FilterRule[T]) FilterRule[T] {
	switch f := filterRule.(type) {
	case FilterRuleLog[T], *FilterRuleLog[T]:
		return filterRule

	case FilterRuleRecursiveWrapper[T]:
		return &FilterRuleLog[T]{f.RecursiveWrap(RecursiveFilterRuleLog[T])}

	default:
		return &FilterRuleLog[T]{filterRule}
	}
}

// Filter an input channel
//
// Filtering is done on an un-buffered channel, one by one, so it will
// block until the next item can be consumed.
//
// Filtering may be early stopped by context.Context.Done(), see
// context.WithCancel(), context.WithTimeout() and context.WithDeadline()
func Filter[T any](
	ctx context.Context,
	inputChan <-chan T,
	filterRule FilterRule[T],
) <-chan T {
	logger := FromContext(ctx).Named("Filter").With(
		"filterRule", filterRule,
	)
	ctx = NewContext(ctx, logger)

	return Process[T, T](ctx, inputChan, func(ctx context.Context, input T) (output T, status ProcessStatus) {
		if filterRule.Filter(ctx, input) == FilterExclude {
			return input, ProcessSkip
		}
		return input, ProcessOutput
	}, nil)
}

// Utilities to filter entries by name (glob, regular expressions) and type (file/directory)

type filterWalkDirEntryIncludeFile struct{}

func (r filterWalkDirEntryIncludeFile) Filter(ctx context.Context, entry WalkDirEntry) FilterStatus {
	if entry.DirEntry().Type().IsRegular() {
		return FilterInclude
	}
	return FilterUnknown
}

func (r filterWalkDirEntryIncludeFile) MarshalJSON() ([]byte, error) {
	return ([]byte)("\"FilterWalkDirEntryIncludeFile\""), nil
}

var _ FilterRule[WalkDirEntry] = (*filterWalkDirEntryIncludeFile)(nil)
var _ json.Marshaler = (*filterWalkDirEntryIncludeFile)(nil)

// Include regular files, others are left as unknown
var FilterWalkDirEntryIncludeFile = &filterWalkDirEntryIncludeFile{}

// Exclude regular files, others are left as unknown
var FilterWalkDirEntryExcludeFile = &FilterRuleNot[WalkDirEntry]{FilterWalkDirEntryIncludeFile}

type filterWalkDirEntryIncludeDir struct{}

func (r filterWalkDirEntryIncludeDir) Filter(ctx context.Context, entry WalkDirEntry) FilterStatus {
	if entry.DirEntry().IsDir() {
		return FilterInclude
	}
	return FilterUnknown
}

func (r filterWalkDirEntryIncludeDir) MarshalJSON() ([]byte, error) {
	return ([]byte)("\"FilterWalkDirEntryIncludeDir\""), nil
}

var _ FilterRule[WalkDirEntry] = (*filterWalkDirEntryIncludeDir)(nil)
var _ json.Marshaler = (*filterWalkDirEntryIncludeDir)(nil)

// Include directories, others are left as unknown
var FilterWalkDirEntryIncludeDir = &filterWalkDirEntryIncludeDir{}

// Exclude directories, others are left as unknown
var FilterWalkDirEntryExcludeDir = &FilterRuleNot[WalkDirEntry]{FilterWalkDirEntryIncludeDir}

// Include entries if name matches the regular expression
//
// # Note that only file names (basename) are matched, not the whole path
//
// Non-matching files are NOT excluded, they are just left as unknown.
// to exclude them, use &FilterRuleNot[WalkDirEntry]{FilterWalkDirEntryIncludeRegExp{...}}
type FilterWalkDirEntryIncludeRegExp struct {
	Regexp *regexp.Regexp
}

func (r FilterWalkDirEntryIncludeRegExp) Filter(ctx context.Context, entry WalkDirEntry) FilterStatus {
	if r.Regexp.MatchString(entry.DirEntry().Name()) {
		return FilterInclude
	}
	return FilterUnknown
}

func (r FilterWalkDirEntryIncludeRegExp) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"Regexp": r.Regexp.String()})
}

var _ FilterRule[WalkDirEntry] = (*FilterWalkDirEntryIncludeRegExp)(nil)
var _ json.Marshaler = (*FilterWalkDirEntryIncludeRegExp)(nil)

// Include entries if name matches the glob pattern
//
// # Note that only file names (basename) are matched, not the whole path
//
// Non-matching files are NOT excluded, they are just left as unknown.
// to exclude them, use &FilterRuleNot[WalkDirEntry]{FilterWalkDirEntryIncludeGlobMatch{...}}
type FilterWalkDirEntryIncludeGlobMatch struct {
	Pattern string
}

func (r FilterWalkDirEntryIncludeGlobMatch) Filter(ctx context.Context, entry WalkDirEntry) FilterStatus {
	match, _ := filepath.Match(r.Pattern, entry.DirEntry().Name())
	if match {
		return FilterInclude
	}
	return FilterUnknown
}

var _ FilterRule[WalkDirEntry] = (*FilterWalkDirEntryIncludeGlobMatch)(nil)
