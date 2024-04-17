package sdk

import (
	"os"
	"strings"
)

var version string = func() string {
	if vv := os.Getenv("VERSION"); vv != "" {
		if vv, ok := strings.CutPrefix(vv, "v"); ok {
			return vv
		}
		return vv
	}

	return "v0.0.0"
}()
