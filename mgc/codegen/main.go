package main

import (
	"fmt"
	"os"

	"github.com/MagaluCloud/magalu/mgc/codegen/cmd"
)

func main() {
	rootCmd := cmd.NewRoot()
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
