package schema_flags

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"golang.org/x/exp/constraints"
)

func addXOfSchemaConstraints(message string, refs mgcSchemaPkg.SchemaRefs, dst *[]string) {
	constraints := make([]string, 0, len(refs))
	for _, ref := range refs {
		if ref == nil || ref.Value == nil || ref.Value.Type == nil {
			continue
		}
		if desc := getDescriptionConstraints((*mgcSchemaPkg.Schema)(ref.Value)); desc != "" {
			constraints = append(constraints, desc)
			continue
		}
		if t := getFlagType((*mgcSchemaPkg.Schema)(ref.Value)); t != "" {
			constraints = append(constraints, t)
			continue
		}
	}

	*dst = append(*dst, message+": "+formatAlternatives("%s", "%s or %s", constraints))
}

func addSchemaConstraints(s *mgcSchemaPkg.Schema, dst *[]string) {
	if len(s.OneOf) > 0 {
		addXOfSchemaConstraints("exactly one of", s.OneOf, dst)
		return
	}
	if len(s.AnyOf) > 0 {
		addXOfSchemaConstraints("at least one of", s.AnyOf, dst)
		return
	}

	if s.Type != nil {
		switch {
		case s.Type.Includes("string"):
			addStringConstraints(s, dst)
		case s.Type.Includes("integer"), s.Type.Includes("number"):
			addNumberConstraints(s, dst)
		case s.Type.Includes("array"):
			addArrayConstraints(s, dst)
		case s.Type.Includes("object"):
			addObjectConstraints(s, dst)
		}
	}
}

