package parser

import (
	"go/ast"
	"strings"
)

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
