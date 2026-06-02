package pipeline

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNewDeeplTranslator(t *testing.T) {
	tests := []struct {
		apiKey           string
		expectedEndpoint string
	}{
		{"key123", "https://api.deepl.com/v2/translate"},
		{"key123:fx", "https://api-free.deepl.com/v2/translate"},
		{":fx", "https://api-free.deepl.com/v2/translate"},
	}
	for _, tt := range tests {
		t.Run(tt.apiKey, func(t *testing.T) {
			tr := newDeeplTranslator(tt.apiKey)
			if tr.endpoint != tt.expectedEndpoint {
				t.Errorf("got %q, want %q", tr.endpoint, tt.expectedEndpoint)
			}
		})
	}
}

func TestResolveType(t *testing.T) {
	tests := []struct {
		name     string
		prop     paramProperty
		expected string
	}{
		{"string", paramProperty{Type: "string"}, "string"},
		{"boolean", paramProperty{Type: "boolean"}, "boolean"},
		{"integer", paramProperty{Type: "integer"}, "integer"},
		{"array with items", paramProperty{Type: "array", Items: &propertyItems{Type: "string"}}, "array(string)"},
		{"array without items", paramProperty{Type: "array"}, "array"},
		{"anyOf", paramProperty{AnyOf: json.RawMessage(`[{}]`)}, "object"},
		{"oneOf", paramProperty{OneOf: json.RawMessage(`[{}]`)}, "object"},
		{"empty", paramProperty{}, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveType(tt.prop); got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestExtractFlags(t *testing.T) {
	tests := []struct {
		name      string
		params    *nodeParams
		wantNil   bool
		wantNames []string
		wantDesc  string
		wantReq   bool
		checkReq  bool
	}{
		{
			name:    "nil params returns nil",
			params:  nil,
			wantNil: true,
		},
		{
			name:      "underscore converted to dash",
			params:    &nodeParams{Properties: map[string]paramProperty{"my_flag": {Type: "string", Description: "desc"}}},
			wantNames: []string{"my-flag"},
		},
		{
			name:      "leading underscore becomes control. prefix",
			params:    &nodeParams{Properties: map[string]paramProperty{"_limit": {Type: "integer", Description: "limit"}}},
			wantNames: []string{"control.limit"},
		},
		{
			name:      "double leading underscore becomes control. prefix",
			params:    &nodeParams{Properties: map[string]paramProperty{"__offset": {Type: "integer", Description: "offset"}}},
			wantNames: []string{"control.offset"},
		},
		{
			name: "required flag is marked required",
			params: &nodeParams{
				Properties: map[string]paramProperty{"name": {Type: "string", Description: "the name"}},
				Required:   []string{"name"},
			},
			wantNames: []string{"name"},
			wantReq:   true,
			checkReq:  true,
		},
		{
			name:      "non-required flag is not marked required",
			params:    &nodeParams{Properties: map[string]paramProperty{"name": {Type: "string", Description: "the name"}}},
			wantNames: []string{"name"},
			wantReq:   false,
			checkReq:  true,
		},
		{
			name:      "description falls back to title",
			params:    &nodeParams{Properties: map[string]paramProperty{"flag": {Type: "string", Title: "My Title"}}},
			wantNames: []string{"flag"},
			wantDesc:  "My Title",
		},
		{
			name:      "description takes priority over title",
			params:    &nodeParams{Properties: map[string]paramProperty{"flag": {Type: "string", Description: "My Desc", Title: "My Title"}}},
			wantNames: []string{"flag"},
			wantDesc:  "My Desc",
		},
		{
			name: "flags sorted alphabetically",
			params: &nodeParams{Properties: map[string]paramProperty{
				"zebra": {Type: "string"},
				"apple": {Type: "string"},
				"mango": {Type: "string"},
			}},
			wantNames: []string{"apple", "mango", "zebra"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := extractFlags(tt.params)
			if tt.wantNil {
				if flags != nil {
					t.Error("expected nil")
				}
				return
			}
			if len(flags) != len(tt.wantNames) {
				t.Fatalf("got %d flags, want %d", len(flags), len(tt.wantNames))
			}
			for i, name := range tt.wantNames {
				if flags[i].Name != name {
					t.Errorf("flags[%d].Name = %q, want %q", i, flags[i].Name, name)
				}
			}
			if tt.wantDesc != "" && flags[0].Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", flags[0].Description, tt.wantDesc)
			}
			if tt.checkReq && flags[0].Required != tt.wantReq {
				t.Errorf("Required = %v, want %v", flags[0].Required, tt.wantReq)
			}
		})
	}
}

func TestGetCommands(t *testing.T) {
	flagProp := map[string]paramProperty{"id": {Type: "string", Description: "id"}}

	tests := []struct {
		name         string
		nodes        []treeNode
		path         []string
		wantCount    int
		wantCommands map[string]string
	}{
		{
			name: "internal nodes are skipped",
			nodes: []treeNode{
				{Name: "internal", IsInternal: true, Parameters: &nodeParams{Properties: flagProp}},
			},
			wantCount: 0,
		},
		{
			name: "command with flags has mgc prefix and [flags] suffix",
			nodes: []treeNode{
				{Name: "create", Parameters: &nodeParams{Properties: flagProp}},
			},
			path:      []string{"vm", "instances"},
			wantCount: 1,
			wantCommands: map[string]string{
				"vm.instances.create": "mgc vm instances create [flags]",
			},
		},
		{
			name: "command without flags has no [flags] suffix",
			nodes: []treeNode{
				{Name: "list", Parameters: &nodeParams{Properties: map[string]paramProperty{}}},
			},
			path:      []string{"vm"},
			wantCount: 1,
			wantCommands: map[string]string{
				"vm.list": "mgc vm list",
			},
		},
		{
			name: "nested path builds correct key",
			nodes: []treeNode{
				{Name: "instances", Children: []treeNode{
					{Name: "create", Parameters: &nodeParams{Properties: flagProp}},
				}},
			},
			path:      []string{"vm"},
			wantCount: 1,
			wantCommands: map[string]string{
				"vm.instances.create": "mgc vm instances create [flags]",
			},
		},
		{
			name:      "node without parameters is not added to result",
			nodes:     []treeNode{{Name: "group"}},
			wantCount: 0,
		},
		{
			name: "sibling nodes do not corrupt each other's path",
			nodes: []treeNode{
				{Name: "create", Parameters: &nodeParams{Properties: flagProp}},
				{Name: "delete", Parameters: &nodeParams{Properties: flagProp}},
				{Name: "update", Parameters: &nodeParams{Properties: flagProp}},
			},
			path:      []string{"vm"},
			wantCount: 3,
			wantCommands: map[string]string{
				"vm.create": "mgc vm create [flags]",
				"vm.delete": "mgc vm delete [flags]",
				"vm.update": "mgc vm update [flags]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]commandResult)
			getCommands(tt.nodes, tt.path, result)

			if len(result) != tt.wantCount {
				t.Fatalf("got %d entries, want %d", len(result), tt.wantCount)
			}
			for key, wantCmd := range tt.wantCommands {
				entry, ok := result[key]
				if !ok {
					t.Errorf("missing key %q", key)
					continue
				}
				if entry.Command != wantCmd {
					t.Errorf("result[%q].Command = %q, want %q", key, entry.Command, wantCmd)
				}
			}
		})
	}
}

func TestLoadExistingOutput(t *testing.T) {
	validContent, _ := json.Marshal(map[string]commandResult{
		"vm.create": {Flags: []flagResult{
			{Name: "name", Description: "Nome", OriginalDescription: "Name"},
			{Name: "region", Description: "Região", OriginalDescription: "Region"},
		}},
	})

	tests := []struct {
		name         string
		content      []byte
		noFile       bool
		wantNil      bool
		wantKey      string
		wantFlagDesc string
		wantFlagOrig string
	}{
		{
			name:    "file does not exist",
			noFile:  true,
			wantNil: true,
		},
		{
			name:    "invalid JSON",
			content: []byte("not json"),
			wantNil: true,
		},
		{
			name:         "returns full command result map",
			content:      validContent,
			wantKey:      "vm.create",
			wantFlagDesc: "Nome",
			wantFlagOrig: "Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.noFile {
				path = "/nonexistent/path.json"
			} else {
				path = writeTempFile(t, tt.content)
			}

			result := loadExistingOutput(path)

			if tt.wantNil {
				if result != nil {
					t.Error("expected nil")
				}
				return
			}
			entry, ok := result[tt.wantKey]
			if !ok {
				t.Fatalf("missing key %q", tt.wantKey)
			}
			flags := flagsByName(entry.Flags)
			if flags["name"].Description != tt.wantFlagDesc {
				t.Errorf("Description = %q, want %q", flags["name"].Description, tt.wantFlagDesc)
			}
			if flags["name"].OriginalDescription != tt.wantFlagOrig {
				t.Errorf("OriginalDescription = %q, want %q", flags["name"].OriginalDescription, tt.wantFlagOrig)
			}
		})
	}
}

