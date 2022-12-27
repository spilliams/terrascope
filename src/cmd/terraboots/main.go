package main

import (
	"fmt"
	"os"

	"github.com/spilliams/terraboots/internal/cli"
)

func main() {
	err := cli.NewTerrabootsCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
