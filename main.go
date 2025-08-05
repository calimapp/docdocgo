// DocDocGO is a golang documentation generator
package main

import (
	"embed"
	"fmt"

	"github.com/calimapp/docdocgo/parser"
	"github.com/spf13/cobra"
)

//go:embed src/*
var templates embed.FS

var (
	output     string
	modVersion string

	cmd = &cobra.Command{
		Use:   "docdocgo <module-path>",
		Short: "A golang documentation generator",
		Long:  "Docdocgo render a golang documentation in a single html file",
		Args:  cobra.ExactArgs(1),
		RunE:  run,
	}
)

func main() {
	cmd.Flags().StringVarP(&output, "output", "o", "out.html", "html output filepath")
	cmd.Flags().StringVar(&modVersion, "mod-version", "", "manually set module version")
	cmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	modulePath := args[0]

	documentation, err := parser.ParseModule(modulePath)
	if err != nil {
		return fmt.Errorf("error generating doc: %w", err)
	}
	if modVersion != "" {
		documentation.Version = modVersion
	}

	if err = documentation.ToHTML(templates, output); err != nil {
		return fmt.Errorf("error rendering doc html: %w", err)
	}
	return nil
}
