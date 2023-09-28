package utils

import "github.com/mitchellh/mapstructure"

// Attempts to decode the value passed in as input to the structure specified by the template
// Structs can be converted to maps and vice-versa, strings to integers, etc...
// The Tag used for specifying how to decode/encode is centralized in 'json', no need to use
// 'mapstructure'.
func DecodeValue[T any, U any](input T, output *U) (err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &output,
		TagName:          "json",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.RecursiveStructToMapHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
		),
	})
	if err != nil {
		return
	}
	err = decoder.Decode(input)
	return
}

// Attempts to decode the value passed in as input to the structure specified by the template
// Structs can be converted to maps and vice-versa, strings to integers, etc...
// The Tag used for specifying how to decode/encode is centralized in 'json', no need to use
// 'mapstructure'.
func DecodeNewValue[T any](input any) (result *T, err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &result,
		TagName:          "json",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.RecursiveStructToMapHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
		),
	})
	if err != nil {
		return
	}
	err = decoder.Decode(input)
	return
}
