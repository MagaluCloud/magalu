package transform

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"go.uber.org/zap"
)

func doTransformsToValue(logger *zap.SugaredLogger, transformers []transformer, value any) (result any, err error) {
	result = value
	for _, t := range transformers {
		result, err = t.TransformValue(result)
		if err != nil {
			logger.Debugw("transformation attempt failed", "value", value)
			return
		}
	}
	if result != value {
		logger.Debugw("transformed value", "input", value, "output", result)
	}
	return
}

// Recursively transforms the value based on the schema that may contain transformations
// If the schema doesn't contain any transformation, then the value is unchanged
func transformValue(logger *zap.SugaredLogger, schema *core.Schema, transformationKey string, value any) (any, error) {
	t := &commonSchemaTransformer[any]{
		logger:               logger,
		tKey:                 transformationKey,
		transform:            doTransformsToValue,
		transformArray:       transformArrayValue,
		transformObject:      transformObjectValue,
		transformConstraints: transformConstraintsValue,
	}
	return mgcSchemaPkg.Transform[any](t, schema, value)
}

func transformArrayValue(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[any], schema *core.Schema, itemSchema *core.Schema, value any) (any, error) {
	valueSlice, ok := value.([]any)
	if !ok {
		if value == nil {
			if schema.Nullable {
				return value, nil
			}
			return value, fmt.Errorf("received null for non-nullable schema: %#v", schema)
		}
		return value, fmt.Errorf("expected []any, got %T %#v", value, value)
	}

	cs := utils.NewCOWSliceFunc(valueSlice, utils.IsSameValueOrPointer)
	for i, itemValue := range valueSlice {
		logger.Debugw("transform array item...", "index", i, "itemValue", itemValue)
		convertedValue, err := mgcSchemaPkg.Transform(t, itemSchema, itemValue)
		logger.Debugw("transformed array item", "index", i, "itemValue", itemValue, "error", err)
		if err != nil {
			return value, err
		}
		cs.Set(i, convertedValue)
	}

	valueSlice, _ = cs.Release()
	return valueSlice, nil
}

func transformObjectValue(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[any], schema *core.Schema, value any) (any, error) {
	valueMap, ok := value.(map[string]any)
	if !ok {
		if value == nil {
			if schema.Nullable {
				return value, nil
			}
			return value, fmt.Errorf("received null for non-nullable schema: %#v", schema)
		}
		return value, fmt.Errorf("expected map[string]any, got %T %#v", value, value)
	}
	cm, err := mgcSchemaPkg.TransformObjectProperties(
		schema,
		utils.NewCOWMapFunc(valueMap, utils.IsSameValueOrPointer),
		func(propName string, propSchema *core.Schema, cm *utils.COWMap[string, any],
		) (*utils.COWMap[string, any], error) {
			propValue, ok := valueMap[propName]
			if !ok {
				return cm, nil
			}

			logger.Debugw("transform object property...", "propName", propName, "propValue", propValue)
			convertedFieldValue, err := mgcSchemaPkg.Transform(t, propSchema, propValue)
			logger.Debugw("transformed object property", "propName", propName, "propValue", propValue, "error", err)
			if err != nil {
				return cm, err
			}
			cm.Set(propName, convertedFieldValue)
			return cm, nil
		},
	)
	if err != nil {
		return value, err
	}

	valueMap, _ = cm.Release()
	return valueMap, nil
}

func transformConstraintsValue(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[any], kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value any) (any, error) {
	// TODO: handle kind properly, see https://swagger.io/docs/specification/data-models/oneof-anyof-allof-not/
	return mgcSchemaPkg.TransformSchemasArray(t, schemaRefs, value)
}
