package preview

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core/dataloader"
	"gopkg.in/yaml.v3"
)

const (
	indexFileName = "index.openapi.yaml"
	indexVersion  = "1.0.0"
	specSuffix    = ".openapi.yaml"
)

// storeLoader feeds a Store to the openapi source. It synthesizes the index
// from the store names, so every package "<name>" becomes the module "<name>".
type storeLoader struct {
	level     Level
	store     Store
	validator Validator
}

func (l *storeLoader) Load(file string) ([]byte, error) {
	if file == indexFileName {
		return l.buildIndex()
	}

	name, ok := strings.CutSuffix(file, specSuffix)
	if !ok {
		return nil, fmt.Errorf("%s spec %q must be named <module>%s", l.level.Name, file, specSuffix)
	}

	pkg, err := l.store.Read(name)
	if err != nil {
		return nil, err
	}
	if err := l.validator.Validate(l.level, pkg); err != nil {
		return nil, err
	}
	return pkg.Spec, nil
}

type indexModule struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Path        string `json:"path"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type index struct {
	Version string        `json:"version"`
	Modules []indexModule `json:"modules"`
}

// buildIndex is written as JSON, which is also valid YAML for the openapi source.
func (l *storeLoader) buildIndex() ([]byte, error) {
	names, err := l.store.Names()
	if err != nil {
		return nil, err
	}

	result := index{Version: indexVersion, Modules: make([]indexModule, 0, len(names))}
	for _, name := range names {
		info := l.readInfo(name)

		description := info.Description
		if description == "" {
			description = info.Title
		}
		if description == "" {
			description = fmt.Sprintf("%s %s", l.level.Name, name)
		}

		result.Modules = append(result.Modules, indexModule{
			Name:        name,
			Url:         l.level.Name + "://" + name,
			Path:        name + specSuffix,
			Version:     info.Version,
			Description: description,
			Summary:     description,
		})
	}
	return json.Marshal(result)
}

type specInfo struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// readInfo is best effort: a package that does not read or parse is still
// listed, and fails only when its module is used.
func (l *storeLoader) readInfo(name string) specInfo {
	var spec struct {
		Info specInfo `yaml:"info"`
	}
	if pkg, err := l.store.Read(name); err == nil {
		_ = yaml.Unmarshal(pkg.Spec, &spec)
	}
	return spec.Info
}

var _ dataloader.Loader = (*storeLoader)(nil)
