package schema_flags

import (
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
	flag "github.com/spf13/pflag"
	"magalu.cloud/core"
	"magalu.cloud/core/config"
)

type SchemaFlagValue interface {
	flag.Value

	Desc() SchemaFlagValueDesc

	RawDefaultValue() string
	RawNoOptDefVal() string
	Usage() string
	Parse() (any, error)

	Changed() bool // if the value was explicitly set
}

func newSchemaFlagValue(desc SchemaFlagValueDesc) SchemaFlagValue {
	switch desc.Schema.Type {
	case "array":
		return newSchemaFlagValueArray(desc)
	case "boolean":
		return newSchemaFlagValueBool(desc)
	case "string":
		return newSchemaFlagValueString(desc)
	case "object":
		return newSchemaFlagValueObject(desc)
	default:
		return newSchemaFlagValueCommon(desc)
	}
}

func newFlag(value SchemaFlagValue) *flag.Flag {
	return &flag.Flag{
		Name:        string(value.Desc().FlagName),
		DefValue:    value.RawDefaultValue(),
		NoOptDefVal: value.RawNoOptDefVal(),
		Usage:       value.Usage(),
		Value:       value,
	}
}

func NewSchemaFlag(
	container *core.Schema, // the object that contains this property
	propName string, // name inside ParametersSchema()/ConfigsSchema() Properties
	flagName flag.NormalizedName, // public/user-visible name after normalization
	isRequired bool,
	isConfig bool,
) *flag.Flag {
	schema := (*core.Schema)(container.Properties[propName].Value)
	return newFlag(newSchemaFlagValue(SchemaFlagValueDesc{
		container,
		schema,
		propName,
		flagName,
		isRequired,
		isConfig,
	}))
}

var (
	ErrNoFlagValue  = errors.New("no value")
	ErrRequiredFlag = errors.New("is required")
)

func getUnchangedFlagValue(fv SchemaFlagValue, cfg *config.Config) (value any, err error) {
	desc := fv.Desc()

	if desc.IsConfig {
		if err = cfg.Get(desc.PropName, &value); err == nil && value != nil {
			return
		}
	}

	value = desc.Schema.Default
	return
}

func getFlagValue(fv SchemaFlagValue, cfg *config.Config) (value any, err error) {
	if !fv.Changed() {
		return getUnchangedFlagValue(fv, cfg)
	}

	return fv.Parse()
}

func GetFlagValue(f *flag.Flag, cfg *config.Config) (value any, err error) {
	fv := f.Value.(SchemaFlagValue)
	value, err = getFlagValue(fv, cfg)
	if err != nil {
		return
	}

	desc := fv.Desc()

	if value == nil && !desc.Schema.Nullable {
		if desc.IsRequired {
			err = ErrRequiredFlag
		} else {
			err = ErrNoFlagValue
		}
		return
	}

	err = desc.Schema.VisitJSON(value, openapi3.MultiErrors())

	return
}
