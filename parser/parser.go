package parser

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path/filepath"
)

func parsePackage(packagePath string) (*goPackage, error) {
	fset := token.NewFileSet()

	astPkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, astPkg := range astPkgs {
		pkgDoc := doc.New(astPkg, packagePath, doc.AllDecls)

		return &goPackage{
			Name:      pkgDoc.Name,
			Doc:       pkgDoc.Doc,
			Ref:       pkgDoc.ImportPath,
			Functions: parsePackageFunctions(pkgDoc),
			Types:     parsePackageTypes(pkgDoc),
		}, nil
	}
	return nil, fmt.Errorf("no packages found at %s", packagePath)
}

func parsePackageFunctions(pkgDocumentation *doc.Package) []goFunction {
	goFunctions := []goFunction{}
	for _, fn := range pkgDocumentation.Funcs {
		goFunctions = append(goFunctions, goFunction{Name: fn.Name, Doc: fn.Doc})
	}
	return goFunctions
}

func parsePackageTypes(pkgDocumentation *doc.Package) []goType {
	goTypes := []goType{}
	for _, ty := range pkgDocumentation.Types {
		goTypes = append(goTypes, goType{Name: ty.Name, Doc: ty.Doc})
	}
	return goTypes
}

func (m *goModule) ToHTML(outputPath string) error {
	tmpl, err := template.ParseFiles("src/doc.html")
	if err != nil {
		return err
	}
	outputFile, _ := os.Create(outputPath)
	return tmpl.Execute(outputFile, m)
}

// return .go sources files of a given directory
// if no sourceFiles found or error, it return an empty list
func getSourceFiles(dirPath string) []string {
	sourceFiles := make([]string, 0)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return sourceFiles
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".go" {
			sourceFiles = append(sourceFiles, entry.Name())
		}
	}
	return sourceFiles
}
