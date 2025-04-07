package schema_flags

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

const (
	ValueLoadJSONFromFilePrefix     = "@"
	ValueLoadVerbatimFromFilePrefix = "%"
	ValueVerbatimStringPrefix       = "#"
	ValueHelpIsRequired             = "help"
)

// handles special cases, in order:
//  1. "": empty value is returned. No error.
//  2. "@filename": load JSON from file.
//     Returns error if file was not found or it's not a valid JSON for type "T".
//  3. try to JSON parse as value type "T"
func parseJSONFlagValue[T any](rawValue string) (value T, err error) {
	switch {
	case rawValue == "":
		return

	case strings.HasPrefix(rawValue, ValueLoadJSONFromFilePrefix):
		return loadJSONFromFile[T](rawValue[1:])

	default:
		err = json.Unmarshal([]byte(rawValue), &value)
		return
	}
}

// handles special cases targeted at strings, in addition to parseJSON(), in order:
//  1. "%filename": load raw string from filename (no trim or parsing is done).
//     Returns error if file was not found.
//  2. "#string": use the rest of the string verbatim (no trim or parsing is done). No error.
//  3. "help": show the flag help. To provide the value "help", use "#help" or provide it quoted.
//  4. parseJSON(), if errors then use the rawValue instead. No error
//
// To pass a string with leading-and-trailing quotes (`"value"`)
// one must either provide a version with escaped quotes (`"\"value\""`)
// or use the verbatim prefix (`#"value"`)
func parseStringFlagValue(rawValue string) (value string, err error) {
	switch {
	case strings.HasPrefix(rawValue, ValueLoadVerbatimFromFilePrefix):
		return loadVerbatimFromFile(rawValue[1:])

	case strings.HasPrefix(rawValue, ValueVerbatimStringPrefix):
		return rawValue[1:], nil

	case rawValue == ValueHelpIsRequired:
		return "", ErrWantHelp

	default:
		value, err = parseJSONFlagValue[string](rawValue)
		if err != nil {
			value = rawValue
			err = nil
		}
		return
	}
}

func parseBoolFlagValue(rawValue string) (value bool, err error) {
	switch rawValue {
	case "":
		return // default to false

	case ValueHelpIsRequired:
		return false, ErrWantHelp

	default:
		return strconv.ParseBool(rawValue) // mimics pflag's boolValue
	}
}

func isWhiteSpace(c rune) bool {
	return unicode.IsSpace(c)
}

func isCSVDelimiter(c rune) bool {
	return c == ',' || c == ';'
}

func isCSVDelimiterOrWhiteSpace(c rune) bool {
	return isCSVDelimiter(c) || isWhiteSpace(c)
}

func isObjectKeyDelimiter(c rune) bool {
	return c == ':' || c == '='
}

func skipLeadingMatches(s string, isMatch func(rune) bool) (end int) {
	for i, c := range s {
		if !isMatch(c) {
			return i
		}
	}

	return len(s)
}

func skipWhiteSpaces(s string) (end int) {
	return skipLeadingMatches(s, isWhiteSpace)
}

func skipCSVDelimitersOrWhiteSpaces(s string) (end int) {
	return skipLeadingMatches(s, isCSVDelimiterOrWhiteSpace)
}

func parseNonWhitespaceCSVItem(s string) (item string, end int) {
	for i, c := range s {
		if isCSVDelimiterOrWhiteSpace(c) {
			end = i + 1
			return s[:i], end
		}
	}

	return s, len(s)
}

func parseNumberItem(s string) (value any, end int, err error) {
	text, end := parseNonWhitespaceCSVItem(s)
	value, err = strconv.ParseFloat(text, 64)
	return
}

func parseIntegerItem(s string) (value any, end int, err error) {
	text, end := parseNonWhitespaceCSVItem(s)
	i, err := strconv.ParseInt(text, 10, 64)
	if math.MinInt <= i && i <= math.MaxInt {
		value = int(i)
	} else {
		value = i
	}
	return
}

func parseBooleanItem(s string) (value any, end int, err error) {
	text, end := parseNonWhitespaceCSVItem(s)
	value, err = strconv.ParseBool(text)
	return
}

func parseQuotedStringItem(s string) (value string, end int, err error) {
	var quote rune
	escaped := false

	for i, c := range s {
		if i == 0 {
			quote = c
			continue
		}

		if escaped {
			escaped = false
		} else if c == '\\' {
			escaped = true
			continue
		} else if c == quote {
			end = i + 1
			end += skipCSVDelimitersOrWhiteSpaces(s[end:])
			return
		}

		value += string(c)
	}

	return value, len(s), fmt.Errorf("missing end quote char %c: %s", quote, s)
}

