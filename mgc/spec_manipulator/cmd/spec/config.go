package spec

import (
	"fmt"
	"slices"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

type removeItem struct {
	Path       *string  `json:"path,omitempty" yaml:"path,omitempty"`
	PathPrefix *string  `json:"path_prefix,omitempty" yaml:"path_prefix,omitempty"`
	Method     []string `json:"method,omitempty" yaml:"method,omitempty"`
}

type specList struct {
	Url      string       `json:"url" yaml:"url"`
	File     string       `json:"file" yaml:"file"`
	Menu     string       `json:"menu" yaml:"menu"`
	Enabled  bool         `json:"enabled" yaml:"enabled"`
	ToRemove []removeItem `json:"to_remove,omitempty" yaml:"to_remove,omitempty"`
}

func interfaceToMap(i interface{}) (map[string]interface{}, bool) {
	mapa, ok := i.(map[string]interface{})
	if !ok {
		fmt.Println("A interface não é um mapa ou mapa de interfaces.")
		return nil, false
	}
	return mapa, true
}

func saveConfig(config []specList) {
	viper.Set("jaxyendy", config)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func loadList(specificMenu string) ([]specList, error) {
	var currentConfig []specList
	config := viper.Get("jaxyendy")

	if config != nil {
		for _, v := range config.([]interface{}) {
			vv, ok := interfaceToMap(v)
			if !ok {
				return currentConfig, fmt.Errorf("fail to load current config")
			}
			if specificMenu != "" && vv["menu"].(string) != specificMenu {
				continue
			}

			toRemove := []removeItem{}
			if vv["to_remove"] != nil {

				for _, g := range vv["to_remove"].([]interface{}) {
					gg, ok := interfaceToMap(g)
					if !ok {
						return currentConfig, fmt.Errorf("fail to load current config")
					}

					var path *string
					var pathPrefix *string
					var method []string

					if gg["path"] != nil {
						path = new(string)
						*path = gg["path"].(string)
					}
					if gg["method"] != nil {
						for _, g := range gg["method"].([]interface{}) {
							method = append(method, g.(string))
						}
					}
					if gg["path_prefix"] != nil {
						pathPrefix = new(string)
						*pathPrefix = gg["path_prefix"].(string)
					}

					toRemove = append(toRemove, removeItem{
						Path:       path,
						Method:     method,
						PathPrefix: pathPrefix,
					})
				}
			}
			currentConfig = append(currentConfig, specList{
				Url:      vv["url"].(string),
				Menu:     vv["menu"].(string),
				File:     vv["file"].(string),
				Enabled:  vv["enabled"].(bool),
				ToRemove: toRemove,
			})

		}

	}
	return currentConfig, nil
}

func loadListMap() ([]string, []specList, error) {
	currentConfig, err := loadList("")
	if err != nil {
		return nil, nil, err
	}

	menus := []string{}
	for _, v := range currentConfig {
		menus = append(menus, v.Menu)
	}

	slices.Sort(menus)

	return menus, currentConfig, nil
}

func getConfigToRun() ([]specList, error) {
	menus, currentConfig, err := loadListMap()
	if err != nil {
		return nil, err
	}

	ms := pterm.DefaultInteractiveMultiselect.
		WithDefaultText("Select products").
		WithMaxHeight(14).
		WithOptions(menus)

	op, err := ms.Show()
	if err != nil {
		return nil, err
	}

	if len(op) == 0 {
		return nil, fmt.Errorf("no products selected")
	}

	configToRun := []specList{}
	for _, v := range op {
		for _, v2 := range currentConfig {
			if v2.Menu == v {
				configToRun = append(configToRun, v2)
			}
		}
	}
	pterm.Println()
	return configToRun, nil
}
