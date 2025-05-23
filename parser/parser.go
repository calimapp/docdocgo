package parser

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path/filepath"
)

func getPackageDocumentation(packagePath string) (*goDocumentation, error) {
	fset := token.NewFileSet()

	astPkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, astPkg := range astPkgs {
		pkgDoc := doc.New(astPkg, packagePath, doc.AllDecls)
		return &goDocumentation{
			Overview:  pkgDoc.Doc,
			Constants: parsePackageConstants(pkgDoc),
			Variables: parsePackageVariables(pkgDoc),
			Types:     parsePackageTypes(pkgDoc),
			Functions: parsePackageFunctions(pkgDoc),
		}, nil
	}
	return nil, fmt.Errorf("no packages found at %s", packagePath)
}

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
		// parse params
		args := make([]valueTypePair, 0)
		if fn.Decl.Type.Params != nil {
			for _, field := range fn.Decl.Type.Params.List {
				args = append(args, valueTypePair{Name: field.Names[0].Name, Type: field.Type.(*ast.Ident).Name})
			}
		}
		// parse return
		results := make([]valueTypePair, 0)
		if fn.Decl.Type.Results != nil {
			for _, field := range fn.Decl.Type.Results.List {
				results = append(results, valueTypePair{Type: field.Type.(*ast.Ident).Name})
			}
		}
		goFunctions = append(goFunctions, goFunction{Name: fn.Name, Doc: fn.Doc, Arguments: args, Results: results})
	}
	return goFunctions
}

func parsePackageTypes(pkgDocumentation *doc.Package) []goType {
	goTypes := []goType{}
	for _, ty := range pkgDocumentation.Types {
		fields := make([]valueTypePair, 0)
		// methods := make([]goFunction, 0)
		var typeType string
		switch t := ty.Decl.Specs[0].(*ast.TypeSpec).Type.(type) {
		case *ast.StructType:
			typeType = "struct"
			for _, field := range t.Fields.List {
				fields = append(fields, valueTypePair{Name: field.Names[0].Name, Type: field.Type.(*ast.Ident).Name})
			}
		case *ast.InterfaceType:
			typeType = "interface"
			// for _, method := range t.Methods.List {
			// methods = append(methods, goFunction{Name: method.Names[0].Name, Argument})
			// TODO
			// }
		case *ast.Ident:
			typeType = t.String()
		}
		goTypes = append(goTypes, goType{Name: ty.Name, Type: typeType, Doc: ty.Doc, Fields: fields})
	}
	return goTypes
}

// retrive name, value, type and documentation of all constants in given package
func parsePackageConstants(pkgDocumentation *doc.Package) []goVar {
	consts := make([]goVar, 0)
	for _, cst := range pkgDocumentation.Consts {
		goConst := goVar{
			Name:  cst.Names[0],
			Value: cst.Decl.Specs[0].(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value,
			Type:  cst.Decl.Specs[0].(*ast.ValueSpec).Type.(*ast.Ident).Name,
			Doc:   cst.Doc,
		}
		consts = append(consts, goConst)
	}
	return consts
}

// retrive name, value, type and documentation of all variables in given package
func parsePackageVariables(pkgDocumentation *doc.Package) []goVar {
	vars := make([]goVar, 0)
	for _, v := range pkgDocumentation.Vars {
		goVar := goVar{
			Name:  v.Names[0],
			Value: v.Decl.Specs[0].(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value,
			Type:  v.Decl.Specs[0].(*ast.ValueSpec).Type.(*ast.Ident).Name,
			Doc:   v.Doc,
		}
		vars = append(vars, goVar)
	}
	return vars
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