func parseString(s string, isDelimiter func(rune) bool) (value string, end int, err error) {
	start := skipWhiteSpaces(s)
	remaining := s[start:]
	if len(remaining) > 1 && (remaining[0] == '"' || remaining[0] == '\'') {
		value, end, err = parseQuotedStringItem(remaining)
		end += start
		if err == nil && end < len(s) && isDelimiter(rune(s[end])) {
			end++
		}
		return
	}

	if remaining == "" {
		end = start
		return
	}

	lastNonWhitespace := 0
	for i, c := range remaining {
		if isDelimiter(c) {
			end = start + i + 1
			value = remaining[:lastNonWhitespace+1]
			return
		} else if !isWhiteSpace(c) {
			lastNonWhitespace = i
		}
	}

	return remaining[:lastNonWhitespace+1], len(s), nil
}

func parseStringItem(s string) (value any, end int, err error) {
	return parseString(s, isCSVDelimiter)
}

func parseObjectItem(schema *core.Schema, s string) (value any, end int, err error) {
	value, end, err = parseObjectItemCSV(schema, s, nil)
	if err == nil {
		return
	}

	return parseAnyItem(s)
}

func parseAnyItem(s string) (value any, end int, err error) {
	start := skipWhiteSpaces(s)
	end = start
	if start >= len(s) {
		return
	}

	if err = json.Unmarshal([]byte(s[start:]), &value); err == nil {
		end = len(s)
		return
	}

	var syntaxError = new(json.SyntaxError)
	if !errors.As(err, &syntaxError) {
		return
	}

	errorAt := int(syntaxError.Offset - 1)
	if errorAt <= 0 {
		return parseStringItem(s)
	}

	var c rune
	if n, _ := fmt.Sscanf(syntaxError.Error(), "invalid character '%c' after top-level value", &c); n != 1 || !isCSVDelimiter(c) {
		return
	}

	err = json.Unmarshal([]byte(s[start:errorAt]), &value)
	end = start + errorAt + 1
	return
}

func getItemTypeParser(schema *core.Schema) (itemParser func(s string) (any, int, error)) {
	var itemType string
	if schema != nil && len(schema.Type.Slice()) > 0 {
		itemType = schema.Type.Slice()[0]
	}

	switch itemType {
	case "number":
		return parseNumberItem
	case "integer":
		return parseIntegerItem
	case "boolean":
		return parseBooleanItem
	case "string":
		return parseStringItem
	case "object":
		return func(s string) (any, int, error) {
			return parseObjectItem(schema, s)
		}
	default:
		return parseAnyItem
	}

}

func parseArrayCSV(itemsSchema *core.Schema, rawValue string) (value []any, err error) {
	itemParser := getItemTypeParser(itemsSchema)

	for {
		end := skipCSVDelimitersOrWhiteSpaces(rawValue)
		if end >= len(rawValue) {
			return
		}
		rawValue = rawValue[end:]

		var item any
		item, end, err = itemParser(rawValue)
		if err != nil {
			return
		}
		rawValue = rawValue[end:]

		value = append(value, item)
	}
}

func parseArrayFlagValueSingle(itemsSchema *core.Schema, rawValue string) (value []any, err error) {
	value, err = parseJSONFlagValue[[]any](rawValue)
	if err == nil {
		return
	}

	value, csvErr := parseArrayCSV(itemsSchema, rawValue)
	if csvErr == nil {
		err = nil
		return
	}

	return
}

func parseArrayFlagValue(itemsSchema *core.Schema, rawValues []string) (items []any, err error) {
	if len(rawValues) == 1 && rawValues[0] == ValueHelpIsRequired {
		return nil, ErrWantHelp
	}

	for i, rawValue := range rawValues {
		value, err := parseArrayFlagValueSingle(itemsSchema, rawValue)
		if err != nil {
			if len(rawValues) > 0 {
				err = &utils.ChainedError{Name: fmt.Sprint(i), Err: err}
			}
			return items, err
		}
		items = append(items, value...)
	}

	return
}

func parseObjectKey(s string) (key string, end int, err error) {
	return parseString(s, isObjectKeyDelimiter)
}

func parseObjectValue(propSchema *core.Schema, s string) (value any, end int, err error) {
	start := skipWhiteSpaces(s)
	if start >= len(s) {
		return
	}

	itemParser := getItemTypeParser(propSchema)
	value, end, err = itemParser(s[start:])
	end += start
	if err != nil {
		return
	}

	return
}

const objectKeyPathDelimiter = '.'

func isObjectKeyPathDelimiter(r rune) bool {
	return r == objectKeyPathDelimiter
}

func parseObjectPropertyValue(schema *core.Schema, propName string, s string) (propValue any, end int, err error) {
	for _, os := range mgcSchemaPkg.CollectObjectPropertySchemas(schema, propName) {
		propValue, end, err = parseObjectValue(os.PropSchema, s)
		if err == nil {
			return
		}
	}
	return nil, 0, fmt.Errorf("could not find property %q", propName)
}

func parseObjectValueFromNextKeyPath(schema *core.Schema, propName, keyPath, s string) (nextPropName string, propValue any, end int, err error) {
	for _, os := range mgcSchemaPkg.CollectObjectPropertySchemas(schema, propName) {
		nextPropName, propValue, end, err = parseObjectValueFromPath(os.PropSchema, keyPath, s)
		if err == nil {
			return
		}
	}
	return "", nil, 0, fmt.Errorf("could not find property %q", propName)
}

