package cmd

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestSplitUnquoted(t *testing.T) {
	testCases := []struct {
		in       string
		out      []string
		maxSplit int
	}{
		{`a:b`, []string{"a", "b"}, 0},
		{`a":b"`, []string{"a:b"}, 0},
		{`"a:b"`, []string{"a:b"}, 0},
		{`"a":"b"`, []string{"a", "b"}, 0},
		{`"a:\"b"`, []string{`a:"b`}, 0},
		{``, []string{}, 0},
		{`a:b:c`, []string{"a", "b", "c"}, 0},
		{`a:b:c`, []string{"a:b:c"}, 1},
		{`a:b:c`, []string{"a", "b:c"}, 2},
		{`a:b:c`, []string{"a", "b", "c"}, 3},
	}
	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Convert \"%s\" to %#v", tc.in, tc.out),
			func(t *testing.T) {
				result, err := splitUnquoted(tc.in, ":", tc.maxSplit)
				if err != nil {
					t.Errorf("Expected function to work, found error: %v", err)
					return
				}
				if !reflect.DeepEqual(tc.out, result) {
					t.Errorf("Convesion failed: for input %s, expected %#v, got %#v", tc.in, tc.out, result)
				}
			},
		)
	}
}

func TestSplitUnquotedShouldFail(t *testing.T) {
	testCases := []string{
		`a:"b`,
		`"a:b`,
	}
	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Should fail to convert \"%s\"", tc),
			func(t *testing.T) {
				_, err := splitUnquoted(tc, ":", 0)
				if err == nil {
					t.Error("Function should error and did not")
					return
				}
				if !errors.Is(err, strconv.ErrSyntax) {
					t.Errorf("Expected %v, found %v", strconv.ErrSyntax, err)
				}
			},
		)
	}
}
