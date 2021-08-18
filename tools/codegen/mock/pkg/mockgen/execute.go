package mockgen

import (
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/mavolin/dismock/tools/codegen/mock/assets"
)

func Generate(files []File) error {
	log.Println("generating generated_mocks.go")

	f, err := os.Create("generated_mocks.go")
	if err != nil {
		return err
	}

	t := template.New("generated_mocks.tmpl")
	t.Funcs(template.FuncMap{
		"indent": func(lvl int, s string) string {
			split := strings.SplitAfter(s, "\n")

			for i, s := range split {
				split[i] = strings.Repeat("    ", lvl) + s
			}

			return strings.Join(split, "")
		},
	})

	if _, err := t.ParseFS(assets.FS, "templates/*.tmpl"); err != nil {
		return err
	}

	return t.Execute(f, files)
}
