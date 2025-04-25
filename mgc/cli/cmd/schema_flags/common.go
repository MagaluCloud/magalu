package schema_flags

import "slices"

type schemaFlagValueCommon struct {
	desc     SchemaFlagValueDesc
	rawValue string
	changed  bool
}

var _ SchemaFlagValue = (*schemaFlagValueCommon)(nil)

func initSchemaFlagValueCommon(desc SchemaFlagValueDesc) (o schemaFlagValueCommon) {
	o.desc = desc
	o.rawValue = o.RawDefaultValue()
	return o
}

func newSchemaFlagValueCommon(desc SchemaFlagValueDesc) *schemaFlagValueCommon {
	o := initSchemaFlagValueCommon(desc)
	return &o
}

func (o *schemaFlagValueCommon) String() string {
	return o.rawValue
}

func (o *schemaFlagValueCommon) Set(rawValue string) error {
	o.rawValue = rawValue
	o.changed = true
	return nil
}

func (o *schemaFlagValueCommon) Type() string {
	return o.desc.FlagType()
}

func (o *schemaFlagValueCommon) Desc() SchemaFlagValueDesc {
	return o.desc
}

func (o *schemaFlagValueCommon) RawDefaultValue() string {
	result := o.desc.RawDefaultValue()
	allowedValues := []string{"\"br-se1\"", "\"prod\"", "[\"network\",\"image\",\"machine-type\"]"}
	if slices.Contains(allowedValues, result) || o.desc.IsRequired {
		return result
	}
	return ""
}

func (o *schemaFlagValueCommon) RawNoOptDefVal() string {
	return ""
}

func (o *schemaFlagValueCommon) Usage() string {
	return o.desc.Usage()
}

func (o *schemaFlagValueCommon) Parse() (value any, err error) {
	if o.rawValue == ValueHelpIsRequired {
		return nil, ErrWantHelp
	}

	value, err = parseJSONFlagValue[any](o.rawValue)
	if err != nil {
		value = o.rawValue
		err = nil
	}

	return
}

func (o *schemaFlagValueCommon) Changed() bool {
	return o.changed
}
