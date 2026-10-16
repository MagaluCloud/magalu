package posthog

import "fmt"

type debugEntry struct {
	msg string
	kv  []any
}

func (d debugEntry) value(key string) string {
	for i := 0; i+1 < len(d.kv); i += 2 {
		if d.kv[i] == key {
			return fmt.Sprint(d.kv[i+1])
		}
	}
	return ""
}

func envFrom(vars map[string]string) func(string) string {
	return func(k string) string { return vars[k] }
}
