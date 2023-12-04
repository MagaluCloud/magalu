package schema_flags

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
	return o.desc.RawDefaultValue()
}

func (o *schemaFlagValueCommon) RawNoOptDefVal() string {
	return ""
}

func (o *schemaFlagValueCommon) Usage() string {
	return o.desc.Usage()
}

func (o *schemaFlagValueCommon) Parse() (value any, err error) {
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
