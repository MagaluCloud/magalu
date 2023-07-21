package mgc_openapi

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

func getExtensionBool(prefix *string, name string, extensions map[string]any, def bool) (b bool, ok bool) {
	value, _ := getExtension(prefix, name, extensions, def)
	b, ok = value.(bool)
	return
}

func getNameExtension(prefix *string, extensions map[string]any, def string) string {
	str, _ := getExtensionString(prefix, "name", extensions, def)
	return str
}

func getHiddenExtension(prefix *string, extensions map[string]any) bool {
	b, _ := getExtensionBool(prefix, "hidden", extensions, false)
	return b
}
