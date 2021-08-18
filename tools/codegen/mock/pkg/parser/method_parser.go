package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type methodParser struct {
	fun *ast.FuncDecl
	Method
}

func newMethodParser(fun *ast.FuncDecl) *methodParser {
	return &methodParser{fun: fun}
}

func (p *methodParser) Parse() (*Method, error) {
	p.Name = p.fun.Name.Name

	if err := p.resolveParams(); err != nil {
		return nil, errors.Wrap(err, p.fun.Name.Name)
	}

	if err := p.resolveReturnType(); err != nil {
		return nil, errors.Wrap(err, p.fun.Name.Name)
	}

	return &p.Method, errors.Wrap(p.resolveRequestMeta(), p.fun.Name.Name)
}

func (p *methodParser) resolveParams() error {
	p.Params = make([]Param, 0, len(p.fun.Type.Params.List))

	for _, param := range p.fun.Type.Params.List {
		// The reason why there are multiple names for a single param is
		// undocumented.
		// I suspect this will hold multiple values if multiple parameters use
		// the short notation to share the same type, e.g.
		// func foo(a, b string).
		// Therefore, for now, multiple names are treated as multiple params.
		for _, name := range param.Names {
			t, err := resolveType(param.Type)
			if err != nil {
				return errors.Wrapf(err, "param %d", len(p.Params)+1)
			}

			p.Params = append(p.Params, Param{
				Name: name.Name,
				Type: *t,
			})
		}
	}

	return nil
}

func (p *methodParser) resolveReturnType() error {
	returns := p.fun.Type.Results.List
	if len(returns) <= 1 {
		return nil
	} else if len(returns) > 2 {
		return errors.New("methods with more than 2 returns are not supported")
	}

	t, err := resolveType(returns[0].Type)
	if err != nil {
		return errors.Wrap(err, "return")
	}

	p.ReturnType = t
	return nil
}

func (p *methodParser) resolveRequestMeta() error {
	lastStmt := p.fun.Body.List[len(p.fun.Body.List)-1]

	returnStmt, ok := lastStmt.(*ast.ReturnStmt)
	if !ok {
		return fmt.Errorf("expected last statement be a return statement, but got %T", lastStmt)
	}

	call, ok := returnStmt.Results[len(returnStmt.Results)-1].(*ast.CallExpr)
	if !ok {
		return fmt.Errorf("expected last return value to be generated using a function call, but got %T", lastStmt)
	}

	funSel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return fmt.Errorf("expected return function to call a function from a different package, but found %T",
			call.Fun)
	}

	xident, ok := funSel.X.(*ast.Ident)
	if !ok {
		return fmt.Errorf("expected return function call to either be function of package sendpart or method of "+
			"api.Client, but found %+v", funSel.X)
	}

	if xident.Name == "sendpart" {
		return p.resolveMultipart(funSel.Sel.Name, call)
	}

	var offset int

	switch funSel.Sel.Name {
	case "FastRequest":
	case "RequestJSON":
		offset = 1
	default:
		return fmt.Errorf("unexpected return function call %+v", funSel)
	}

	if err := p.resolveHTTPMethod(call.Args[0+offset]); err != nil {
		return err
	}

	p.EndpointExpr = call.Args[1+offset]

	if err := p.resolveOptions(call.Args[2+offset:]...); err != nil {
		return err
	}

	return nil
}

func (p *methodParser) resolveHTTPMethod(methodExpr ast.Expr) error {
	method, ok := methodExpr.(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("expected http method to be a literal, but found %T", methodExpr)
	}

	if method.Kind != token.STRING {
		return fmt.Errorf("expected http method to be a string literal, but found %s", method.Kind.String())
	}

	p.HTTPMethod = strings.Trim(method.Value, `"`)
	return nil
}

func (p *methodParser) resolveOptions(optionExprs ...ast.Expr) error {
	for _, expr := range optionExprs {
		call, ok := expr.(*ast.CallExpr)
		if !ok {
			return nil
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return nil
		}

		xident, ok := sel.X.(*ast.Ident)
		if !ok || xident.Name != "httputil" {
			return nil
		}

		switch sel.Sel.Name {
		case "WithJSONBody":
			bodyIdent, ok := call.Args[0].(*ast.Ident)
			if !ok {
				return fmt.Errorf("expected the WithJSONBody option to take in a variable, but found %+v",
					call.Args[0])
			}

			p.JSONParam = bodyIdent.Name
		case "WithSchema":
			schemaIdent, ok := call.Args[1].(*ast.Ident)
			if ok { // no error checks, since sometimes url.Values literals are used
				p.QueryParam = schemaIdent.Name
			}
		case "WithHeaders":
			headerSel, ok := call.Args[0].(*ast.SelectorExpr)
			if !ok {
				return fmt.Errorf("expected the WithHeaders option to take in a selector, but found %+v", call.Args[0])
			}

			p.ReasonParam = headerSel.Sel.Name
		}
	}

	return nil
}

func (p *methodParser) resolveMultipart(httpMethod string, call *ast.CallExpr) error {
	p.HTTPMethod = httpMethod

	bodyIdent, ok := call.Args[1].(*ast.Ident)
	if !ok {
		return fmt.Errorf("expected the data given to sendpart to be a variable, but found %+v", call.Args[1])
	}

	p.JSONParam = bodyIdent.Name
	p.EndpointExpr = call.Args[3]

	return nil
}
