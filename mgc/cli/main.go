package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"magalu.cloud/cli/cmd"
	mgcSdk "magalu.cloud/sdk"
)

func main() {
	defer panicRecover()
	mgcSdk.SetUserAgent("MgcCLI")

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	debug.SetGCPercent(10)
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func panicRecover() {
	err := recover()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n😔 Oops! Something went wrong.\nPlease help us improve by sending the error report to: https://help.magalu.cloud/hc/pt-br/requests/new\n\n  Version: %s\n  SO: %s / %s\n  Args: %s\n  Error: %s\n\nThank you for your cooperation!\n\n",
			mgcSdk.Version,
			runtime.GOOS,
			runtime.GOARCH,
			strings.Join(os.Args, " "),
			err)
		os.Exit(1)
	}
}