func TestTranslateAll(t *testing.T) {
	tests := []struct {
		name            string
		itemCount       int
		serverError     bool
		wantCallCount   int
		wantResultCount int
		wantErr         bool
	}{
		{
			name:            "splits into chunks of deeplChunkSize",
			itemCount:       deeplChunkSize + 1,
			wantCallCount:   2,
			wantResultCount: deeplChunkSize + 1,
		},
		{
			name:            "exact chunk size needs one call",
			itemCount:       deeplChunkSize,
			wantCallCount:   1,
			wantResultCount: deeplChunkSize,
		},
		{
			name:        "propagates chunk error",
			itemCount:   1,
			serverError: true,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				if tt.serverError {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				var req deeplRequest
				json.NewDecoder(r.Body).Decode(&req)
				resp := deeplResponse{}
				for _, text := range req.Text {
					resp.Translations = append(resp.Translations, struct {
						Text string `json:"text"`
					}{Text: "t_" + text})
				}
				json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			texts := make([]string, tt.itemCount)
			for i := range texts {
				texts[i] = fmt.Sprintf("text_%d", i)
			}

			tr := &deeplTranslator{apiKey: "key", endpoint: server.URL}
			result, err := tr.translateAll(texts)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if callCount != tt.wantCallCount {
				t.Errorf("API call count = %d, want %d", callCount, tt.wantCallCount)
			}
			if len(result) != tt.wantResultCount {
				t.Errorf("got %d results, want %d", len(result), tt.wantResultCount)
			}
		})
	}
}

