package openapi

func getExtension(prefix *string, name string, extensions map[string]any, def any) (value any, ok bool) {
	if prefix == nil || *prefix == "" {
		return def, false
	}
	key := *prefix + "-" + name
	value, ok = extensions[key]
	if !ok {
		value = def
	}
	return
}

func getExtensionString(prefix *string, name string, extensions map[string]any, def string) (str string, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	str, ok = value.(string)
	return
}

func getExtensionObject(prefix *string, name string, extensions map[string]any, def map[string]any) (m map[string]any, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	m, ok = value.(map[string]any)
	return
}

func getExtensionArray(prefix *string, name string, extensions map[string]any, def []any) (m []any, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	m, ok = value.([]any)
	return
}

func getExtensionBool(prefix *string, name string, extensions map[string]any, def bool) (b bool, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	b, ok = value.(bool)
	return
}

func getNameExtension(prefix *string, extensions map[string]any, def string) string {
	str, _ := getExtensionString(prefix, "name", extensions, def)
	return str
}

func getDescriptionExtension(prefix *string, extensions map[string]any, def string) string {
	str, _ := getExtensionString(prefix, "description", extensions, def)
	return str
}

func getHiddenExtension(prefix *string, extensions map[string]any) bool {
	b, _ := getExtensionBool(prefix, "hidden", extensions, false)
	return b
}

func getRemoveFromCLIFlag(prefix *string, extensions map[string]any) bool {
	b, _ := getExtensionBool(prefix, "remove-from-cli", extensions, false)
	return b
}
