package schema_flags

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
)

type schemaFlagValueArray struct {
	schemaFlagValueMulti
	itemsSchema *core.Schema
}

var _ SchemaFlagValue = (*schemaFlagValueArray)(nil)

func newSchemaFlagValueArray(desc SchemaFlagValueDesc) *schemaFlagValueArray {
	var itemsSchema *core.Schema
	if desc.Schema.Items != nil {
		itemsSchema = (*core.Schema)(desc.Schema.Items.Value)
	}

	return &schemaFlagValueArray{
		initSchemaFlagValueMulti(desc),
		itemsSchema,
	}
}

func (o *schemaFlagValueArray) Type() (s string) {
	s = "array"
	if o.itemsSchema != nil {
		s += fmt.Sprintf("(%s)", getFlagType(o.itemsSchema))
	}
	return s
}

func (o *schemaFlagValueArray) Parse() (value any, err error) {
	if len(o.setValues) == 0 {
		return parseArrayFlagValueSingle(o.itemsSchema, o.rawValue)
	}

	return parseArrayFlagValue(o.itemsSchema, o.setValues)
}
