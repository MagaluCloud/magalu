package cmd

import (
	"reflect"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/spf13/cobra"
)

func TestTelemetryTreeNode(t *testing.T) {
	wantTelemetry := map[string]any{
		"name": telemetryCmdName,
		"children": []map[string]any{
			{"name": "disable"},
			{"name": "enable"},
			{"name": "status"},
		},
	}

	testCases := []struct {
		name          string
		withTelemetry bool
		wantFound     bool
		wantNode      map[string]any
	}{
		{"telemetry registered in the root", true, true, wantTelemetry},
		{"root without telemetry", false, false, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root := &cobra.Command{Use: "mgc"}
			if tc.withTelemetry {
				root.AddCommand(newTelemetryCmd(func() *telemetry.Service { return nil }))
			}

			node, found := telemetryTreeNode(root)

			if found != tc.wantFound {
				t.Fatalf("found = %v, want %v", found, tc.wantFound)
			}
			if !reflect.DeepEqual(node, tc.wantNode) {
				t.Errorf("got  %v\nwant %v", node, tc.wantNode)
			}
		})
	}
}
