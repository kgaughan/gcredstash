package templating

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/mattn/go-shellwords"
)

func MakeTemplate(driver *internal.Driver, table string) *template.Template {
	return template.New("template").Funcs(template.FuncMap{
		"get": func(credential string, cxt ...string) (string, error) {
			context, err := internal.ParseContext(cxt)
			if err != nil {
				return "", fmt.Errorf("could not parse context: %w", err)
			}

			value, err := driver.GetSecret(credential, "", table, context)
			if err != nil {
				return "", fmt.Errorf("could not fetch credentials: %w", err)
			}

			return value, nil
		},
		"env": func(key string) (string, error) {
			return os.Getenv(key), nil
		},
		"sh": func(line string) (string, error) {
			cmd, err := shellwords.Parse(line)
			if err != nil {
				return "", fmt.Errorf("could not parse command: %w", err)
			}

			var out []byte

			switch len(cmd) {
			case 0:
				out = []byte{}
			case 1:
				//nolint:gosec
				out, err = exec.Command(cmd[0]).Output()
			default:
				//nolint:gosec
				out, err = exec.Command(cmd[0], cmd[1:]...).Output()
			}

			if err != nil {
				return "", fmt.Errorf("could not run command: %w", err)
			}

			str := string(out)

			return strings.TrimRight(str, "\n"), nil
		},
	})
}
