package schema_flags

import (
	"errors"
	"reflect"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/getkin/kin-openapi/openapi3"
	flag "github.com/spf13/pflag"
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

	if desc.Schema.Type != nil {
		switch {
		case desc.Schema.Type.Includes("array"):
			return newSchemaFlagValueArray(desc)
		case desc.Schema.Type.Includes("boolean"):
			return newSchemaFlagValueBool(desc)
		case desc.Schema.Type.Includes("string"):
			return newSchemaFlagValueString(desc)
		case desc.Schema.Type.Includes("object"):
			return newSchemaFlagValueObject(desc)
		}
	}
	return newSchemaFlagValueCommon(desc)

}

func newFlag(value SchemaFlagValue) *flag.Flag {
	return &flag.Flag{
		Name:        string(value.Desc().FlagName),
		DefValue:    value.RawDefaultValue(),
		NoOptDefVal: value.RawNoOptDefVal(),
		Usage:       value.Usage(),
		Value:       value,
		Hidden:      value.Desc().IsHidden,
	}
}

func NewSchemaFlag(
	container *core.Schema, // the object that contains this property
	propName string, // name inside ParametersSchema()/ConfigsSchema() Properties
	flagName flag.NormalizedName, // public/user-visible name after normalization
	isRequired bool,
	isConfig bool,
	isHidden bool,
) *flag.Flag {
	schema := (*core.Schema)(container.Properties[propName].Value)
	return newFlag(newSchemaFlagValue(SchemaFlagValueDesc{
		container,
		schema,
		propName,
		flagName,
		isRequired,
		isConfig,
		isHidden,
	}))
}

// Proxy set/parse calls to the given functions
func NewProxyFlag(
	desc SchemaFlagValueDesc,
	proxy ProxyFlagSpec,
) *flag.Flag {
	return newFlag(newSchemaFlagValueProxy(desc, proxy))
}

var (
	ErrNoFlagValue  = errors.New("no value")
	ErrRequiredFlag = errors.New("is required")
	ErrWantHelp     = errors.New("help is needed")
)

func getUnchangedFlagValue(fv SchemaFlagValue, cfg *config.Config) (value any, err error) {
	desc := fv.Desc()

	if desc.IsConfig {
		if err = cfg.Get(desc.PropName, &value); err == nil && value != nil {
			return
		}
	}

	if desc.IsRequired {
		value = desc.Schema.Default
	}

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

	// mgc config set [key] [value-integer]
	if desc.Schema.OneOf != nil {
		for _, oneOf := range desc.Schema.OneOf {
			typeOfValue := reflect.TypeOf(value).String()
			if typeOfValue == "float64" {
				typeOfValue = "integer"
			}
			if typeOfValue == oneOf.Value.Type.Slice()[0] {
				desc.Schema.Type = oneOf.Value.Type
				break
			}
		}
	}

	err = desc.Schema.VisitJSON(value, openapi3.MultiErrors())

	return
}
