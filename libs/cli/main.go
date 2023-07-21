package main

import (
	"fmt"

	"github.com/profusion/magalu/libs/cli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Error running the CLI: %s\n", err)
	}
}
