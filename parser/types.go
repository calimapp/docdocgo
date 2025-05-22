package parser

import "html/template"

type goFunction struct {
	Name string
	Doc  string
}

type goType struct {
	Name string
	Doc  string
}

type goPackage struct {
	Name      string
	Ref       string
	Doc       string
	Functions []goFunction
	Types     []goType
}

type goModule struct {
	Name          string
	Readme        template.HTML
	Documentation string
	Packages      []goPackage
	SourceFiles   []string
}
