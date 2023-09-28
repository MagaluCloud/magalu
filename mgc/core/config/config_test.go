package config

import (
	"encoding/json"
	"path"
	"reflect"
	"testing"

	"github.com/spf13/afero"
	"magalu.cloud/core/utils"
)

type test struct {
	key      string
	fileData []byte
	expected any
}

func setupWithoutFile() *Config {
	path, _ := utils.BuildMGCPath()
	c := New()
	c.init(path, afero.NewMemMapFs())

	return c
}

func setupWithFile(testFileData []byte) (*Config, error) {
	file, err := utils.BuildMGCFilePath(CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, file, testFileData, 0644); err != nil {
		return nil, err
	}

	c := New()
	c.init(path.Dir(file), fs)

	return c, nil
}

func TestGetWithoutFile(t *testing.T) {
	tests := []test{
		{key: "foo", fileData: []byte{}, expected: nil},
	}

	for _, tc := range tests {
		c := setupWithoutFile()

		var out any
		if err := c.Get(tc.key, &out); err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		if out != tc.expected {
			t.Errorf("expected %v, found %v", tc.expected, out)
		}
	}
}

func TestGetWithFile(t *testing.T) {
	tests := []test{
		{key: "foo", fileData: []byte(`foo: bar`), expected: "bar"},
		{key: "foo", fileData: []byte(`foo:`), expected: nil},
		{key: "foo", fileData: []byte(``), expected: nil},
	}

	for _, tc := range tests {
		c, err := setupWithFile(tc.fileData)
		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		var out any
		if err := c.Get(tc.key, &out); err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		if out != tc.expected {
			t.Errorf("expected %v, found %v", tc.expected, out)
		}
	}
}

func TestGet(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("decode to no pointer", func(t *testing.T) {
		c, err := setupWithFile([]byte(`{ "foo": "bar" }`))

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		var p person
		err = c.Get("foo", p)

		if err == nil {
			t.Errorf("expected err != nil, found: %v", err)
		}
	})

	t.Run("decode to nil pointer", func(t *testing.T) {
		c, err := setupWithFile([]byte(`{ "foo": "bar" }`))

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		var p person
		err = c.Get("foo", p)

		if err == nil {
			t.Errorf("expected err != nil, found: %v", err)
		}
	})

	t.Run("decode from config file to pointer", func(t *testing.T) {
		data := `{
			"foo": {
				"name": "jon",
				"age": 5
			}
		}`
		expected := person{Name: "jon", Age: 5}

		c, err := setupWithFile([]byte(data))

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		p := new(person)
		err = c.Get("foo", p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if !reflect.DeepEqual(*p, expected) {
			t.Errorf("expected p == %+v, found: %+v", expected, p)
		}
	})

	t.Run("decode env var to pointer", func(t *testing.T) {
		data := `{
			"name": "jon",
			"age": 5
		}`
		var expected person
		_ = json.Unmarshal([]byte(data), &expected)

		t.Setenv("MGC_FOO", data)

		c := setupWithoutFile()

		p := new(person)
		err := c.Get("foo", p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if !reflect.DeepEqual(*p, expected) {
			t.Errorf("expected p == %v, found: %v", expected, p)
		}
	})

	t.Run("decode string in config file to string", func(t *testing.T) {
		data := `{ "foo": "bar" }`
		expected := "bar"

		c, err := setupWithFile([]byte(data))
		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		var p string
		err = c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if p != expected {
			t.Errorf("expected p == bar, found: %+v", p)
		}
	})

	t.Run("decode string in env var to string", func(t *testing.T) {
		expected := "bar"
		t.Setenv("MGC_FOO", expected)

		c := setupWithoutFile()

		var p string
		err := c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if p != expected {
			t.Errorf("expected p == %v, found: %v", expected, p)
		}
	})

	t.Run("decode object in config file to struct", func(t *testing.T) {
		data := `{
			"foo": {
				"name": "jon",
				"age": 5
			}
		}`
		expected := person{Name: "jon", Age: 5}

		c, err := setupWithFile([]byte(data))

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		var p person
		err = c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected p == %+v, found: %+v", expected, p)
		}
	})

	t.Run("decode object in env var to struct", func(t *testing.T) {
		data := `{
			"name": "jon",
			"age": 5
		}`
		var expected person
		_ = json.Unmarshal([]byte(data), &expected)

		t.Setenv("MGC_FOO", data)

		c := setupWithoutFile()

		var p person
		err := c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected p == %+v, found: %+v", expected, p)
		}
	})

	t.Run("decode string in config file to any", func(t *testing.T) {
		data := `{ "foo": "bar" }`
		expected := "bar"

		c, err := setupWithFile([]byte(data))
		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		var p any
		err = c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if p != expected {
			t.Errorf("expected p == bar, found: %+v", p)
		}
	})

	t.Run("decode object in config file to any", func(t *testing.T) {
		data := `{
			"foo": {
				"name": "jon",
				"age": 5
			}
		}`

		expected := map[string]any{
			"name": "jon",
			"age":  5,
		}

		c, err := setupWithFile([]byte(data))

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		var p any
		err = c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected p == %+v, found: %+v", expected, p)
		}
	})

	t.Run("decode string in env var to any", func(t *testing.T) {
		expected := "bar"

		t.Setenv("MGC_FOO", expected)

		c := setupWithoutFile()

		var p any
		err := c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if p != expected {
			t.Errorf("expected p == %+v, found: %+v", expected, p)
		}
	})

	t.Run("decode object in env var to any", func(t *testing.T) {
		expected := `{
			"name": "jon",
			"age": 5
		}`

		t.Setenv("MGC_FOO", expected)

		c := setupWithoutFile()

		var p any
		err := c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected p == %+v, found: %+v", expected, p)
		}
	})
}

