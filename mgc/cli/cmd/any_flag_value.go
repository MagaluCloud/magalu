package cmd

import (
	"encoding/json"
	"os"
	"strings"
)

type AnyFlagValue struct {
	value    any
	typeName string
}

func (f *AnyFlagValue) String() string {
	str, err := json.Marshal(f.value)
	if err != nil {
		return ""
	}
	return string(str)
}

func (f *AnyFlagValue) Set(val string) error {
	err := json.Unmarshal([]byte(val), &f.value)
	if err != nil {
		if !strings.HasPrefix(val, "@") {
			f.value = val
			return nil
		} else if f.value, err = loadFromFile(val[1:]); err != nil {
			return err
		}
	}
	return err
}

func (f *AnyFlagValue) Type() string {
	return f.typeName
}

func (f *AnyFlagValue) Value() any {
	return f.value
}

func loadFromFile(filename string) (value any, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// TODO: eventually we want to load raw... could we determine from schema?
	err = json.Unmarshal(data, &value)
	if err != nil {
		value = string(data)
	}
	return value, nil
}
