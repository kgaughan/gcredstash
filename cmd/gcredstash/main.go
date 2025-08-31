package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kgaughan/gcredstash/internal/command"
)

func main() {
	if err := command.Root.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
