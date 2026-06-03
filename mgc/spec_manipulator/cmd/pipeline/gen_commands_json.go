package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type treeNode struct {
	Name       string      `json:"name"`
	IsInternal bool        `json:"isInternal"`
	Children   []treeNode  `json:"children"`
	Parameters *nodeParams `json:"parameters"`
}

type nodeParams struct {
	Properties map[string]paramProperty `json:"properties"`
	Required   []string                 `json:"required"`
}

type paramProperty struct {
	Type        string                   `json:"type"`
	Description string                   `json:"description"`
	Title       string                   `json:"title"`
	Default     any                      `json:"default"`
	Items       *propertyItems           `json:"items"`
	AnyOf       json.RawMessage          `json:"anyOf"`
	OneOf       json.RawMessage          `json:"oneOf"`
	Properties  map[string]paramProperty `json:"properties"`
}

type propertyItems struct {
	Type string `json:"type"`
}

type commandResult struct {
	Command string       `json:"command"`
	Flags   []flagResult `json:"flags"`
}

type flagResult struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	OriginalDescription string `json:"original_description"`
	Type                string `json:"type"`
	Required            bool   `json:"required"`
	DefaultValue        any    `json:"default_value"`
}

const deeplChunkSize = 50

type deeplTranslator struct {
	apiKey   string
	endpoint string
}

type deeplRequest struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

type deeplResponse struct {
	Translations []struct {
		Text string `json:"text"`
	} `json:"translations"`
}

type genCommandsJSONParams struct {
	dumpFile   string
	outputFile string
	translate  bool
	deeplKey   string
}

func newDeeplTranslator(apiKey string) *deeplTranslator {
	endpoint := "https://api.deepl.com/v2/translate"
	if strings.HasSuffix(apiKey, ":fx") {
		endpoint = "https://api-free.deepl.com/v2/translate"
	}
	return &deeplTranslator{apiKey: apiKey, endpoint: endpoint}
}

func (t *deeplTranslator) translateChunk(texts []string) ([]string, error) {
	body, err := json.Marshal(deeplRequest{Text: texts, TargetLang: "PT-BR"})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, t.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "DeepL-Auth-Key "+t.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deepl API returned status %d: %s", resp.StatusCode, bytes.TrimSpace(body))
	}

	var result deeplResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	translated := make([]string, len(result.Translations))
	for i, tr := range result.Translations {
		translated[i] = tr.Text
	}
	return translated, nil
}

func (t *deeplTranslator) translateAll(texts []string) ([]string, error) {
	all := make([]string, 0, len(texts))
	for i := 0; i < len(texts); i += deeplChunkSize {
		chunk, err := t.translateChunk(texts[i:min(i+deeplChunkSize, len(texts))])
		if err != nil {
			return nil, fmt.Errorf("translating chunk starting at %d: %w", i, err)
		}
		all = append(all, chunk...)
	}
	return all, nil
}

func loadExistingOutput(outputFile string) map[string]commandResult {
	data, err := os.ReadFile(outputFile)
	if err != nil {
		return nil
	}
	var existing map[string]commandResult
	if err := json.Unmarshal(data, &existing); err != nil {
		return nil
	}
	return existing
}

func flagsByName(flags []flagResult) map[string]flagResult {
	m := make(map[string]flagResult, len(flags))
	for _, f := range flags {
		m[f.Name] = f
	}
	return m
}

func applyTranslations(result map[string]commandResult, apiKey, outputFile string) error {
	existing := loadExistingOutput(outputFile)

	seen := make(map[string]struct{})
	var toTranslate []string

	for key, entry := range result {
		existingFlags := flagsByName(existing[key].Flags)
		for _, flag := range entry.Flags {
			if flag.Description == "" {
				continue
			}
			if ex, ok := existingFlags[flag.Name]; ok && ex.OriginalDescription == flag.Description {
				continue
			}
			if _, ok := seen[flag.Description]; !ok {
				seen[flag.Description] = struct{}{}
				toTranslate = append(toTranslate, flag.Description)
			}
		}
	}

	translationMap := make(map[string]string)
	switch {
	case len(toTranslate) == 0:
		fmt.Println("All descriptions already translated, skipping API call.")
	case apiKey == "":
		fmt.Fprintf(os.Stderr, "warning: %d description(s) need translation but no API key available, skipping\n", len(toTranslate))
	default:
		fmt.Printf("Translating %d new/changed descriptions via DeepL...\n", len(toTranslate))
		translated, err := newDeeplTranslator(apiKey).translateAll(toTranslate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: translation failed (%v), keeping original descriptions\n", err)
		} else {
			for i, original := range toTranslate {
				translationMap[original] = translated[i]
			}
		}
	}

	for key, entry := range result {
		existingFlags := flagsByName(existing[key].Flags)
		for i, flag := range entry.Flags {
			if flag.Description == "" {
				if ex, ok := existingFlags[flag.Name]; ok && ex.Description != "" {
					entry.Flags[i].Description = ex.Description
					entry.Flags[i].OriginalDescription = ex.OriginalDescription
				}
				continue
			}
			ex, unchanged := existingFlags[flag.Name]
			if unchanged && ex.OriginalDescription == flag.Description {
				entry.Flags[i].OriginalDescription = ex.OriginalDescription
				entry.Flags[i].Description = ex.Description
			} else if pt, ok := translationMap[flag.Description]; ok {
				entry.Flags[i].OriginalDescription = flag.Description
				entry.Flags[i].Description = pt
			} else {
				entry.Flags[i].OriginalDescription = flag.Description
			}
		}
		result[key] = entry
	}

	return nil
}

