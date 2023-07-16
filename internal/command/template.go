package command

import (
	"bytes"
	"fmt"
	"os"

	gcredstash "github.com/kgaughan/gcredstash/internal"
	"github.com/kgaughan/gcredstash/internal/templating"
	"github.com/spf13/cobra"
)

func MakeTemplateCmd(driver *gcredstash.Driver, common *CommonFlags) *cobra.Command {
	var inplace bool

	cmd := &cobra.Command{
		Use:   "template",
		Short: "Parse a tempalte file with credentials",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tmplFile := args[0]

			var content string
			if tmplFile == "-" {
				content = gcredstash.ReadStdin()
			} else {
				var err error
				content, err = gcredstash.ReadFile(tmplFile)
				if err != nil {
					return fmt.Errorf("cannot read %q: %w", tmplFile, err)
				}
			}
			tmpl, err := templating.MakeTemplate(driver, common.Table).Parse(content)
			if err != nil {
				return fmt.Errorf("cannot parse %q template: %w", tmplFile, err)
			}

			buf := &bytes.Buffer{}
			if err = tmpl.Execute(buf, nil); err != nil {
				return fmt.Errorf("cannot execute %q template: %w", tmplFile, err)
			}

			if inplace {
				if err := os.WriteFile(tmplFile, buf.Bytes(), 0o644); err != nil { //nolint:gosec
					return fmt.Errorf("cannot write to %q: %w", tmplFile, err)
				}
			}

			fmt.Print(buf.String())
			return nil
		},
	}
	cmd.Flags().BoolVarP(&inplace, "inplace", "i", false, "overwrite the template file")

	return cmd
}
