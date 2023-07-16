package command

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func MakeListCmd(driver *internal.Driver, common *CommonFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List credentials and their version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := driver.ListSecrets(common.Table)
			if err != nil {
				return err //nolint:wrapcheck
			}

			maxKeyLen := internal.MaxKeyLen(items)
			lines := make([]string, 0, len(items))
			for name, version := range items {
				versionNum, err := strconv.Atoi(*version)
				if err != nil {
					fmt.Fprintf(os.Stderr, "bad version for %q: %q\n", *name, *version)
				} else {
					lines = append(lines, fmt.Sprintf("%-*s -- version: %d\n", maxKeyLen, *name, versionNum))
				}
			}
			sort.Strings(lines)
			for line := range lines {
				fmt.Println(line)
			}
			return nil
		},
	}
}
