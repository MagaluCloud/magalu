package cmd

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/core"

	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/objects"
)

func TestSyncBatchSizeFlag(t *testing.T) {
	descriptor, err := objects.GetGroup().GetChildByName("sync")
	if err != nil {
		t.Fatal(err)
	}
	executor := descriptor.(core.Executor)
	schema := executor.ParametersSchema()
	for _, tc := range []struct {
		name, value string
		want        int
		wantErr     bool
	}{
		{name: "omitted", want: 1000},
		{name: "minimum", value: "1", want: 1},
		{name: "maximum", value: "1000", want: 1000},
		{name: "zero", value: "0", wantErr: true},
		{name: "negative", value: "-1", wantErr: true},
		{name: "above maximum", value: "1001", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := schema_flags.NewSchemaFlag(schema, "batch_size", "batch-size", slices.Contains(schema.Required, "batch_size"), false, false)
			if tc.value != "" {
				if err := f.Value.Set(tc.value); err != nil {
					t.Fatal(err)
				}
			}
			value, err := schema_flags.GetFlagValue(f, nil)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected invalid batch size")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			var got int
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}
