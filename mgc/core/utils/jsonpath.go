package utils

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"github.com/dustin/go-humanize"
)

var hasKeyFunc = gval.Function("hasKey", func(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("'hasKey' jsonpath function expects two arguments, a map[string]any and a string (key). Got %v instead", args...)
	}

	m, ok := args[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("'hasKey' jsonpath function expected map[string]any as first argument, got %T instead: %v", args[0], args[0])
	}

	k, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("'hasKey' jsonpath function expected string as second argument, got %T instead: %v", args[1], args[1])
	}

	_, hasKey := m[k]
	return hasKey, nil
})

var fileSizeFunc = gval.Function("fileSize", func(args ...any) (result any, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("'fileSize' jsonpath function expects a single argument with the value or slice to be formatted. Got %#v instead", args)
	}

	fileSize := func(val any) (string, error) {
		switch v := val.(type) {
		case int:
			return humanize.Bytes(uint64(v)), nil
		case int8:
			return humanize.Bytes(uint64(v)), nil
		case int16:
			return humanize.Bytes(uint64(v)), nil
		case int32:
			return humanize.Bytes(uint64(v)), nil
		case int64:
			return humanize.Bytes(uint64(v)), nil
		case uint:
			return humanize.Bytes(uint64(v)), nil
		case uint8:
			return humanize.Bytes(uint64(v)), nil
		case uint16:
			return humanize.Bytes(uint64(v)), nil
		case uint32:
			return humanize.Bytes(uint64(v)), nil
		case uint64:
			return humanize.Bytes(v), nil
		case float32:
			bi, _ := big.NewFloat(float64(v)).Int(nil)
			return humanize.BigBytes(bi), nil
		case float64:
			bi, _ := big.NewFloat(v).Int(nil)
			return humanize.BigBytes(bi), nil

		default:
			return "", fmt.Errorf("fileSize can't handle type %T (%#v)", v, v)
		}
	}

	switch val := args[0].(type) {
	case []any:
		r := make([]any, len(val))
		for i, o := range val {
			r[i], err = fileSize(o)
			if err != nil {
				return
			}
		}
		result = r
		return

	default:
		return fileSize(val)
	}
})

var humanTimeStringParseLayouts = []string{
	// in order, most common first:
	time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"

	time.RFC822,  // "02 Jan 06 15:04 MST"
	time.RFC822Z, // "02 Jan 06 15:04 -0700" // RFC822 with numeric zone

	time.RFC850, // "Monday, 02-Jan-06 15:04:05 MST"

	time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone

	time.Layout,   // "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
	time.ANSIC,    // "Mon Jan _2 15:04:05 2006"
	time.UnixDate, // "Mon Jan _2 15:04:05 MST 2006"
	time.RubyDate, // "Mon Jan 02 15:04:05 -0700 2006"

	time.DateTime, // "2006-01-02 15:04:05"
	time.DateOnly, // "2006-01-02"
	time.TimeOnly, // "15:04:05"

	time.Kitchen,    // "3:04PM"
	time.Stamp,      // "Jan _2 15:04:05"
	time.StampMilli, // "Jan _2 15:04:05.000"
	time.StampMicro, // "Jan _2 15:04:05.000000"
	time.StampNano,  // "Jan _2 15:04:05.000000000"
}

var humanTimeFunc = gval.Function("humanTime", func(args ...any) (result any, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("'humanTime' jsonpath function expects a single argument with the value or slice to be formatted. Got %#v instead", args)
	}

	asTime := func(val any) (time.Time, error) {
		switch v := val.(type) {
		case int:
			return time.UnixMilli(int64(v)), nil
		case int8:
			return time.UnixMilli(int64(v)), nil
		case int16:
			return time.UnixMilli(int64(v)), nil
		case int32:
			return time.UnixMilli(int64(v)), nil
		case int64:
			return time.UnixMilli(v), nil
		case uint:
			return time.UnixMilli(int64(v)), nil
		case uint8:
			return time.UnixMilli(int64(v)), nil
		case uint16:
			return time.UnixMilli(int64(v)), nil
		case uint32:
			return time.UnixMilli(int64(v)), nil
		case uint64:
			return time.UnixMilli(int64(v)), nil
		case float32:
			i, _ := big.NewFloat(float64(v)).Int64()
			return time.UnixMilli(i), nil
		case float64:
			i, _ := big.NewFloat(v).Int64()
			return time.UnixMilli(i), nil

		case string:
			for _, layout := range humanTimeStringParseLayouts {
				t, err := time.Parse(layout, v)
				if err == nil {
					return t, nil
				}
			}
			return time.Time{}, fmt.Errorf("humanTime can't parse string: %q", v)

		default:
			return time.Time{}, fmt.Errorf("humanTime can't handle type %T (%#v)", v, v)
		}
	}

	humanTime := func(val any) (string, error) {
		t, err := asTime(val)
		if err != nil {
			return "", err
		}
		return humanize.Time(t), nil
	}

	switch val := args[0].(type) {
	case []any:
		r := make([]any, len(val))
		for i, o := range val {
			r[i], err = humanTime(o)
			if err != nil {
				return
			}
		}
		result = r
		return

	default:
		return humanTime(val)
	}
})

var jsonPathBuilder = gval.Full(
	jsonpath.PlaceholderExtension(),
	hasKeyFunc,
	fileSizeFunc,
	humanTimeFunc,
)

func NewJsonPath(expression string) (jp gval.Evaluable, err error) {
	return jsonPathBuilder.NewEvaluable(expression)
}

func GetJsonPath(expression string, document any) (result any, err error) {
	jp, err := NewJsonPath(expression)
	if err != nil {
		return nil, err
	}
	return jp(context.Background(), document)
}

func CreateJsonPathChecker(expression string) (checker func(document any) (bool, error), err error) {
	jp, err := NewJsonPath(expression)
	if err != nil {
		return nil, err
	}
	return CreateJsonPathCheckerFromEvaluable(jp), nil
}

func CreateJsonPathCheckerFromEvaluable(jp gval.Evaluable) (checker func(document any) (bool, error)) {
	return func(value any) (ok bool, err error) {
		v, err := jp(context.Background(), value)
		if err != nil {
			return false, err
		}

		if v == nil {
			return false, nil
		} else if lst, ok := v.([]any); ok {
			return len(lst) > 0, nil
		} else if m, ok := v.(map[string]any); ok {
			return len(m) > 0, nil
		} else if b, ok := v.(bool); ok {
			return b, nil
		} else {
			return false, fmt.Errorf("unknown jsonpath result. Expected list, map or boolean. Got %#v", value)
		}
	}
}
