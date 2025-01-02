package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
)

func main() {
	defer panicRecover()
	mgcSdk.SetUserAgent("MgcCLI")

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func panicRecover() {
	err := recover()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n😔 Oops! Something went wrong.\nPlease help us improve by sending the error report to: https://help.github.com/MagaluCloud/magalu/mgc/hc/pt-br/requests/new\n\n  Version: %s\n  SO: %s / %s\n  Args: %s\n  Error: %s\n\nThank you for your cooperation!\n\n",
			mgcSdk.Version,
			runtime.GOOS,
			runtime.GOARCH,
			strings.Join(os.Args, " "),
			err)
		os.Exit(1)
	}
}
