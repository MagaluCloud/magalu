package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"magalu.cloud/core"
	coreLogger "magalu.cloud/core/logger"
	"magalu.cloud/core/schema"
)

type pkgSymbol struct{}

var pkgLogger *zap.SugaredLogger

func logger() *zap.SugaredLogger {
	if pkgLogger == nil {
		pkgLogger = coreLogger.New[pkgSymbol]()
	}
	return pkgLogger
}

func loggerSchema() (*core.Schema, error) {
	reflector := jsonschema.Reflector{Mapper: zapMapper}
	s, err := schema.ToCoreSchema(reflector.Reflect(zap.Config{}))
	if err != nil {
		return nil, fmt.Errorf("unable to create JSON Schema for type '%T': %w", zap.Config{}, err)
	}

	removeRequired(s)

	s.Description = "Logger configuration. For more information see https://pkg.go.dev/go.uber.org/zap#Config"

	return s, nil
}

func removeRequired(s *core.Schema) {
	if len(s.Required) != 0 {
		s.Required = []string{}
	}
	for _, ref := range s.Properties {
		if ref.Value != nil {
			removeRequired(((*core.Schema)(ref.Value)))
		}
	}
}

func levelEncoder(zapcore.Level, zapcore.PrimitiveArrayEncoder)        {}
func timeEncoder(time.Time, zapcore.PrimitiveArrayEncoder)             {}
func durationEncoder(time.Duration, zapcore.PrimitiveArrayEncoder)     {}
func callerEncoder(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder) {}
func nameEncoder(string, zapcore.PrimitiveArrayEncoder)                {}

// The zapMapper function is necessary because some zapcore.EncoderConfig
// fields are functions, and it's not possiple to reflect these fields
// to a jsonschema.Schema. So we have to tell jsonschema how to handle them

func zapMapper(t reflect.Type) *jsonschema.Schema {
	switch t {
	case reflect.TypeOf(zap.AtomicLevel{}):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"debug",
				"info",
				"warn",
				"error",
				"panic",
				"dpanic",
				"fatal",
			},
		}
	case reflect.TypeOf(zapcore.LevelEncoder(levelEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"capital",
				"color",
				"capitalColor",
				"default",
			}}
	case reflect.TypeOf(zapcore.TimeEncoder(timeEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"rfc3339nano", "RFC3339Nano",
				"rfc3339", "RFC3339",
				"iso8601", "ISO8601",
				"millis",
				"nanos",
				"default",
			}}
	case reflect.TypeOf(zapcore.DurationEncoder(durationEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"string",
				"nanos",
				"ms",
				"default",
			},
		}
	case reflect.TypeOf(zapcore.CallerEncoder(callerEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"string",
				"default",
			},
		}
	case reflect.TypeOf(zapcore.NameEncoder(nameEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"full",
				"default",
			},
		}
	default:
		return nil
	}
}
