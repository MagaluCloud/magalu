package openapi

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func testCompatLogger() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

func TestAddLegacyRoutePortID(t *testing.T) {
	tests := []struct {
		name     string
		body     map[string]core.Value
		expected map[string]core.Value
	}{
		{
			name: "derives port_id from a port target",
			body: map[string]core.Value{
				"cidr_destination": "192.168.1.0/24",
				"targets":          map[string]any{"id": "port-uuid", "type": "port_id"},
			},
			expected: map[string]core.Value{
				"cidr_destination": "192.168.1.0/24",
				"targets":          map[string]any{"id": "port-uuid", "type": "port_id"},
				"port_id":          "port-uuid",
			},
		},
		{
			// The comparison is lenient, but the value sent to the API is not
			// normalized: "targets" must go on the wire exactly as typed.
			name: "matches the target type regardless of case and preserves it",
			body: map[string]core.Value{
				"targets": map[string]any{"id": "port-uuid", "type": "PORT_ID"},
			},
			expected: map[string]core.Value{
				"targets": map[string]any{"id": "port-uuid", "type": "PORT_ID"},
				"port_id": "port-uuid",
			},
		},
		{
			// vpc_peering has no legacy representation, so it must reach old
			// regions unadorned and fail there.
			name: "ignores targets that are not ports",
			body: map[string]core.Value{
				"targets": map[string]any{"id": "peering-uuid", "type": "vpc_peering"},
			},
			expected: map[string]core.Value{
				"targets": map[string]any{"id": "peering-uuid", "type": "vpc_peering"},
			},
		},
		{
			// The case while the API contract is not migrated yet: the shim is
			// inert until "targets" shows up in the request body.
			name:     "ignores a body without targets",
			body:     map[string]core.Value{"cidr_destination": "192.168.1.0/24"},
			expected: map[string]core.Value{"cidr_destination": "192.168.1.0/24"},
		},
		{
			name:     "ignores null targets",
			body:     map[string]core.Value{"targets": nil},
			expected: map[string]core.Value{"targets": nil},
		},
		{
			name:     "ignores an empty target",
			body:     map[string]core.Value{"targets": map[string]any{}},
			expected: map[string]core.Value{"targets": map[string]any{}},
		},
		{
			name:     "ignores targets that are not an object",
			body:     map[string]core.Value{"targets": []any{map[string]any{"id": "port-uuid", "type": "port_id"}}},
			expected: map[string]core.Value{"targets": []any{map[string]any{"id": "port-uuid", "type": "port_id"}}},
		},
		{
			name:     "ignores a target without a type",
			body:     map[string]core.Value{"targets": map[string]any{"id": "port-uuid"}},
			expected: map[string]core.Value{"targets": map[string]any{"id": "port-uuid"}},
		},
		{
			name:     "ignores a target whose type is not a string",
			body:     map[string]core.Value{"targets": map[string]any{"id": "port-uuid", "type": 1}},
			expected: map[string]core.Value{"targets": map[string]any{"id": "port-uuid", "type": 1}},
		},
		{
			name:     "ignores a target without an id",
			body:     map[string]core.Value{"targets": map[string]any{"type": "port_id"}},
			expected: map[string]core.Value{"targets": map[string]any{"type": "port_id"}},
		},
		{
			name:     "ignores a target with an empty id",
			body:     map[string]core.Value{"targets": map[string]any{"id": "", "type": "port_id"}},
			expected: map[string]core.Value{"targets": map[string]any{"id": "", "type": "port_id"}},
		},
		{
			name:     "ignores a target whose id is not a string",
			body:     map[string]core.Value{"targets": map[string]any{"id": 1, "type": "port_id"}},
			expected: map[string]core.Value{"targets": map[string]any{"id": 1, "type": "port_id"}},
		},
		{
			name: "never overwrites an explicit port_id",
			body: map[string]core.Value{
				"port_id": "explicit-port",
				"targets": map[string]any{"id": "derived-port", "type": "port_id"},
			},
			expected: map[string]core.Value{
				"port_id": "explicit-port",
				"targets": map[string]any{"id": "derived-port", "type": "port_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addLegacyRoutePortID(testCompatLogger(), tt.body)
			assert.Equal(t, tt.expected, tt.body)
		})
	}
}

// The unit tests above assert what the shim does with a body shaped by hand,
// which is worth nothing if the real spec shapes it differently. This one goes
// through the production path -- embedded spec, newRequestBody, create() -- so
// that a "targets" whose schema drifts (object to array, renamed enum value)
// fails here instead of only against a non-migrated region.
func TestAddLegacyRoutePortIDAgainstEmbeddedSpec(t *testing.T) {
	const (
		routePath   = "/v1/vpcs/{vpc_id}/route_table/routes"
		portID      = "899672c1-bd54-4c98-9363-c886d00fc2a5"
		cidr        = "10.13.14.0/24"
		specName    = "openapis/network.openapi.yaml"
		extPrefixed = "x-mgc"
	)

	data, err := folder.ReadFile(specName)
	require.NoError(t, err)

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	require.NoError(t, err)

	pathItem := doc.Paths.Find(routePath)
	require.NotNil(t, pathItem, "%s is not in %s", routePath, specName)
	require.NotNil(t, pathItem.Post, "%s has no POST", routePath)

	extensionPrefix := extPrefixed
	rb := newRequestBody("POST", pathItem.Post, testCompatLogger(), &extensionPrefix)

	// Mirrors: --targets.id=<uuid> --targets.type port_id --cidr-destination <cidr>
	// (the kebab-case flag names belong to the CLI layer; here the parameter
	// names are the schema's own).
	mimeType, _, reader, requestBody, err := rb.create(core.Parameters{
		"cidr_destination": cidr,
		"targets":          map[string]any{"id": portID, "type": "port_id"},
	})
	require.NoError(t, err)
	assert.Equal(t, "application/json", mimeType)

	raw, err := io.ReadAll(reader)
	require.NoError(t, err)

	var wire map[string]any
	require.NoError(t, json.Unmarshal(raw, &wire))

	// Non-migrated regions read this one, and reject the request without it.
	assert.Equal(t, portID, wire["port_id"], "port_id was not derived; wire body: %s", raw)
	// Migrated regions read this one, untouched.
	assert.Equal(t, map[string]any{"id": portID, "type": "port_id"}, wire["targets"])
	assert.Equal(t, cidr, wire["cidr_destination"])

	// create() also returns the body for link expressions and debug logs; it
	// must agree with what actually went on the wire.
	assert.Equal(t, portID, requestBody.(map[string]core.Value)["port_id"])
}

// Every shim is keyed by operationId, so a spec refresh that renames or removes
// the operation would silently disable it. Fail here instead.
func TestRequestBodyCompatShimsExistInEmbeddedSpecs(t *testing.T) {
	operationIDs := map[string]string{}

	entries, err := folder.ReadDir("openapis")
	assert.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "index.openapi.yaml" {
			continue
		}

		data, err := folder.ReadFile("openapis/" + entry.Name())
		assert.NoError(t, err)

		var doc struct {
			Paths map[string]map[string]any `yaml:"paths"`
		}
		assert.NoError(t, yaml.Unmarshal(data, &doc), entry.Name())

		for _, operations := range doc.Paths {
			for _, operation := range operations {
				operationMap, ok := operation.(map[string]any)
				if !ok {
					continue
				}
				if id, ok := operationMap["operationId"].(string); ok {
					operationIDs[id] = entry.Name()
				}
			}
		}
	}

	assert.NotEmpty(t, operationIDs, "no operationId found in the embedded specs")

	for operationID := range requestBodyCompatShims {
		if _, ok := operationIDs[operationID]; !ok {
			t.Errorf(
				"request body compatibility shim is keyed by operationId %q, which is not in the embedded specs. "+
					"Either a spec refresh renamed/removed the operation and the shim is now silently disabled, "+
					"or the operation is gone and the shim should be deleted.",
				operationID,
			)
		}
	}
}