func collectObjectSubProperties(prop paramProperty) map[string]paramProperty {
	if len(prop.Properties) > 0 {
		return prop.Properties
	}
	merged := make(map[string]paramProperty)
	for _, raw := range []json.RawMessage{prop.AnyOf, prop.OneOf} {
		if raw == nil {
			continue
		}
		var schemas []paramProperty
		if err := json.Unmarshal(raw, &schemas); err != nil {
			continue
		}
		for _, schema := range schemas {
			if schema.Type != "object" {
				continue
			}
			for k, v := range schema.Properties {
				merged[k] = v
			}
		}
	}
	return merged
}

func expandObjectFlags(prefix string, prop paramProperty) []flagResult {
	var flags []flagResult
	for subName, subProp := range collectObjectSubProperties(prop) {
		subFlagName := prefix + "." + strings.ReplaceAll(subName, "_", "-")
		subDesc := subProp.Description
		if subDesc == "" {
			subDesc = subProp.Title
		}
		subType := resolveType(subProp)
		flags = append(flags, flagResult{
			Name:        subFlagName,
			Description: subDesc,
			Type:        subType,
			Required:    false,
		})
		if subType == "object" {
			flags = append(flags, expandObjectFlags(subFlagName, subProp)...)
		}
	}
	return flags
}

func resolveType(prop paramProperty) string {
	if prop.Type == "" {
		if prop.AnyOf != nil || prop.OneOf != nil {
			return "object"
		}
		return "unknown"
	}
	if prop.Type == "array" && prop.Items != nil && prop.Items.Type != "" {
		return fmt.Sprintf("array(%s)", prop.Items.Type)
	}
	return prop.Type
}

func extractFlags(params *nodeParams) []flagResult {
	if params == nil {
		return nil
	}

	required := make(map[string]bool, len(params.Required))
	for _, r := range params.Required {
		required[r] = true
	}

	flags := make([]flagResult, 0, len(params.Properties))
	for name, prop := range params.Properties {
		flagName := strings.ReplaceAll(name, "_", "-")
		if strings.HasPrefix(flagName, "-") {
			flagName = "control." + strings.TrimLeft(flagName, "-")
		}

		desc := prop.Description
		if desc == "" {
			desc = prop.Title
		}

		flagType := resolveType(prop)

		flags = append(flags, flagResult{
			Name:         flagName,
			Description:  desc,
			Type:         flagType,
			Required:     required[name],
			DefaultValue: prop.Default,
		})

		if flagType == "object" {
			flags = append(flags, expandObjectFlags(flagName, prop)...)
		}
	}

	sort.Slice(flags, func(i, j int) bool { return flags[i].Name < flags[j].Name })
	return flags
}

func getCommands(nodes []treeNode, path []string, result map[string]commandResult) {
	for _, node := range nodes {
		if node.IsInternal {
			continue
		}

		currentPath := append(path[:len(path):len(path)], node.Name)
		if node.Parameters != nil {
			flags := extractFlags(node.Parameters)
			command := "mgc " + strings.Join(currentPath, " ")
			if len(flags) > 0 {
				command += " [flags]"
			}
			result[strings.Join(currentPath, ".")] = commandResult{
				Command: command,
				Flags:   flags,
			}
		}
		getCommands(node.Children, currentPath, result)
	}
}

func loadAPIKeyFromDotEnv() string {
	env, err := godotenv.Read(".env")
	if err != nil {
		return ""
	}
	return env["DEEPL_API_KEY"]
}

func deeplAPIKey(params genCommandsJSONParams) string {
	apiKey := params.deeplKey
	if apiKey == "" {
		apiKey = os.Getenv("DEEPL_API_KEY")
	}
	if apiKey == "" {
		apiKey = loadAPIKeyFromDotEnv()
	}

	return apiKey
}

func runGenCommandsJSON(params genCommandsJSONParams) error {
	data, err := os.ReadFile(params.dumpFile)
	if err != nil {
		return fmt.Errorf("reading dump file %q: %w", params.dumpFile, err)
	}

	var tree []treeNode
	if err := json.Unmarshal(data, &tree); err != nil {
		return fmt.Errorf("parsing dump file: %w", err)
	}

	result := make(map[string]commandResult)
	getCommands(tree, nil, result)

	if params.translate {
		apiKey := deeplAPIKey(params)

		if err := applyTranslations(result, apiKey, params.outputFile); err != nil {
			return fmt.Errorf("translating descriptions: %w", err)
		}
	}

	tmp, err := os.CreateTemp(filepath.Dir(params.outputFile), "commands-*.json.tmp")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := func() error {
		defer tmp.Close()
		enc := json.NewEncoder(tmp)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		return enc.Encode(result)
	}(); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	if err := os.Rename(tmpPath, params.outputFile); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("Generated %d commands -> %s\n", len(result), params.outputFile)
	return nil
}

func GenCommandsJSONCmd() *cobra.Command {
	params := &genCommandsJSONParams{}
	cmd := &cobra.Command{
		Use:   "gen-commands-json",
		Short: "Generates a JSON with all CLI commands and their flags from a dump-tree file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenCommandsJSON(*params)
		},
	}
	cmd.Flags().StringVarP(&params.dumpFile, "dump", "d", "../cli/cli-dump-tree.json", "Path to the CLI dump tree JSON file")
	cmd.Flags().StringVarP(&params.outputFile, "output", "o", "../cli/commands.json", "Path for the output JSON file")
	cmd.Flags().BoolVarP(&params.translate, "translate", "t", false, "Translate flag descriptions to PT-BR via DeepL API")
	cmd.Flags().StringVar(&params.deeplKey, "deepl-key", "", "DeepL API key (fallback: DEEPL_API_KEY env var or .env file)")

	return cmd
}
