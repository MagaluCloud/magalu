package transform

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"go.uber.org/zap"
)

func doTransformsToSchema(logger *zap.SugaredLogger, transformers []transformer, value *mgcSchemaPkg.COWSchema, transformationKey string) (result *mgcSchemaPkg.COWSchema, err error) {
	result = value
	for _, t := range transformers {
		result, err = t.TransformSchema(result)
		if err != nil {
			logger.Debugw("transformation attempt failed", "value", value)
			return
		}
	}
	value.ExtensionsCOW().Delete(transformationKey)
	if result != value {
		logger.Debugw("transformed schema", "input", value.Peek(), "output", result.Peek())
	}
	return
}

func transformSchema(logger *zap.SugaredLogger, schema *core.Schema, transformationKey string, value *core.Schema) (*core.Schema, error) {
	t := &commonSchemaTransformer[*mgcSchemaPkg.COWSchema]{
		logger: logger,
		tKey:   transformationKey,
		transform: func(logger *zap.SugaredLogger, transformers []transformer, value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
			return doTransformsToSchema(logger, transformers, value, transformationKey)
		},
		transformArray:       transformArraySchema,
		transformObject:      transformObjectSchema,
		transformConstraints: transformConstraintsSchema,
	}
	cowSchema := mgcSchemaPkg.NewCOWSchema(value)
	cowSchema, err := mgcSchemaPkg.Transform[*mgcSchemaPkg.COWSchema](t, schema, cowSchema)
	if err != nil {
		return value, err
	}
	return cowSchema.Peek(), nil
}

func transformArraySchema(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[*mgcSchemaPkg.COWSchema], schema *core.Schema, itemSchema *core.Schema, value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
	itemsCow := value.ItemsCOW().ValueCOW()
	_, err := mgcSchemaPkg.Transform(t, itemSchema, itemsCow)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func transformObjectSchema(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[*mgcSchemaPkg.COWSchema], schema *core.Schema, value *mgcSchemaPkg.COWSchema) (*mgcSchemaPkg.COWSchema, error) {
	_, err := mgcSchemaPkg.TransformObjectProperties(
		schema,
		value.PropertiesCOW(),
		func(propName string, propSchema *core.Schema, propertiesCow *utils.COWMapOfCOW[string, *mgcSchemaPkg.SchemaRef, *mgcSchemaPkg.COWSchemaRef],
		) (*utils.COWMapOfCOW[string, *mgcSchemaPkg.SchemaRef, *mgcSchemaPkg.COWSchemaRef], error) {
			propSchemaCow, ok := propertiesCow.GetCOW(propName)
			if !ok {
				return nil, fmt.Errorf("schema missing property %q", propName) // this should never happen
			}

			_, err := mgcSchemaPkg.Transform(t, propSchema, propSchemaCow.ValueCOW())
			if err != nil {
				return nil, err
			}
			return propertiesCow, nil
		},
	)
	if err != nil {
		return value, err
	}
	return value, nil
}

func transformConstraintsSchema(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[*mgcSchemaPkg.COWSchema], kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value *mgcSchemaPkg.COWSchema) (result *mgcSchemaPkg.COWSchema, err error) {
	result = value

	if kind == mgcSchemaPkg.ConstraintNot {
		_, err = mgcSchemaPkg.Transform(t, (*mgcSchemaPkg.Schema)(schemaRefs[0].Value), value.NotCOW().ValueCOW())
		return
	}

	var constraintCow *utils.COWSliceOfCOW[*mgcSchemaPkg.SchemaRef, *mgcSchemaPkg.COWSchemaRef]
	switch kind {
	case mgcSchemaPkg.ConstraintAllOf:
		constraintCow = value.AllOfCOW()
	case mgcSchemaPkg.ConstraintAnyOf:
		constraintCow = value.AnyOfCOW()
	case mgcSchemaPkg.ConstraintOneOf:
		constraintCow = value.OneOfCOW()
	default:
		return value, fmt.Errorf("unknown constraint kind: %q", kind)
	}

	constraintCow.ForEachCOW(func(i int, cowRef *mgcSchemaPkg.COWSchemaRef) (run bool) {
		itemSchema := cowRef.Peek()
		if itemSchema == nil {
			return true
		}

		_, err = mgcSchemaPkg.Transform(t, (*mgcSchemaPkg.Schema)(itemSchema.Value), cowRef.ValueCOW())
		return err == nil
	})

	return
}
