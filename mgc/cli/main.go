package main

import (
	"fmt"
	"net/url"
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

		Url := "https://github.com/MagaluCloud/mgccli/issues/new"
		args := strings.Join(os.Args, " ")

		query := url.Values{}
		query.Add("title", fmt.Sprintf("Error report at '%s'", args))
		query.Add("body", fmt.Sprintf("Version: %s\nSO: %s / %s\nArgs: %s\nError: %s\n",
			Version,
			runtime.GOOS,
			runtime.GOARCH,
			args,
			err))
		Url = Url + "?" + query.Encode()

		fmt.Fprintf(os.Stderr, `
😔 Oops! Something went wrong.
     Version: %s
     SO: %s / %s  
     Args: %s 
     Error: %s

Please help us improve by sending the error report to our repository:
	%s

Thank you for your cooperation!
`,
			Version,
			runtime.GOOS,
			runtime.GOARCH,
			args,
			err,
			Url)
		os.Exit(1)
	}
}
