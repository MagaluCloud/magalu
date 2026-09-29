package preview

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type testIndex struct {
	Version string `yaml:"version"`
	Modules []struct {
		Name        string `yaml:"name"`
		Url         string `yaml:"url"`
		Path        string `yaml:"path"`
		Version     string `yaml:"version"`
		Description string `yaml:"description"`
		Summary     string `yaml:"summary"`
	} `yaml:"modules"`
}

func TestStoreLoader(t *testing.T) {
	t.Run("synthesizes the index from the store names", func(t *testing.T) {
		store := newMemStore(map[string]string{
			"foo":    fooSpec,
			"bar":    "openapi: 3.0.3\ninfo:\n  title: Bar\n",
			"broken": brokenSpec,
		})
		loader := &storeLoader{level: Alpha, store: store, validator: &fakeValidator{}}

		data, err := loader.Load(indexFileName)
		require.NoError(t, err)

		var index testIndex
		require.NoError(t, yaml.Unmarshal(data, &index))
		assert.Equal(t, "1.0.0", index.Version)
		require.Len(t, index.Modules, 3)

		bar, broken, foo := index.Modules[0], index.Modules[1], index.Modules[2]

		assert.Equal(t, "foo", foo.Name)
		assert.Equal(t, "foo.openapi.yaml", foo.Path)
		assert.Equal(t, "alpha://foo", foo.Url)
		assert.Equal(t, "1.0", foo.Version)
		assert.Equal(t, "Foo preview API.", foo.Description)
		assert.Equal(t, "Foo preview API.", foo.Summary)

		// falls back to the title when there is no description
		assert.Equal(t, "bar", bar.Name)
		assert.Equal(t, "Bar", bar.Description)

		// a spec that does not parse still gets listed, so the error shows up
		// only when that module is used
		assert.Equal(t, "broken", broken.Name)
		assert.Equal(t, "alpha broken", broken.Description)
	})

	t.Run("loads validated specs", func(t *testing.T) {
		validator := &fakeValidator{}
		loader := &storeLoader{level: Beta, store: newMemStore(map[string]string{"foo": fooSpec}), validator: validator}

		data, err := loader.Load("foo.openapi.yaml")
		require.NoError(t, err)
		assert.Equal(t, fooSpec, string(data))
		assert.Equal(t, []string{"beta/foo"}, validator.seen)
	})

	t.Run("refuses specs the validator refuses", func(t *testing.T) {
		validator := &fakeValidator{refused: map[string]string{"foo": "nope"}}
		loader := &storeLoader{level: Beta, store: newMemStore(map[string]string{"foo": fooSpec}), validator: validator}

		_, err := loader.Load("foo.openapi.yaml")
		assert.EqualError(t, err, "nope")
	})

	t.Run("refuses files that are not module specs", func(t *testing.T) {
		loader := &storeLoader{level: Beta, store: newMemStore(map[string]string{"foo": fooSpec}), validator: &fakeValidator{}}

		_, err := loader.Load("foo.json")
		assert.Error(t, err)
	})
}
