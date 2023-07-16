package main

import (
	"fmt"
	"os"

	"github.com/kgaughan/gcredstash/internal/command"
)

func main() {
	if err := command.Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
