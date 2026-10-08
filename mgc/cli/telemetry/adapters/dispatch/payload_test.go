package dispatch

import (
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"testing"
)

func TestPayloadRoundTrip(t *testing.T) {
	testCases := []struct {
		name         string
		debugLogPath string
	}{
		{"without debug", ""},
		{"with debug log", "/home/user/.config/mgc/telemetry-debug.txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := testEvent()

			data, err := Encode(event, tc.debugLogPath)
			if err != nil {
				t.Fatal(err)
			}
			payload, err := Decode(data)
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}

			if payload.DebugLogPath != tc.debugLogPath {
				t.Errorf("debug log path = %q, want %q", payload.DebugLogPath, tc.debugLogPath)
			}
			if payload.Event.InstallationID != event.InstallationID || payload.Event.Product != event.Product {
				t.Errorf("installation id and product = %q %q, want %q %q",
					payload.Event.InstallationID, payload.Event.Product, event.InstallationID, event.Product)
			}
			if !reflect.DeepEqual(payload.Event, event) {
				t.Errorf("event changed in the round trip\ngot  %+v\nwant %+v", payload.Event, event)
			}
		})
	}
}

func TestPayloadDecodeErrors(t *testing.T) {
	valid, err := Encode(testEvent(), "")
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(valid, &fields); err != nil {
		t.Fatal(err)
	}
	fields["version"] = payloadVersion + 1
	futureVersion, _ := json.Marshal(fields)

	testCases := []struct {
		name        string
		data        string
		wantVersion bool
	}{
		{"unknown version", string(futureVersion), true},
		{"missing version", `{"event":{"timestamp":"2026-09-22T14:30:45.000Z"}}`, true},
		{"invalid json", `{"version":`, false},
		{"invalid event", `{"version":1,"event":{"timestamp":"yesterday"}}`, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Decode([]byte(tc.data))
			if err == nil {
				t.Fatalf("Decode(%s) = nil, want an error", tc.data)
			}
			if got := errors.Is(err, ErrUnknownPayloadVersion); got != tc.wantVersion {
				t.Errorf("errors.Is(ErrUnknownPayloadVersion) = %v, want %v (err %v)", got, tc.wantVersion, err)
			}
		})
	}
}

// TestPayloadAddsNothingToTheEvent garante que o payload não acrescenta dados ao
// evento. O campo event é o JSON do modelo, cuja anonimização é testada no núcleo, e
// fora dele só existem as chaves conhecidas do protocolo
func TestPayloadAddsNothingToTheEvent(t *testing.T) {
	testCases := []struct {
		name         string
		debugLogPath string
		wantKeys     []string
	}{
		{"without debug", "", []string{"event", "installationId", "product", "version"}},
		{"with debug", "/home/user/.config/mgc/telemetry-debug.txt", []string{"debugLogPath", "event", "installationId", "product", "version"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := testEvent()
			data, err := Encode(event, tc.debugLogPath)
			if err != nil {
				t.Fatal(err)
			}

			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			keys := make([]string, 0, len(fields))
			for key := range fields {
				keys = append(keys, key)
			}
			slices.Sort(keys)
			if !slices.Equal(keys, tc.wantKeys) {
				t.Errorf("payload keys = %v, want %v", keys, tc.wantKeys)
			}

			model, err := json.Marshal(event)
			if err != nil {
				t.Fatal(err)
			}
			if string(fields["event"]) != string(model) {
				t.Errorf("event in the payload differs from the event model\npayload %s\nmodel   %s", fields["event"], model)
			}
		})
	}
}
