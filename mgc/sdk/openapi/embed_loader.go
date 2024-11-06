package openapi

import (
	"embed"
	"fmt"
	"os"
	"syscall"

	"gopkg.in/yaml.v3"
	"magalu.cloud/core/dataloader"
)

//go:embed openapis/*.yaml
var folder embed.FS

type embedLoader map[string][]byte

func GetEmbedLoader() dataloader.Loader {
	return embedLoaderInstance()
}

func (f embedLoader) Load(name string) ([]byte, error) {
	if data, ok := embedLoaderInstance()[name]; ok {
		return data, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: syscall.ENOENT}
}

func (f embedLoader) String() string {
	return "embedLoader"
}

type module struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Summary     string `yaml:"summary"`
	URL         string `yaml:"url"`
	Version     string `yaml:"version"`
}

type configIndex struct {
	Modules []module `yaml:"modules"`
	Version string   `yaml:"version"`
}

var embedLoaderInstance = func() embedLoader {
	loader := embedLoader{}
	indexFileName := "index.openapi.yaml"
	dataIndex, err := folder.ReadFile("openapis/" + indexFileName)
	if err != nil {
		fmt.Println("Error reading index file")
		panic(err)
	}

	config := &configIndex{}
	err = yaml.Unmarshal(dataIndex, config)
	if err != nil {
		fmt.Println("Error unmarshalling index file")
		panic(err)
	}

	loader[indexFileName] = dataIndex

	for _, m := range config.Modules {
		data, err := folder.ReadFile("openapis/" + m.Path)
		if err != nil {
			fmt.Println("Error reading module file")
			panic(err)
		}
		loader[m.Path] = data

	}
	return loader
}
