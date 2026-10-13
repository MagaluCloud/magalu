package cmdtelemetry

import (
	"errors"
	"testing"

	flag "github.com/spf13/pflag"
)

func TestWrapFlagErrorKeepsCause(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want error
	}{
		{name: "nil stays nil", err: nil, want: nil},
		{name: "help request still detected", err: flag.ErrHelp, want: flag.ErrHelp},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := WrapFlagError(nil, tc.err)
			if tc.want == nil {
				if got != nil {
					t.Errorf("got %v, want nil", got)
				}
				return
			}
			if !errors.Is(got, tc.want) {
				t.Errorf("errors.Is(%v, %v) = false", got, tc.want)
			}
		})
	}
}
