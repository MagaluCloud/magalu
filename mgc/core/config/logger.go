package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logfilterSchema() *core.Schema {
	s := schema.NewStringSchema()
	s.Pattern = "^((((((debug|info|warn|error|panic|dpanic|fatal)\\+?)(,(debug|info|warn|error|panic|dpanic|fatal)\\+?)*)|\\*)(:\\S+)?\\s?)*)$"
	s.Description = "Default log filter to be used. See https://github.com/moul/zapfilter#zapfilter for reference and examples"
	return s
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
				"lowercase",
			},
		}
	case reflect.TypeOf(zapcore.TimeEncoder(timeEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"rfc3339nano", "RFC3339Nano",
				"rfc3339", "RFC3339",
				"iso8601", "ISO8601",
				"millis",
				"nanos",
				"epoch",
			}}
	case reflect.TypeOf(zapcore.DurationEncoder(durationEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"string",
				"nanos",
				"ms",
				"s",
			},
		}
	case reflect.TypeOf(zapcore.CallerEncoder(callerEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"full",
				"short",
			},
		}
	case reflect.TypeOf(zapcore.NameEncoder(nameEncoder)):
		return &jsonschema.Schema{
			Type: "string",
			Enum: []any{
				"full",
			},
		}
	default:
		return nil
	}
}
