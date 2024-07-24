package provider

import (
	"reflect"
	"testing"
)

type TestConfig struct {
	Name    *string `json:"name"`
	Role    *string `json:"role"`
	Version *int    `json:"version"`
}

type TestConfig2 struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Version int    `json:"version"`
}

func mockFindKey(key string, value any) error {
	values := map[string]any{
		"name":    "TestApp",
		"role":    "TestRole",
		"version": 1,
		"env":     "TestEnv",
	}

	if val, ok := values[key]; ok {
		reflect.ValueOf(value).Elem().Set(reflect.ValueOf(val))
		return nil
	}
	return nil
}

func TestGetConfigsFromTags(t *testing.T) {
	config := TestConfig{}

	name := "TestApp"
	role := "TestRole"
	version := 1

	expected := TestConfig{
		Name:    &name,
		Role:    &role,
		Version: &version,
	}

	result := GetConfigsFromTags(mockFindKey, config)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestGetConfigsFromTags2(t *testing.T) {
	config := TestConfig2{}

	name := "TestApp"
	role := "TestRole"
	version := 1

	expected := TestConfig2{
		Name:    name,
		Role:    role,
		Version: version,
	}

	result := GetConfigsFromTags(mockFindKey, config)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestGetConfigsFromTagsIgnoreEmptyValues(t *testing.T) {
	type MyTestConfig struct {
		Banana  *string `json:"banana"`
		Role    *string `json:"role"`
		Version *int    `json:"version"`
	}

	config := MyTestConfig{}

	role := "TestRole"
	version := 1

	expected := MyTestConfig{
		Role:    &role,
		Version: &version,
	}

	result := GetConfigsFromTags(mockFindKey, config)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestSetField(t *testing.T) {
	type TestStruct struct {
		Name    string
		Role    string
		Version int
	}

	testStruct := TestStruct{
		Name:    "TestApp",
		Role:    "TestRole",
		Version: 1,
	}

	newValue := "NewName"

	err := setField(&testStruct, "Name", newValue)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if testStruct.Name != newValue {
		t.Errorf("Expected %s, got %s", newValue, testStruct.Name)
	}
}

func TestSetField2(t *testing.T) {
	type TestStruct struct {
		Name    *string
		Role    *string
		Version *int
	}

	testStruct := TestStruct{
		Name:    nil,
		Role:    nil,
		Version: nil,
	}

	newValue := "NewName"

	err := setField(&testStruct, "Name", newValue)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if *testStruct.Name != newValue {
		t.Errorf("Expected %s, got %s", newValue, *testStruct.Name)
	}
}
