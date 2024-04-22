package sdk

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var rawVersion string

var version string = func() string {
	return strings.Trim(rawVersion, " \t\n\r")
}()
