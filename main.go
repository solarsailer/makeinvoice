package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/solarsailer/makeinvoice/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		os.Exit(-1)
	}
}