func parseObjectValueFromPath(schema *core.Schema, keyPath string, s string) (propName string, propValue any, end int, err error) {
	propName, nextKey, err := parseString(keyPath, isObjectKeyPathDelimiter)
	if err != nil {
		return
	}
	nextKey += skipWhiteSpaces(keyPath[nextKey:])
	keyPath = keyPath[nextKey:]

	if propName == string(objectKeyPathDelimiter) {
		return parseObjectValueFromPath(schema, keyPath, s)
	}

	if keyPath == "" {
		propValue, end, err = parseObjectPropertyValue(schema, propName, s)
		return
	}

	childName, childValue, end, err := parseObjectValueFromNextKeyPath(schema, propName, keyPath, s)
	if err != nil {
		return
	}
	propValue = map[string]any{childName: childValue}

	return
}

// if both 'a' and 'b' are maps, merge their items recursively.
// if only one of them is a map, then fail.
//
// other kind of types, 'a' is returned
func mergeValue(a, b any) (r any, err error) {
	r = a
	mA, ok := a.(map[string]any)
	if !ok {
		if _, ok := b.(map[string]any); ok {
			err = fmt.Errorf("cannot merge types %T and %T (a: %#v, b: %#v)", a, b, a, b)
			return
		}
		return
	}

	mB, ok := b.(map[string]any)
	if !ok {
		err = fmt.Errorf("cannot merge types %T and %T (a: %#v, b: %#v)", a, b, a, b)
		return
	}

	m := make(map[string]any, len(mA)+len(mB))
	for k, v := range mA {
		m[k] = v
	}
	for k, v := range mB {
		if existing, hasExisting := m[k]; hasExisting {
			v, err = mergeValue(v, existing)
			if err != nil {
				return
			}
		}
		m[k] = v
	}
	r = m

	return
}

func parseObjectItemCSV(schema *core.Schema, rawValue string, v map[string]any) (value map[string]any, end int, err error) {
	value = v
	start := skipCSVDelimitersOrWhiteSpaces(rawValue)
	end += start
	if start >= len(rawValue) {
		return
	}
	rawValue = rawValue[start:]
	if rawValue[0] == '{' {
		parsedValue, start, err := parseAnyItem(rawValue)
		end += start
		if err != nil {
			return v, end, err
		}
		m, ok := parsedValue.(map[string]any)
		if !ok {
			err = fmt.Errorf("expected JSON object at %q", rawValue)
			return v, end, err
		}
		mergedValue, err := mergeValue(m, v)
		if err != nil {
			return v, end, err
		}
		return mergedValue.(map[string]any), end, nil
	}

	var propName string
	propName, start, err = parseObjectKey(rawValue)
	end += start
	if err != nil {
		return
	}
	rawValue = rawValue[start:]
	if rawValue == "" {
		err = fmt.Errorf("missing property %q value", propName)
		return
	}

	propValue, start, err := parseObjectPropertyValue(schema, propName, rawValue)
	if err == nil {
		end += start
	} else {
		if !strings.ContainsFunc(propName, isObjectKeyPathDelimiter) {
			return
		}
		propName, propValue, start, err = parseObjectValueFromPath(schema, propName, rawValue)
		if err != nil {
			return
		}
		end += start
	}

	if existing, hasExisting := value[propName]; hasExisting {
		propValue, err = mergeValue(propValue, existing)
		if err != nil {
			return
		}
	}

	if value == nil {
		value = map[string]any{}
	}
	value[propName] = propValue
	return
}

func parseObjectCSV(schema *core.Schema, rawValue string) (value map[string]any, err error) {
	for {
		var end int
		value, end, err = parseObjectItemCSV(schema, rawValue, value)
		if err != nil {
			return
		}

		rawValue = rawValue[end:]
		if rawValue == "" {
			return
		}
	}
}

func parseObjectFlagValueSingle(schema *core.Schema, rawValue string) (value map[string]any, err error) {
	value, err = parseJSONFlagValue[map[string]any](rawValue)
	if err == nil {
		return
	}

	value, csvErr := parseObjectCSV(schema, rawValue)
	if csvErr == nil {
		err = nil
		return
	}

	return
}

func parseObjectFlagValue(schema *core.Schema, rawValues []string) (items map[string]any, err error) {
	if len(rawValues) == 1 && rawValues[0] == ValueHelpIsRequired {
		return nil, ErrWantHelp
	}

	for i, rawValue := range rawValues {
		value, err := parseObjectFlagValueSingle(schema, rawValue)
		if err != nil {
			if len(rawValues) > 0 {
				err = &utils.ChainedError{Name: fmt.Sprint(i), Err: err}
			}
			return items, err
		}
		if len(value) == 0 {
			continue
		}

		if items == nil {
			items = map[string]any{}
		}

		maps.Copy(items, value)
	}

	return
}

func loadJSONFromFile[T any](filename string) (value T, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &value)
	return
}

func loadVerbatimFromFile(filename string) (value string, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	value = string(data)
	return
}