func addStringConstraints(s *mgcSdk.Schema, dst *[]string) {
	if s.MinLength != 0 && s.MaxLength != nil {
		*dst = append(*dst, fmt.Sprintf("between %v and %v characters", s.MinLength, int(*s.MaxLength)))
	} else if s.MinLength != 0 {
		*dst = append(*dst, fmt.Sprintf("min character count: %v", s.MinLength))
	} else if s.MaxLength != nil {
		*dst = append(*dst, fmt.Sprintf("max character count: %v", int(*s.MaxLength)))
	}

	if s.Pattern != "" {
		*dst = append(*dst, fmt.Sprintf("pattern: %s", s.Pattern))
	}
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func addNumberConstraints(s *mgcSdk.Schema, dst *[]string) {
	min := getMinNumber(s)
	max := getMaxNumber(s)

	if min != nil && max != nil {
		*dst = append(*dst, fmt.Sprintf("range: %s - %s", formatFloat(*min), formatFloat(*max)))
	} else if min != nil {
		*dst = append(*dst, fmt.Sprintf("min: %s", formatFloat(*min)))
	} else if max != nil {
		*dst = append(*dst, fmt.Sprintf("max: %s", formatFloat(*max)))
	}
}

func getPlural[T constraints.Integer](v T, singular, plural string) string {
	if v == 1 {
		return singular
	}
	return plural
}

func addCountConstraints[T constraints.Integer](minimum T, maximum *T, singular, plural string, dst *[]string) {
	if minimum != 0 && maximum != nil {
		*dst = append(*dst, fmt.Sprintf("between %v and %v %s", minimum, *maximum, getPlural(*maximum, singular, plural)))
	} else if minimum != 0 {
		*dst = append(*dst, fmt.Sprintf("at least %v %s", minimum, getPlural(minimum, singular, plural)))
	} else if maximum != nil {
		*dst = append(*dst, fmt.Sprintf("at most %v %s", *maximum, getPlural(*maximum, singular, plural)))
	}
}

func addArrayConstraints(s *mgcSdk.Schema, dst *[]string) {
	addCountConstraints(s.MinItems, s.MaxItems, "item", "items", dst)
}

func addObjectConstraints(s *mgcSdk.Schema, dst *[]string) {
	addCountConstraints(s.MinProps, s.MaxProps, "property", "properties", dst)

	if s.AdditionalProperties.Has != nil && *s.AdditionalProperties.Has {
		return
	}

	if len(s.Properties) == 0 {
		return
	}

	keys := []string{}
	for k := range s.Properties {
		keys = append(keys, k)
	}

	slices.Sort(keys)
	*dst = append(*dst, formatAlternatives("single property: %s", "properties: %s and %s", keys))
}

func getEnumAsString(s *mgcSdk.Schema) (asStrings []string) {
	length := len(s.Enum)
	if length == 0 {
		return
	}

	asStrings = make([]string, 0, length)
	for _, e := range s.Enum {
		data, err := json.Marshal(e)
		var s string
		if err != nil {
			s = fmt.Sprint(e)
		} else {
			s = string(data)
		}

		asStrings = append(asStrings, s)
	}

	slices.Sort(asStrings)

	return
}

func addEnumConstraint(s *mgcSdk.Schema, dst *[]string) {
	asStrings := getEnumAsString(s)
	if len(asStrings) == 0 {
		return
	}

	*dst = append(*dst, formatAlternatives("must be %s", "one of %s or %s", asStrings))
}

// oneFmt takes a single parameter, while multipleFmt takes exactly 2
func formatAlternatives(oneFmt string, multipleFmt string, asStrings []string) string {
	switch len(asStrings) {
	case 0:
		return ""

	case 1:
		return fmt.Sprintf(oneFmt, asStrings[0])

	default:
		lastIndex := len(asStrings) - 1
		commaDelimited := strings.Join(asStrings[:lastIndex], ", ")
		return fmt.Sprintf(multipleFmt, commaDelimited, asStrings[lastIndex])
	}
}

func getDescriptionConstraints(s *mgcSdk.Schema) string {
	constraints := []string{}
	addEnumConstraint(s, &constraints)
	if len(constraints) == 0 {
		addSchemaConstraints(s, &constraints)
	}

	return formatAlternatives("%s", "%s and %s", constraints)
}

func getMinNumber(s *mgcSdk.Schema) *float64 {
	if s.Min == nil {
		return nil
	}

	minInt := *s.Min
	if s.ExclusiveMin {
		minInt += 1
	}

	return &minInt
}

func getMaxNumber(s *mgcSdk.Schema) *float64 {
	if s.Max == nil {
		return nil
	}

	maxInt := *s.Max
	if s.ExclusiveMax {
		maxInt += 1
	}

	return &maxInt
}

func shouldRecommendHelpValue(s *mgcSdk.Schema) bool {
	if s.Type != nil &&
		s.Type.Includes("integer") || s.Type.Includes("number") || s.Type.Includes("boolean") || s.Type.Includes("string") {
		return false
	}

	if s.Type != nil &&
		s.Type.Includes("array") {
		if s.Items != nil && s.Items.Value != nil {
			return shouldRecommendHelpValue((*mgcSdk.Schema)(s.Items.Value))
		}
		return true
	}

	return true
}

type HumanReadableConstraints struct {
	Description     string
	Message         string
	ChildrenMessage string
	Children        []*HumanReadableConstraints
}

func NewHumanReadableConstraints(schema *mgcSchemaPkg.Schema) (c *HumanReadableConstraints) {
	c = specificHumanReadableConstraints(schema)
	if c == nil {
		return
	}
	c.Description = getHumanReadableConstraintsDescription(schema)
	addDefaultHumanReadableConstraint(schema, &c.Children)
	addExampleHumanReadableConstraints(schema, &c.Children)
	return
}

func specificHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	if len(schema.Enum) > 0 {
		return newEnumHumanReadableConstraints(schema)
	}

	if len(schema.OneOf) > 0 {
		return newXOfHumanReadableConstraints("Exactly one of the following must apply", schema, schema.OneOf)
	}
	if len(schema.AnyOf) > 0 {
		return newXOfHumanReadableConstraints("At least one of the following must apply", schema, schema.AnyOf)
	}

	if schema.Type != nil {
		switch {
		case schema.Type.Includes("boolean"):
			return newBooleanHumanReadableConstraints(schema)
		case schema.Type.Includes("string"):
			return newStringHumanReadableConstraints(schema)
		case schema.Type.Includes("integer"), schema.Type.Includes("number"):
			return newNumberHumanReadableConstraints(schema)
		case schema.Type.Includes("array"):
			return newArrayHumanReadableConstraints(schema)
		case schema.Type.Includes("object"):
			return newObjectHumanReadableConstraints(schema)
		default:
			if schema.Not != nil {
				return newNotHumanReadableConstraints(schema)
			}
			return newAnyHumanReadableConstraints(schema)
		}
	}
	if schema.Not != nil {
		return newNotHumanReadableConstraints(schema)
	}
	return newAnyHumanReadableConstraints(schema)

}

