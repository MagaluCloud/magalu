package main

import (
	"fmt"
	"load-specs/cmd"
	"os"
)

func main() {
	// defer panicRecover()
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func panicRecover() {
	// err := recover()
	// if err != nil {
	// 	fmt.Println("Fatal error!")
	// }
}
