package sdk

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

//go:embed version.txt
var rawVersion string

var version string = func() string {
	if vv := os.Getenv("VERSION"); vv != "" {
		return fmt.Sprintf("v%s", strings.TrimPrefix(vv, "v"))
	}
	return strings.Trim(rawVersion, " \t\n\r")
}()
