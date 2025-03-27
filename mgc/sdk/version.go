package sdk

import (
	_ "embed"
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

var RawVersion string

var Version string = func() string {
	if RawVersion == "" {
		return getVCSInfo("v0.0.0")
	}

	// Validate version format using regex
	matched, err := regexp.MatchString(`^v\d+\.\d+\.\d+$`, RawVersion)
	if err != nil || !matched {
		return getVCSInfo(RawVersion)
	}

	return strings.Trim(RawVersion, " \t\n\r")
}()

func getVCSInfo(version string) string {
	if info, ok := debug.ReadBuildInfo(); ok {
		var vcs, rev, status string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs":
				vcs = setting.Value
			case "vcs.revision":
				rev = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					status = " (modified)"
				}
			}
		}

		if vcs != "" {
			return fmt.Sprintf("%s-%s-%s%s", version, vcs, rev, status)
		}
	}
	return "v0.0.0 dev"
}
