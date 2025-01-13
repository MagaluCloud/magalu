package transform

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"go.uber.org/zap"
)

// Common pattern that checks existing specs, if they exist then call transformSpecs(),
// otherwise process Arrays and Objects.
//
// Scalars are passed thru while Constraints() are recursively processed.
type commonSchemaTransformer[T any] struct {
	logger               *zap.SugaredLogger
	tKey                 string
	transform            func(logger *zap.SugaredLogger, transformers []transformer, value T) (T, error)
	transformArray       func(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[T], schema *core.Schema, itemSchema *core.Schema, value T) (T, error)
	transformObject      func(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[T], schema *core.Schema, value T) (T, error)
	transformConstraints func(logger *zap.SugaredLogger, t mgcSchemaPkg.Transformer[T], kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value T) (T, error)
}

func (t *commonSchemaTransformer[T]) Transform(schema *core.Schema, value T) (T, error) {
	transformers, err := getTransformers(schema.Extensions, t.tKey)
	if err != nil {
		t.logger.Warnw("getTransformers() failed", "schema", schema, "value", value, "error", err)
		return value, err
	}
	if len(transformers) > 0 {
		t.logger.Debugw("transform...", "schema", schema, "transformers", transformers, "value", value)
		value, err = t.transform(t.logger, transformers, value)
		if err == nil {
			err = mgcSchemaPkg.TransformStop
		}
		t.logger.Debugw("transformed", "schema", schema, "transformers", transformers, "value", value)
	}
	return value, err
}

func (t *commonSchemaTransformer[T]) Scalar(schema *core.Schema, value T) (T, error) {
	return value, nil
}

func (t *commonSchemaTransformer[T]) Array(schema *core.Schema, itemSchema *core.Schema, value T) (T, error) {
	if itemSchema == nil {
		return value, nil
	}
	t.logger.Debugw("transform array...", "schema", schema, "itemSchema", itemSchema, "value", value)
	value, err := t.transformArray(t.logger, t, schema, itemSchema, value)
	t.logger.Debugw("transformed array", "schema", schema, "itemSchema", itemSchema, "value", value, "error", err)
	return value, err
}

func (t *commonSchemaTransformer[T]) Constraints(kind mgcSchemaPkg.ConstraintKind, schemaRefs mgcSchemaPkg.SchemaRefs, value T) (T, error) {
	t.logger.Debugw("transform constraints...", "kind", kind, "schemaRefs", schemaRefs, "value", value)
	value, err := t.transformConstraints(t.logger, t, kind, schemaRefs, value)
	t.logger.Debugw("transformed constraints", "kind", kind, "schemaRefs", schemaRefs, "value", value, "error", err)
	return value, err
}

func (t *commonSchemaTransformer[T]) Object(schema *core.Schema, value T) (T, error) {
	t.logger.Debugw("transform object...", "schema", schema, "value", value)
	value, err := t.transformObject(t.logger, t, schema, value)
	t.logger.Debugw("transformed object", "schema", schema, "value", value, "error", err)
	return value, err
}

var _ mgcSchemaPkg.Transformer[any] = (*commonSchemaTransformer[any])(nil)