func getHumanReadableConstraintsDescription(schema *mgcSchemaPkg.Schema) string {
	description := strings.Trim(getSchemaDescription(schema), "\t\n\r ")
	if description == "" {
		return ""
	}

	if !unicode.IsPunct(rune(description[len(description)-1])) {
		description += "."
	}

	return description
}

func addAnyValueHumanReadableConstraint(format string, value any, constraints *[]*HumanReadableConstraints) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	*constraints = append(*constraints, &HumanReadableConstraints{
		Message: fmt.Sprintf(format, data),
	})
}

func addDefaultHumanReadableConstraint(schema *mgcSchemaPkg.Schema, constraints *[]*HumanReadableConstraints) {
	if schema.Default != nil && len(schema.Enum) == 0 {
		addAnyValueHumanReadableConstraint("If no value is provided, then %s is used", schema.Default, constraints)
	}
}

func addExampleHumanReadableConstraints(schema *mgcSchemaPkg.Schema, constraints *[]*HumanReadableConstraints) {
	if schema.Example != nil && len(schema.Enum) == 0 {
		addAnyValueHumanReadableConstraint("Example value: %s", schema.Example, constraints)
	}
}

func newRefHumanReadableConstraints(ref string) *HumanReadableConstraints {
	return &HumanReadableConstraints{Message: fmt.Sprintf("Previously described [%s]", ref)}
}

func newSchemaRefHumanReadableConstraints(schemaRef *mgcSchemaPkg.SchemaRef) *HumanReadableConstraints {
	if schemaRef == nil {
		return nil
	}

	if schemaRef.Ref != "" {
		// avoid infinite loops on recursive types
		return newRefHumanReadableConstraints(schemaRef.Ref)
	}
	return NewHumanReadableConstraints((*mgcSchemaPkg.Schema)(schemaRef.Value))
}

func newEnumHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	var defVal string
	if schema.Default != nil {
		data, err := json.Marshal(schema.Default)
		if err == nil {
			defVal = string(data)
		}
	}

	children := make([]*HumanReadableConstraints, len(schema.Enum))
	for i, s := range getEnumAsString(schema) {
		if s == defVal {
			s += " (default value)"
		}
		children[i] = &HumanReadableConstraints{Message: s}
	}

	return &HumanReadableConstraints{
		ChildrenMessage: getPlural(len(children), "Must be exactly", "One of the values"),
		Children:        children,
	}
}

func newBooleanHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	return &HumanReadableConstraints{
		Message: "Boolean value",
	}
}

func getHumanConstraintsChildren(constraints []string) (childrenMessage string, children []*HumanReadableConstraints) {
	if len(constraints) == 0 {
		return
	}

	childrenMessage = "With the following constraints"
	children = make([]*HumanReadableConstraints, len(constraints))
	for i, c := range constraints {
		children[i] = &HumanReadableConstraints{Message: c}
	}
	return
}

func newStringHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	var constraints []string
	addStringConstraints(schema, &constraints)
	childrenMessage, children := getHumanConstraintsChildren(constraints)
	return &HumanReadableConstraints{
		Message:         "String value",
		ChildrenMessage: childrenMessage,
		Children:        children,
	}
}

func newNumberHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	var constraints []string
	addNumberConstraints(schema, &constraints)
	childrenMessage, children := getHumanConstraintsChildren(constraints)

	typeStr := "Number"
	if len(schema.Type.Slice()) > 0 {
		typeStr = schema.Type.Slice()[0]
	}
	return &HumanReadableConstraints{
		Message:         strings.ToUpper(typeStr[:1]) + typeStr[1:] + " value",
		ChildrenMessage: childrenMessage,
		Children:        children,
	}
}

func newArrayHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	var constraints []string
	addArrayConstraints(schema, &constraints)
	childrenMessage, children := getHumanConstraintsChildren(constraints)

	if itemsConstraints := newSchemaRefHumanReadableConstraints(schema.Items); itemsConstraints != nil {
		if childrenMessage == "" {
			childrenMessage = "Array where each item is"
		}
		children = append(children, itemsConstraints)
	}

	return &HumanReadableConstraints{
		Message:         "Array value",
		ChildrenMessage: childrenMessage,
		Children:        children,
	}
}

func newMapHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	var constraints []string
	addCountConstraints(schema.MinProps, schema.MaxProps, "property", "properties", &constraints)
	childrenMessage, children := getHumanConstraintsChildren(constraints)

	if valueConstraints := newSchemaRefHumanReadableConstraints(schema.AdditionalProperties.Schema); valueConstraints != nil {
		if childrenMessage == "" {
			childrenMessage = "Mapping where each value is"
		}
		children = append(children, valueConstraints)
	}

	return &HumanReadableConstraints{
		Message:         "Mapping of string keys to values",
		ChildrenMessage: childrenMessage,
		Children:        children,
	}
}

func newObjectHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	if schema.AdditionalProperties.Has != nil && *schema.AdditionalProperties.Has {
		return newMapHumanReadableConstraints(schema)
	}

	var constraints []string
	addCountConstraints(schema.MinProps, schema.MaxProps, "property", "properties", &constraints)
	childrenMessage, children := getHumanConstraintsChildren(constraints)

	if len(schema.Properties) > 0 {
		properties := make([]*HumanReadableConstraints, 0, len(schema.Properties))
		for k, schemaRef := range schema.Properties {
			if p := newSchemaRefHumanReadableConstraints(schemaRef); p != nil {
				var requiredMarker string
				if slices.Contains(schema.Required, k) {
					requiredMarker = "(required) "
				}
				p.Description = fmt.Sprintf("%q: %s%s", k, requiredMarker, p.Description)
				properties = append(properties, p)
			}
		}
		slices.SortFunc(properties, func(a, b *HumanReadableConstraints) int {
			return strings.Compare(a.Message, b.Message)
		})

		children = append(children, &HumanReadableConstraints{
			Message:  "Object with the following properties",
			Children: properties,
		})
	}

	return &HumanReadableConstraints{
		Message:         "Object value",
		ChildrenMessage: childrenMessage,
		Children:        children,
	}
}

func newXOfHumanReadableConstraints(text string, schema *mgcSchemaPkg.Schema, schemaRefs mgcSchemaPkg.SchemaRefs) *HumanReadableConstraints {
	children := make([]*HumanReadableConstraints, len(schemaRefs))
	for i, schemaRef := range schemaRefs {
		children[i] = newSchemaRefHumanReadableConstraints(schemaRef)
	}

	return &HumanReadableConstraints{
		Message:  text,
		Children: children,
	}
}

func newNotHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	return &HumanReadableConstraints{
		Message:  "Not",
		Children: []*HumanReadableConstraints{newSchemaRefHumanReadableConstraints(schema.Not)},
	}
}

func newAnyHumanReadableConstraints(schema *mgcSchemaPkg.Schema) *HumanReadableConstraints {
	return &HumanReadableConstraints{Message: "Any value"}
}