func TestSetWithoutFile(t *testing.T) {
	tests := []test{
		{key: "foo", fileData: []byte{}, expected: "woo"},
	}

	for _, tc := range tests {
		c := setupWithoutFile()

		if err := c.Set(tc.key, tc.expected); err != nil {
			t.Errorf("expected err == nil , found %v", err)
		}

		var v any
		if err := c.Get(tc.key, &v); err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if v != tc.expected {
			t.Errorf("expected %v, found %v", tc.expected, v)
		}
	}
}

func TestSetWithFile(t *testing.T) {
	tests := []test{
		{key: "foo", fileData: []byte("foo: bar"), expected: "woo"},
		{key: "foo", fileData: []byte("foo:"), expected: "woo"},
		{key: "foo", fileData: []byte(""), expected: "woo"},
	}

	for _, tc := range tests {
		c, err := setupWithFile(tc.fileData)
		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		if err := c.Set(tc.key, tc.expected); err != nil {
			t.Errorf("expected err == nil , found %v", err)
		}

		var v any
		if err := c.Get(tc.key, &v); err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if v != tc.expected {
			t.Errorf("expected %v, found %v", tc.expected, v)
		}
	}
}

func TestDeleteWithoutFile(t *testing.T) {
	tests := []test{
		{key: "foo", fileData: []byte("foo: bar"), expected: nil},
		{key: "foo", fileData: []byte("foo:"), expected: nil},
		{key: "foo", fileData: []byte(""), expected: nil},
	}

	for _, tc := range tests {
		c := setupWithoutFile()

		if err := c.Delete(tc.key); err != nil {
			t.Errorf("expected err == nil, found %v", err)
		}

		var v any
		if err := c.Get(tc.key, &v); err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if v != tc.expected {
			t.Errorf("expected %v, found %v", tc.expected, v)
		}
	}
}

func TestDeleteWithFile(t *testing.T) {
	tests := []test{
		{key: "foo", fileData: []byte("foo: bar"), expected: nil},
		{key: "foo", fileData: []byte("foo:"), expected: nil},
		{key: "foo", fileData: []byte(""), expected: nil},
	}

	for _, tc := range tests {
		c, err := setupWithFile(tc.fileData)
		if err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}

		if err := c.Delete(tc.key); err != nil {
			t.Errorf("expected err == nil, found %v", err)
		}

		var v any
		if err := c.Get(tc.key, &v); err != nil {
			t.Errorf("expected err == nil, found: %v", err)
		}
		if v != tc.expected {
			t.Errorf("expected %v, found %v", tc.expected, v)
		}
	}
}
