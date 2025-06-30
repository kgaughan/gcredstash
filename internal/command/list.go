package command

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func listImpl(cmd *cobra.Command, _ []string, driver *internal.Driver, out io.Writer) error {
	items, err := driver.ListSecrets(table)
	if err != nil {
		return err //nolint:wrapcheck
	}

	maxKeyLen := internal.MaxKeyLen(items)
	lines := make([]string, 0, len(items))
	for name, version := range items {
		versionNum, err := strconv.Atoi(*version)
		if err != nil {
			cmd.PrintErrf("bad version for %q: %q\n", *name, *version)
		} else {
			lines = append(lines, fmt.Sprintf("%-*s -- version: %d", maxKeyLen, *name, versionNum))
		}
	}
	sort.Strings(lines)
	for _, line := range lines {
		fmt.Fprintln(out, line)
	}
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List credentials and their version",
		Args:  cobra.NoArgs,
		RunE:  wrapWithDriver(listImpl),
	}

	Root.AddCommand(cmd)
}
