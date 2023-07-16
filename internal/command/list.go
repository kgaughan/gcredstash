package command

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List credentials and their version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		driver, err := internal.NewDriver()
		if err != nil {
			return err //nolint:wrapcheck
		}

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
				lines = append(lines, fmt.Sprintf("%-*s -- version: %d\n", maxKeyLen, *name, versionNum))
			}
		}
		sort.Strings(lines)
		for line := range lines {
			cmd.Println(line)
		}
		return nil
	},
}

func init() {
	Root.AddCommand(listCmd)
}
