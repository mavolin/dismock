package parser

import (
	"go/ast"
)

type (
	File struct {
		Name    string
		Methods []Method
	}

	Method struct {
		Name       string
		Params     []Param
		ReturnType *Type

		Multipart                          bool
		QueryParam, JSONParam, ReasonParam string

		HTTPMethod   string
		EndpointExpr ast.Expr
	}

	Param struct {
		Name string
		Type Type
	}

	Type struct {
		Slice    uint8
		Variadic bool
		ElemPtr  bool
		Package  string
		Name     string
	}
)
