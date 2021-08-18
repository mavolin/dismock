package parser

import (
	"fmt"
	"go/ast"

	"github.com/pkg/errors"
)

func isPrimitive(i *ast.Ident) bool {
	switch i.Name {
	case "bool":
	case "uint8":
	case "uint16":
	case "uint32":
	case "uint64":
	case "uint":
	case "uintptr":
	case "int8":
	case "int16":
	case "int32":
	case "int64":
	case "int":
	case "float32":
	case "float64":
	case "complex64":
	case "complex128":
	case "string":
	case "byte":
	case "rune":
	default:
		return false
	}

	return true
}

func resolveType(typeExpr ast.Expr) (*Type, error) {
	switch typeExpr := typeExpr.(type) {
	case *ast.Ident:
		t := &Type{Name: typeExpr.Name}

		if !isPrimitive(typeExpr) {
			t.Package = "api"
		}

		return t, nil
	case *ast.SelectorExpr:
		xident, ok := typeExpr.X.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("unknown type: %T", typeExpr)
		}

		return &Type{Package: xident.Name, Name: typeExpr.Sel.Name}, nil
	case *ast.StarExpr:
		t, err := resolveType(typeExpr.X)
		if err != nil {
			return nil, err
		}

		if t.Slice > 0 {
			return nil, errors.New("pointer to slices are not supported")
		}

		t.ElemPtr = true
		return t, nil
	case *ast.ArrayType:
		t, err := resolveType(typeExpr.Elt)
		if err != nil {
			return nil, err
		}

		t.Slice++
		return t, nil
	case *ast.Ellipsis:
		t, err := resolveType(typeExpr.Elt)
		if err != nil {
			return nil, err
		}

		t.Variadic = true
		return t, nil
	}

	return nil, fmt.Errorf("unknown type: %T", typeExpr)
}
