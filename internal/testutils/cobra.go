package testutils

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

func NewDummyCommand(ctx context.Context) (*cobra.Command, *strings.Builder) {
	out := &strings.Builder{}
	cmd := &cobra.Command{}
	cmd.SetContext(ctx)
	cmd.SetOut(out)
	return cmd, out
}
