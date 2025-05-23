package parser

import "html/template"

type valueTypePair struct {
	Name string
	Type string
}

type goFunction struct {
	Name      string
	Doc       string
	Arguments []valueTypePair
	Results   []valueTypePair
}

type goType struct {
	Name    string
	Type    string
	Doc     string
	Fields  []valueTypePair
	Methods []goFunction
}

type goVar struct {
	Name  string
	Value string
	Type  string
	Doc   string
}

type goDocumentation struct {
	Overview  string
	Constants []goVar
	Variables []goVar
	Types     []goType
	Functions []goFunction
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
	Documentation goDocumentation
	Packages      []goPackage
	SourceFiles   []string
}
