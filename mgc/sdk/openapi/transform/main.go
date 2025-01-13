package transform

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	"go.uber.org/zap"
)

func getTransformKey(extensionPrefix *string) string {
	if extensionPrefix == nil || *extensionPrefix == "" {
		return ""
	}
	return *extensionPrefix + "-transforms"
}

// The returned function does NOT and should NOT alter the value that was passed by it
// (maps, for example, when passed as input, won't be altered, a new copy will be made)
func New[T any](logger *zap.SugaredLogger, schema *core.Schema, extensionPrefix *string) (func(value T) (T, error), *core.Schema, error) {
	transformationKey := getTransformKey(extensionPrefix)
	if transformationKey == "" {
		return nil, schema, nil
	}

	needs, err := needsTransformation(logger, schema, transformationKey)
	if err != nil {
		return nil, schema, err
	}
	if !needs {
		return nil, schema, nil
	}

	transformedSchema, err := transformSchema(logger, schema, transformationKey, schema)
	if err != nil {
		return nil, schema, err
	}

	return func(value T) (converted T, err error) {
		r, err := transformValue(logger, schema, transformationKey, value)
		if err != nil {
			return
		}
		converted, ok := r.(T)
		if !ok {
			err = fmt.Errorf("invalid conversion result, expected %T, got %+v", converted, r)
			return
		}
		return
	}, transformedSchema, nil
}
