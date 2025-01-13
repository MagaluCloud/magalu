package transform

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"go.uber.org/zap"
)

// Recursively checks whenever the given schema needs transformation
func needsTransformation(logger *zap.SugaredLogger, schema *core.Schema, transformationKey string) (bool, error) {
	t := &commonSchemaTransformer[bool]{
		logger: logger,
		tKey:   transformationKey,
		transform: func(logger *zap.SugaredLogger, transformers []transformer, value bool) (bool, error) {
			return true, nil
		},
		transformArray:       transformArrayNeedsTransformation,
		transformObject:      transformObjectNeedsTransformation,
		transformConstraints: transformConstraintsNeedsTransformation,
	}
	return mgcSchemaPkg.Transform[bool](t, schema, false)
}

func transformArrayNeedsTransformation(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[bool], schema *core.Schema, itemSchema *core.Schema, value bool) (bool, error) {
	if itemSchema == nil {
		return value, nil
	}
	return mgcSchemaPkg.Transform(t, itemSchema, value)
}

func transformObjectNeedsTransformation(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[bool], schema *core.Schema, value bool) (bool, error) {
	return mgcSchemaPkg.TransformObjectProperties(schema, value, func(propName string, propSchema *core.Schema, value bool) (bool, error) {
		value, err := mgcSchemaPkg.Transform(t, propSchema, value)
		if err != nil {
			return value, err
		}
		if value {
			return true, mgcSchemaPkg.TransformStop
		}
		return false, nil
	})
}

func transformConstraintsNeedsTransformation(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[bool], kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value bool) (bool, error) {
	value, err := mgcSchemaPkg.TransformSchemasArray(t, schemaRefs, value)
	if err != nil {
		return value, err
	}
	if value {
		return true, mgcSchemaPkg.TransformStop
	}
	return false, nil

}
