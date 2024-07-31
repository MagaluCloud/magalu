//go:build !release

package sdk

// in order to use this, build with -buildvcs=true
import (
	"fmt"
	"runtime/debug"
)

var Version = func() string {
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
	return version + "dev"
}()
