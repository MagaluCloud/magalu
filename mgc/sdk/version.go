package sdk

import (
	_ "embed"
	"strings"
)

var RawVersion string

var version string = func() string {
	return strings.Trim(RawVersion, " \t\n\r")
}()
