package main

import (
	"fmt"
	"os"

	"github.com/spilliams/terrascope/internal/cli"
)

func main() {
	err := cli.NewTerrascopeCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
