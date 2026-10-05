package cmd

import (
	"errors"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	flag "github.com/spf13/pflag"
)

func TestInvalidFlagError(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "unknown flag on existing command",
			args: []string{"virtual-machine", "instances", "create", "--flag-inexistente"},
		},
		{
			name: "flag without required value",
			args: []string{"virtual-machine", "instances", "create", "--region"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root := newFakeCommandTree()
			root.SetArgs(tc.args)
			untyped := root.Execute()
			if untyped == nil {
				t.Fatal("expected a flag error")
			}

			root = newFakeCommandTree()
			root.SetFlagErrorFunc(wrapFlagError)
			root.SetArgs(tc.args)
			err := root.Execute()

			if got, want := err.Error(), untyped.Error(); got != want {
				t.Errorf("message changed: got %q, want %q", got, want)
			}
			if got := telemetry.ClassifyError(err); got != telemetry.FailureValidation {
				t.Errorf("ClassifyError = %q, want %q", got, telemetry.FailureValidation)
			}
			if _, track := telemetryCommandInfo(root, tc.args, err); !track {
				t.Error("flag error must generate an event")
			}
		})
	}
}

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
			got := wrapFlagError(nil, tc.err)
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
