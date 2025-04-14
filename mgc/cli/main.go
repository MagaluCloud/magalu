package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
)

var RawVersion string

var Version string = func() string {
	if RawVersion == "" {
		return getVCSInfo("v0.0.0")
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
			return fmt.Sprintf("%s %s%s", version, rev, status)
		}
	}
	return "v0.0.0 dev"
}

func main() {
	defer panicRecover()
	mgcSdk.SetUserAgent("MgcCLI")

	err := cmd.Execute(Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func panicRecover() {
	err := recover()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n😔 Oops! Something went wrong.\nPlease help us improve by sending the error report to: https://help.github.com/MagaluCloud/magalu/mgc/hc/pt-br/requests/new\n\n  Version: %s\n  SO: %s / %s\n  Args: %s\n  Error: %s\n\nThank you for your cooperation!\n\n",
			Version,
			runtime.GOOS,
			runtime.GOARCH,
			strings.Join(os.Args, " "),
			err)
		os.Exit(1)
	}
}
