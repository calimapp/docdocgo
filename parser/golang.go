package parser

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func getPackageDocumentation(packagePath string) (*goDocumentation, error) {
	fset := token.NewFileSet()

	astPkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, astPkg := range astPkgs {
		pkgDoc := doc.New(astPkg, packagePath, doc.AllDecls)
		ref := pkgDoc.ImportPath
		if ref == "." {
			ref = "module"
		}
		return &goDocumentation{
			Ref:       ref,
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
		doc, _ := getPackageDocumentation(packagePath)
		return &goPackage{
			Name:          pkgDoc.Name,
			Ref:           pkgDoc.ImportPath,
			Documentation: *doc,
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
				args = append(args, valueTypePair{Name: field.Names[0].Name, Type: astTypeToString(field.Type)})
			}
		}
		// parse return
		results := make([]valueTypePair, 0)
		if fn.Decl.Type.Results != nil {
			for _, field := range fn.Decl.Type.Results.List {
				results = append(results, valueTypePair{Type: astTypeToString(field.Type)})
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
				fields = append(fields, valueTypePair{Name: field.Names[0].Name, Type: astTypeToString(field.Type)})
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
		for _, spec := range cst.Decl.Specs {
			var value string
			if len(spec.(*ast.ValueSpec).Values) > 0 {
				value = astValueToString(spec.(*ast.ValueSpec).Values[0])
			}
			goConst := goVar{
				Name:  spec.(*ast.ValueSpec).Names[0].Name,
				Value: value,
				Type:  astTypeToString(spec.(*ast.ValueSpec).Type),
				Doc:   spec.(*ast.ValueSpec).Doc.Text(),
			}
			consts = append(consts, goConst)
		}
	}
	return consts
}

// retrive name, value, type and documentation of all variables in given package
func parsePackageVariables(pkgDocumentation *doc.Package) []goVar {
	vars := make([]goVar, 0)
	for _, v := range pkgDocumentation.Vars {
		for _, spec := range v.Decl.Specs {
			var value string
			if len(spec.(*ast.ValueSpec).Values) > 0 {
				value = astValueToString(spec.(*ast.ValueSpec).Values[0])
			}
			goVar := goVar{
				Name:  spec.(*ast.ValueSpec).Names[0].Name,
				Value: value,
				Type:  astTypeToString(spec.(*ast.ValueSpec).Type),
				Doc:   spec.(*ast.ValueSpec).Doc.Text(),
			}
			vars = append(vars, goVar)
		}
	}
	return vars
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

func astTypeToString(value ast.Expr) string {
	switch t := value.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + astTypeToString(t.X)
	case *ast.ArrayType:
		return "[]" + astTypeToString(t.Elt)
	case *ast.SelectorExpr:
		return astTypeToString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return "map[" + astTypeToString(t.Key) + "]" + astTypeToString(t.Value)
	case *ast.ChanType:
		var dir string
		switch t.Dir {
		case ast.SEND:
			dir = "chan<- "
		case ast.RECV:
			dir = "<-chan "
		default:
			dir = "chan "
		}
		return dir + astTypeToString(t.Value)
	case *ast.FuncType:
		return "func(" + paramsToString(t.Params) + ")" + resultsToString(t.Results)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{...}"
	case *ast.Ellipsis:
		return "..." + astTypeToString(t.Elt)
	default:
		return ""
	}
}

func astValueToString(value ast.Expr) string {
	switch t := value.(type) {
	case *ast.BasicLit:
		return t.Value
	case *ast.CompositeLit:
		var values string
		for _, elt := range t.Elts {
			values += " " + astValueToString(elt)
		}
		return values
	default:
		return ""
	}
}

func paramsToString(fields *ast.FieldList) string {
	var params []string
	if fields == nil {
		return ""
	}
	for _, field := range fields.List {
		var names []string
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
		if len(names) == 0 {
			params = append(params, astTypeToString(field.Type))
		} else {
			for _, name := range names {
				params = append(params, name+" "+astTypeToString(field.Type))
			}
		}
	}
	return strings.Join(params, ", ")
}

func resultsToString(fields *ast.FieldList) string {
	if fields == nil || fields.List == nil {
		return ""
	}
	var results []string
	for _, field := range fields.List {
		results = append(results, astTypeToString(field.Type))
	}
	if len(results) > 1 {
		return " (" + strings.Join(results, ", ") + ")"
	} else if len(results) == 1 {
		return " " + results[0]
	}
	return ""
}
