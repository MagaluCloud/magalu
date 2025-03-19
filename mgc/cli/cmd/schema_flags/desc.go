package schema_flags

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	flag "github.com/spf13/pflag"
)

type SchemaFlagValueDesc struct {
	Container *core.Schema
	Schema    *core.Schema
	PropName  string              // name inside ParametersSchema()/ConfigsSchema() Properties
	FlagName  flag.NormalizedName // public/user-visible name after normalization

	IsRequired bool
	IsConfig   bool
	IsHidden   bool
}

const (
	FlagTypeFile      = "file"
	FlagTypeDirectory = "directory"
)

func getFlagType(schema *core.Schema) string {
	if len(schema.Enum) > 0 {
		return "enum"
	}

	if schema.Format != "" {
		return schema.Format
	}

	if mt := schema.Extensions["x-contentMediaType"]; mt != nil {
		if mt == "inode/directory" {
			return FlagTypeDirectory
		}
		return FlagTypeFile
	}

	if schema.Type != nil && len(schema.Type.Slice()) == 0 {
		if (mgcSchemaPkg.CheckSimilarJsonSchemas(schema, &mgcSchemaPkg.Schema{}) || // this is like a bug in the schema, but config set takes it
			mgcSchemaPkg.CheckSimilarJsonSchemas(schema, mgcSchemaPkg.NewAnySchema())) {
			return "anyValue"
		}
		if len(schema.AnyOf) > 0 {
			return "anyOf"
		}
		if len(schema.OneOf) > 0 {
			return "oneOf"
		}
		return "anyValue"
	}
	if len(schema.Type.Slice()) > 1 {
		fmt.Println("REMOVE ME - 20250313-1857   =>", schema.Type.Slice())
	}
	return schema.Type.Slice()[0]
}

// To be used in flag.Value.Type().
//
// Cobra will show this one in their usage and decide some behavior
// based on it, there are special cases for "bool", "string"...
func (d *SchemaFlagValueDesc) FlagType() string {
	return getFlagType(d.Schema)
}

// To be used in flag.DefValue.
func (d *SchemaFlagValueDesc) RawDefaultValue() string {
	if d.Schema.Default == nil {
		return ""
	}

	data, err := json.Marshal(d.Schema.Default)
	if err != nil {
		logger().Warnw(
			"could not convert flag default value to string",
			"defaultValue", d.Schema.Default,
			"flagName", d.FlagName,
			"propName", d.PropName,
			"error", err,
		)
		return ""
	}

	return string(data)
}

var flagDescriptionUsageRe = regexp.MustCompile("([^'\"` ]+)=([^'\"` ]+)")

func (d SchemaFlagValueDesc) fixDescriptionFlagUsage(input string) string {
	enumAsString := getEnumAsString(d.Schema)

	prefixToMatch := d.PropName + "="
	replacedPrefix := string("--" + d.FlagName + "=")

	// replace any mentions propName=value to --flagName=value, so it makes sense
	return flagDescriptionUsageRe.ReplaceAllStringFunc(input, func(s string) string {
		if !strings.HasPrefix(s, prefixToMatch) {
			return s
		}

		value := s[len(prefixToMatch):]

		if len(enumAsString) > 0 {
			data, _ := json.Marshal(value)
			valueAsJson := string(data)
			// some schemas come with a generic description, but enums restrict the fields to be used.
			for _, restrictedValue := range enumAsString {
				if value == restrictedValue || valueAsJson == restrictedValue {
					return replacedPrefix + restrictedValue
				}
			}
			return replacedPrefix + enumAsString[0]
		}

		return replacedPrefix + value
	})
}

func getSchemaDescription(s *mgcSchemaPkg.Schema) string {
	if s.Description == "" {
		return s.Title
	}

	if s.Title == "" {
		return s.Description
	}

	normTitle := strings.ToLower(s.Title)
	normDescription := strings.ToLower(s.Description)

	if strings.Contains(normTitle, normDescription) {
		return s.Title
	}

	if strings.Contains(normDescription, normTitle) {
		return s.Description
	}

	return fmt.Sprintf("%s: %s", s.Title, s.Description)
}

func (d SchemaFlagValueDesc) Description() (description string) {
	description = getSchemaDescription(d.Schema)

	// spf13/pflag have UnquoteUsage() that messes up with back quotes, so remove them
	description = strings.ReplaceAll(description, "`", "'")

	description = d.fixDescriptionFlagUsage(description)

	constraints := getDescriptionConstraints(d.Schema)
	if constraints != "" {
		if description != "" {
			description += " "
		}
		description += fmt.Sprintf("(%s)", constraints)
	}

	return description
}

func (d SchemaFlagValueDesc) Usage() (usage string) {
	usage = d.Description()

	if shouldRecommendHelpValue(d.Schema) {
		if usage != "" {
			usage += "\n"
		}
		usage += fmt.Sprintf("Use --%s=%s for more details", d.FlagName, ValueHelpIsRequired)
	}

	if d.IsRequired {
		usage += " (required)"
	}

	return
}

func (d SchemaFlagValueDesc) HumanReadableConstraints() (c *HumanReadableConstraints) {
	c = NewHumanReadableConstraints(d.Schema)
	if c == nil {
		return
	}

	c.Description = d.fixDescriptionFlagUsage(c.Description)
	return
}

func (d SchemaFlagValueDesc) ChildrenFlags() (children []SchemaFlagValueDesc) {
	properties := utils.SortedMapIterator(mgcSchemaPkg.CollectAllObjectPropertySchemas(d.Schema))
	if len(properties) == 0 {
		return
	}

	children = make([]SchemaFlagValueDesc, 0, len(properties))
	for _, pair := range properties {
		for i, ps := range pair.Value {
			var flagName flag.NormalizedName
			if len(pair.Value) > 1 {
				flagName = flag.NormalizedName(fmt.Sprintf(".%s.%d.", ps.Field, i))
			}

			children = append(children, SchemaFlagValueDesc{
				Container: ps.Container,
				Schema:    ps.PropSchema,
				PropName:  ps.PropName,
				FlagName:  flagName, // actually it's the infix to resolve conflicts

				IsRequired: false,
				IsConfig:   d.IsConfig,
			})
		}
	}

	return
}
