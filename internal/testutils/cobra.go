package testutils

import (
	"strings"

	"github.com/spf13/cobra"
)

func NewDummyCommand() (*cobra.Command, *strings.Builder) {
	out := &strings.Builder{}
	cmd := &cobra.Command{}
	cmd.SetOutput(out)
	return cmd, out
}
