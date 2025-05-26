// DocDocGO is a golang documentation generator
package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.gms.dev.lab/calimap/docdocgo/parser"
)

var (
	varA string = ""
	//test
	varB = ""
)

const (
	a string = ""
	//test
	b = ""
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: docdocgo <module-path>")
		os.Exit(1)
	}
	modulePath := os.Args[1]

	output := flag.String("output", "out.html", "--output <output-path>")
	flag.Parse()

	documentation, err := parser.ParseModule(modulePath)
	if err != nil {
		fmt.Println("Error generating doc:", err)
		os.Exit(1)
	}

	if err = documentation.ToHTML(*output); err != nil {
		fmt.Println("Error rendering doc html template:", err)
		os.Exit(1)
	}
}