func TestApplyTranslations(t *testing.T) {
	tests := []struct {
		name             string
		existingFlags    []flagResult
		inputDesc        string
		wantDesc         string
		wantOriginalDesc string
	}{
		{
			name: "unchanged description keeps existing translation",
			existingFlags: []flagResult{
				{Name: "name", Description: "Nome", OriginalDescription: "Name"},
			},
			inputDesc:        "Name",
			wantDesc:         "Nome",
			wantOriginalDesc: "Name",
		},
		{
			name: "changed description triggers retranslation (API fail fallback)",
			existingFlags: []flagResult{
				{Name: "name", Description: "Nome", OriginalDescription: "Name"},
			},
			inputDesc:        "New Name",
			wantDesc:         "New Name",
			wantOriginalDesc: "New Name",
		},
		{
			name:             "new flag with no existing entry marks original on API failure",
			existingFlags:    nil,
			inputDesc:        "Create a VM",
			wantDesc:         "Create a VM",
			wantOriginalDesc: "Create a VM",
		},
		{
			name: "empty description is skipped entirely",
			existingFlags: []flagResult{
				{Name: "name", Description: "Nome", OriginalDescription: "Name"},
			},
			inputDesc:        "",
			wantDesc:         "",
			wantOriginalDesc: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outputFile string
			if len(tt.existingFlags) > 0 {
				existing := map[string]commandResult{
					"vm.create": {Flags: tt.existingFlags},
				}
				outputFile = writeTempFile(t, mustMarshal(t, existing))
			} else {
				outputFile = "/nonexistent/output.json"
			}

			result := map[string]commandResult{
				"vm.create": {Flags: []flagResult{{Name: "name", Description: tt.inputDesc}}},
			}

			_ = applyTranslations(result, "fake-key", outputFile)

			flag := result["vm.create"].Flags[0]
			if flag.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", flag.Description, tt.wantDesc)
			}
			if flag.OriginalDescription != tt.wantOriginalDesc {
				t.Errorf("OriginalDescription = %q, want %q", flag.OriginalDescription, tt.wantOriginalDesc)
			}
		})
	}
}

func TestDeeplAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		envKey  string
		flagKey string
		want    string
	}{
		{"flag takes priority over env var", "from-env", "from-flag", "from-flag"},
		{"env var used when flag is empty", "from-env", "", "from-env"},
		{"returns empty when nothing is set", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DEEPL_API_KEY", tt.envKey)
			if got := deeplAPIKey(genCommandsJSONParams{deeplKey: tt.flagKey}); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunGenCommandsJSON(t *testing.T) {
	validTree, _ := json.Marshal([]treeNode{
		{
			Name: "vm",
			Children: []treeNode{
				{
					Name: "create",
					Parameters: &nodeParams{
						Properties: map[string]paramProperty{
							"name":   {Type: "string", Description: "Instance name"},
							"region": {Type: "string", Description: "Region"},
						},
						Required: []string{"name"},
					},
				},
				{Name: "internal-op", IsInternal: true, Parameters: &nodeParams{
					Properties: map[string]paramProperty{"x": {Type: "string"}},
				}},
			},
		},
	})

	tests := []struct {
		name        string
		dumpContent []byte
		wantErr     bool
		wantKey     string
		wantCommand string
		wantFlags   int
		wantMissing string
	}{
		{
			name:        "generates correct output",
			dumpContent: validTree,
			wantKey:     "vm.create",
			wantCommand: "mgc vm create [flags]",
			wantFlags:   2,
			wantMissing: "vm.internal-op",
		},
		{
			name:    "returns error for missing dump file",
			wantErr: true,
		},
		{
			name:        "returns error for invalid dump JSON",
			dumpContent: []byte("not json"),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			outputFile := filepath.Join(dir, "commands.json")

			var dumpFile string
			if tt.dumpContent != nil {
				dumpFile = filepath.Join(dir, "dump.json")
				if err := os.WriteFile(dumpFile, tt.dumpContent, 0644); err != nil {
					t.Fatal(err)
				}
			} else {
				dumpFile = "/nonexistent/dump.json"
			}

			err := runGenCommandsJSON(genCommandsJSONParams{
				dumpFile:   dumpFile,
				outputFile: outputFile,
			})

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var result map[string]commandResult
			if err := json.Unmarshal(mustReadFile(t, outputFile), &result); err != nil {
				t.Fatalf("invalid JSON output: %v", err)
			}
			if tt.wantKey != "" {
				entry, ok := result[tt.wantKey]
				if !ok {
					t.Fatalf("missing key %q", tt.wantKey)
				}
				if entry.Command != tt.wantCommand {
					t.Errorf("Command = %q, want %q", entry.Command, tt.wantCommand)
				}
				if len(entry.Flags) != tt.wantFlags {
					t.Errorf("Flags count = %d, want %d", len(entry.Flags), tt.wantFlags)
				}
			}
			if tt.wantMissing != "" {
				if _, ok := result[tt.wantMissing]; ok {
					t.Errorf("key %q should not be present", tt.wantMissing)
				}
			}
		})
	}
}

func writeTempFile(t *testing.T, data []byte) string {
	t.Helper()
	f, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
