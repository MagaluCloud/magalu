package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/spf13/afero"
)

type testCaseConfig struct {
	name       string
	run        func(c *Config) error
	expectedFs []utils.TestFsEntry
	providedFs []utils.TestFsEntry
}

func setupWithoutFile(path string) (*Config, afero.Fs) {
	pf, fs := profile_manager.NewInMemoryProfileManager()
	c := New(pf)
	c.init()

	return c, fs
}

func setupWithFile(testFileData []byte, path string) (*Config, error, afero.Fs) {
	m, fs := profile_manager.NewInMemoryProfileManager()
	if err := m.Current().Write(CONFIG_FILE, testFileData); err != nil {
		return nil, err, fs
	}

	c := New(m)

	return c, nil, fs
}

type unmarshalerField int

func (f *unmarshalerField) UnmarshalText(data []byte) error {
	str := string(data)
	if str == "valid" {
		*f = 100
		return nil
	}

	return fmt.Errorf("'unmarshalerField' only accepts 'valid' keyword. Got %q instead", str)
}

func TestGet(t *testing.T) {
	type person struct {
		Name          string `json:"name"`
		Age           int    `json:"age"`
		CaseSensitive string `json:"caseSensitive"`
	}

	type unmarshalerPerson struct {
		Name        string           `json:"name"`
		Age         int              `json:"age"`
		Unmarshaler unmarshalerField `json:"unmarshaler"`
	}

	type unmarshalerSubObject struct {
		Person unmarshalerPerson `json:"person"`
	}

	fsPath := "/"
	t.Run("decode to no pointer", func(t *testing.T) {
		c, err, _ := setupWithFile([]byte(`{ "foo": "bar" }`), fsPath)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		var p person
		err = c.Get("foo", p)

		if err == nil {
			t.Errorf("expected err != nil, found: %#v", err)
		}
	})

	t.Run("decode to nil pointer", func(t *testing.T) {
		c, err, _ := setupWithFile([]byte(`{ "foo": "bar" }`), fsPath)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		var p person
		err = c.Get("foo", p)

		if err == nil {
			t.Errorf("expected err != nil, found: %#v", err)
		}
	})

	t.Run("decode partial to non-zero struct", func(t *testing.T) {
		data := `{"person":{"name":"Josh"}}`
		c, err, _ := setupWithFile([]byte(data), fsPath)
		if err != nil {
			t.Errorf("setting up file expected err == nil, found: %#v", err)
		}

		p := person{Age: 20}
		err = c.Get("person", &p)
		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		expected := person{Age: 20, Name: "Josh"}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected %#v when decoding %s, found %#v instead", expected, data, p)
		}
	})

	t.Run("decode from config file to pointer", func(t *testing.T) {
		data := `{
			"foo": {
				"name": "jon",
				"age": 5,
				"caseSensitive": "some"
			}
		}`
		expected := person{Name: "jon", Age: 5, CaseSensitive: "some"}

		c, err, _ := setupWithFile([]byte(data), fsPath)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		p := new(person)
		err = c.Get("foo", p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if !reflect.DeepEqual(*p, expected) {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})

	t.Run("decode env var to pointer", func(t *testing.T) {
		data := `{
			"name": "jon",
			"age": 5
		}`
		// TODO: The data above should have `"caseSensitive": "some"`, but we currently
		// don't have support for case sensitive env vars...
		var expected person
		_ = json.Unmarshal([]byte(data), &expected)

		t.Setenv("MGC_FOO", data)

		c, _ := setupWithoutFile(fsPath)

		p := new(person)
		err := c.Get("foo", p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if !reflect.DeepEqual(*p, expected) {
			t.Errorf("expected p == %#v, found: %#v", expected, *p)
		}
	})

	t.Run("decode string in config file to string", func(t *testing.T) {
		data := `{ "foo": "bar" }`
		expected := "bar"

		c, err, _ := setupWithFile([]byte(data), fsPath)
		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		var p string
		err = c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if p != expected {
			t.Errorf("expected p == bar, found: %#v", p)
		}
	})

	t.Run("decode string in env var to string", func(t *testing.T) {
		expected := "bar"
		t.Setenv("MGC_FOO", expected)

		c, _ := setupWithoutFile(fsPath)

		var p string
		err := c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if p != expected {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})

	t.Run("decode object in config file to struct", func(t *testing.T) {
		data := `{
			"foo": {
				"name": "jon",
				"age": 5,
				"caseSensitive": "some"
			}
		}`
		expected := person{Name: "jon", Age: 5, CaseSensitive: "some"}

		c, err, _ := setupWithFile([]byte(data), fsPath)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		var p person
		err = c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})

	t.Run("decode object in env var to struct", func(t *testing.T) {
		data := `{
			"name": "jon",
			"age": 5
		}`
		// TODO: The data above should have `"caseSensitive": "some"`, but we currently
		// don't have support for case sensitive env vars...
		var expected person
		_ = json.Unmarshal([]byte(data), &expected)

		t.Setenv("MGC_FOO", data)

		c, _ := setupWithoutFile(fsPath)

		var p person
		err := c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})

	t.Run("decode string in config file to any", func(t *testing.T) {
		data := `{ "foo": "bar" }`
		expected := "bar"

		c, err, _ := setupWithFile([]byte(data), fsPath)
		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		var p any
		err = c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if p != expected {
			t.Errorf("expected p == bar, found: %#v", p)
		}
	})

	t.Run("decode object in config file with unmarshaler types", func(t *testing.T) {
		// We save objects in Config File as strings...
		data := `{
			"foo": "{\"name\":\"jon\",\"age\":5,\"unmarshaler\":\"valid\"}"
		}`

		expected := unmarshalerPerson{
			Name:        "jon",
			Age:         5,
			Unmarshaler: 100,
		}

		c, err, _ := setupWithFile([]byte(data), fsPath)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		p := new(unmarshalerPerson)
		err = c.Get("foo", p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if !reflect.DeepEqual(*p, expected) {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})

	t.Run("decode object in config file with subfield unmarshaler types", func(t *testing.T) {
		// We save objects in Config File as strings...
		data := `{
			"foo": "{\"person\":{\"name\":\"jon\",\"age\":5,\"unmarshaler\":\"valid\"}}"
		}`

		expected := unmarshalerSubObject{
			Person: unmarshalerPerson{
				Name:        "jon",
				Age:         5,
				Unmarshaler: 100,
			},
		}

		c, err, _ := setupWithFile([]byte(data), fsPath)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		p := new(unmarshalerSubObject)
		err = c.Get("foo", p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if !reflect.DeepEqual(*p, expected) {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})
	t.Run("decode object in config file to any", func(t *testing.T) {
		// We save objects in Config File as strings...
		data := `{
			"foo": "{\"name\":\"jon\",\"age\":5,\"caseSensitive\":\"some\"}"
		}`

		expected := map[string]any{
			"name":          "jon",
			"age":           float64(5),
			"caseSensitive": "some",
		}

		c, err, _ := setupWithFile([]byte(data), fsPath)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}

		var p any
		err = c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})

	t.Run("decode string in env var to any", func(t *testing.T) {
		expected := "bar"

		t.Setenv("MGC_FOO", expected)

		c, _ := setupWithoutFile(fsPath)

		var p any
		err := c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if p != expected {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})

	t.Run("decode object in env var to any", func(t *testing.T) {
		data := `{
			"name": "jon",
			"age": 5,
			"caseSensitive": "some"
		}`

		t.Setenv("MGC_FOO", data)

		expected := map[string]any{
			"name":          "jon",
			"age":           float64(5),
			"caseSensitive": "some",
		}

		c, _ := setupWithoutFile(fsPath)

		var p any
		err := c.Get("foo", &p)

		if err != nil {
			t.Errorf("expected err == nil, found: %#v", err)
		}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("expected p == %#v, found: %#v", expected, p)
		}
	})
}

func deleteTest(name string, key string, expected any, provided []utils.TestFsEntry, expectedfs []utils.TestFsEntry) testCaseConfig {
	provided = utils.AutoMkdirAll(provided)
	expectedfs = utils.AutoMkdirAll(expectedfs)
	return testCaseConfig{
		name:       fmt.Sprintf("Config.DeleteWithFile(%q)", name),
		providedFs: provided,
		expectedFs: expectedfs,
		run: func(c *Config) error {
			return c.Delete(key)
		},
	}
}

func getTest(name string, key string, expected any, provided []utils.TestFsEntry, expectedfs []utils.TestFsEntry) testCaseConfig {
	provided = utils.AutoMkdirAll(provided)
	expectedfs = utils.AutoMkdirAll(expectedfs)
	return testCaseConfig{
		name:       fmt.Sprintf("Config.Get(%q)", name),
		providedFs: provided,
		expectedFs: expectedfs,
		run: func(c *Config) error {
			var out any
			if err := c.Get(key, &out); err != nil {
				return fmt.Errorf("expected err == nil, found: %#v", err)
			}

			if out != expected {
				return fmt.Errorf("expected %#v, found %#v", expected, out)
			}
			return nil
		},
	}
}

func setTest(name string, key string, expected any, provided []utils.TestFsEntry, expectedfs []utils.TestFsEntry) testCaseConfig {
	provided = utils.AutoMkdirAll(provided)
	expectedfs = utils.AutoMkdirAll(expectedfs)
	return testCaseConfig{
		name:       fmt.Sprintf("Config.SetWithFile(%q)", name),
		providedFs: provided,
		expectedFs: expectedfs,
		run: func(c *Config) error {
			return c.Set(key, expected)

		},
	}
}

func TestConfigManagerWithFile(t *testing.T) {
	tests := []testCaseConfig{
		deleteTest("test1", "foo", nil,
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo: bar`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`{}
`),
				},
			}),
		deleteTest("test2", "foo", nil,
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo:`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo:`),
				},
			}),
		deleteTest("test3", "foo", nil,
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte{},
				},
			}),
		deleteTest("withoutFile3", "foo", nil,
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`""`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`""`),
				},
			}),

		getTest("test1", "foo", "bar",
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo: bar`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo: bar`),
				},
			}),
		getTest("test2", "foo", nil,
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo`),
				},
			}),
		getTest("test3", "foo", nil,
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`""`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`""`),
				},
			}),
		setTest("test1", "foo", "woo",
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo: bar`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo: woo
`),
				},
			}),

		setTest("test2", "foo", "woo",
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo:`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo: woo
`),
				},
			}),
		setTest("test3", "foo", "woo",
			[]utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`""`),
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo: woo
`),
				},
			}),
		setTest("test1", "foo", "woo",
			[]utils.TestFsEntry{
				{
					Path: "/default/",
					Mode: utils.DIR_PERMISSION,
					Data: []byte{},
				},
			}, []utils.TestFsEntry{
				{
					Path: "/default/cli.yaml",
					Mode: utils.FILE_PERMISSION,
					Data: []byte(`foo: woo
`)},
			}),
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			m, fs := profile_manager.NewInMemoryProfileManager()
			fs_err := utils.PrepareFs(fs, tc.providedFs)

			if fs_err != nil {
				t.Errorf("could not prepare provided FS: %s", fs_err.Error())
			}
			c := New(m)
			run_error := tc.run(c)

			if run_error != nil {
				t.Errorf("expected err == nil, found: %v", run_error)
			}

			fs_err = utils.CheckFs(fs, tc.expectedFs)

			if fs_err != nil {
				t.Errorf("unexpected FS state: %s", fs_err.Error())
			}
		})
	}
}
