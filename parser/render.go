package parser

import (
	"embed"
	"html/template"
	"os"
)

func (m *goModule) ToHTML(templates embed.FS, outputPath string) error {
	tmpl, err := template.ParseFS(templates, "src/*")
	if err != nil {
		return err
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(outputFile, "index.html", m)
}
