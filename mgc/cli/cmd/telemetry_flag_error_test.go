package cmd

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/failure"
	"github.com/MagaluCloud/magalu/mgc/cli/cmd/telemetryinstrument"
	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

// TestInvalidFlagError fica em cmd, e não em telemetryinstrument, porque depende de
// newFakeCommandTree e telemetryinstrument.BuildCommandInfo(..., getNextUnknownCommand), que por
// sua vez depende de normalizeFlagName e applyPositionalArgs, privadas do cmd.
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
			root.SetFlagErrorFunc(failure.WrapFlagError)
			root.SetArgs(tc.args)
			err := root.Execute()

			if got, want := err.Error(), untyped.Error(); got != want {
				t.Errorf("message changed: got %q, want %q", got, want)
			}
			if got := telemetry.ClassifyError(err); got != telemetry.FailureValidation {
				t.Errorf("ClassifyError = %q, want %q", got, telemetry.FailureValidation)
			}
			if _, track := telemetryinstrument.BuildCommandInfo(root, tc.args, err, nil, getNextUnknownCommand); !track {
				t.Error("flag error must generate an event")
			}
		})
	}
}
