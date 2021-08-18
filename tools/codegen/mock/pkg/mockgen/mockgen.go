// Package mockgen generates mocks for api actions.
package mockgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"net/http"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/mavolin/dismock/tools/codegen/mock/config"
	"github.com/mavolin/dismock/tools/codegen/mock/pkg/parser"
)

const returnParamName = "_ret"

func FilesFromParserFiles(pfiles ...parser.File) ([]File, error) {
	files := make([]File, len(pfiles))

	for i, f := range pfiles {
		actions, err := ActionsFromMethods(f.Methods...)
		if err != nil {
			return nil, errors.Wrap(err, f.Name)
		}

		files[i] = File{
			Name:    f.Name,
			Actions: actions,
		}
	}

	return files, nil
}

type File struct {
	Name    string
	Actions []Action
}

func ActionsFromMethods(methods ...parser.Method) ([]Action, error) {
	actions := make([]Action, len(methods))

	for i, m := range methods {
		cfg := config.C[m.Name]

		a := Action{
			Name:                 m.Name,
			Params:               make([]Param, 0, len(m.Params)),
			Multipart:            m.Multipart,
			WrappedResponseField: cfg.WrappedResponse,
			HTTPMethod:           HTTPMethod(m.HTTPMethod),
		}

		for _, methodParams := range m.Params {
			if ip := cfg.InferParam(methodParams.Name); ip != nil {
				continue
			}

			a.Params = append(a.Params, varFromParserParam(methodParams))
		}

		if m.ReturnType != nil {
			t := typeFromParserType(*m.ReturnType)
			a.ReturnType = &t

			retParam := Param{
				Name: returnParamName,
				Type: t,
			}
			// remove possible pointers
			retParam.Type.ElemPtr = false

			a.Params = append(a.Params, retParam)
		}

		if cfg.URLParams != nil {
			a.URLParams = make([]URLParam, len(cfg.URLParams))

			for j, mURLParam := range cfg.URLParams {
				var inferParam string
				if ip := cfg.InferParam(mURLParam.Param); ip != nil {
					inferParam = returnParamName + "." + ip.From
				}

				var paramType Type
				if mURLParam.Param != "" {
					for _, p := range m.Params {
						if p.Name == mURLParam.Param {
							paramType = typeFromParserType(p.Type)
							break
						}
					}
				}

				aURLParam, err := urlParamFromConfigURLParam(mURLParam, &paramType, inferParam)
				if err != nil {
					return nil, errors.Wrap(err, m.Name)
				}

				a.URLParams[j] = aURLParam
			}
		} else if m.QueryParam != "" {
			if ip := cfg.InferParam(m.QueryParam); ip != nil {
				a.QueryParam = returnParamName + "." + ip.From
			} else {
				a.QueryParam = m.QueryParam
			}
		}

		if cfg.JSONBody != nil {
			a.JSONBody = make([]JSONField, len(cfg.JSONBody))

			for j, mJSONField := range cfg.JSONBody {
				var inferParam string
				if ip := cfg.InferParam(mJSONField.Param); ip != nil {
					inferParam = returnParamName + "." + ip.From
				}

				var paramType Type
				if mJSONField.Param != "" {
					for _, p := range m.Params {
						if p.Name == mJSONField.Param {
							paramType = typeFromParserType(p.Type)
							break
						}
					}
				}

				aJSONField, err := jsonFieldFromConfigJSONField(mJSONField, &paramType, inferParam)
				if err != nil {
					return nil, errors.Wrap(err, m.Name)
				}

				a.JSONBody[j] = aJSONField
			}
		} else if m.JSONParam != "" {
			if ip := cfg.InferParam(m.JSONParam); ip != nil {
				a.JSONParam = returnParamName + "." + ip.From
			} else {
				a.JSONParam = m.JSONParam
			}
		}

		if m.ReasonParam != "" {
			if ip := cfg.InferParam(m.ReasonParam); ip != nil {
				a.ReasonParam = returnParamName + "." + ip.From
			} else {
				a.ReasonParam = m.ReasonParam
			}
		}

		if err := resolveEndpointExpr(m.EndpointExpr, cfg, &a.EndpointExpr); err != nil {
			return nil, errors.Wrap(err, m.Name)
		}

		actions[i] = a
	}

	return actions, nil
}

func resolveEndpointExpr(expr ast.Expr, cfg config.Item, ep *string) error {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		*ep = expr.Value + *ep
		return nil
	case *ast.CallExpr:
		if len(expr.Args) != 0 {
			return errors.New("function call with args when constructing endpoint")
		}

		funSel, ok := expr.Fun.(*ast.SelectorExpr)
		if !ok {
			return fmt.Errorf(
				"expected function call when constructing endpoint to be of type *ast.SelectorExpr, but is %T",
				expr.Fun)
		}

		xident := funSel.X.(*ast.Ident).Name

		if ip := cfg.InferParam(xident); ip != nil {
			xident = returnParamName + "." + ip.From
		}

		*ep = fmt.Sprintf("%s.%s()%s", xident, funSel.Sel.Name, *ep)
		return nil
	case *ast.Ident:
		relEndpoint, err := endpointVarForRelative(expr.Name)
		if err != nil {
			// it's a var other than an endpoint var

			if ip := cfg.InferParam(expr.Name); ip != nil {
				*ep = returnParamName + "." + ip.From + *ep
				return nil
			}

			*ep = expr.Name + *ep
			return nil
		}

		if len(relEndpoint) > 0 {
			*ep = `"` + relEndpoint + `"` + *ep
		} else {
			// remove the '+', since this was part of a BinaryExpr
			*ep = (*ep)[1:]
		}

	case *ast.BinaryExpr:
		if err := resolveEndpointExpr(expr.Y, cfg, ep); err != nil {
			return err
		}

		*ep = "+" + *ep

		return resolveEndpointExpr(expr.X, cfg, ep)
	default:
		return fmt.Errorf("unexpected expression in enpoint expression of type %T", expr)
	}

	return nil
}

type Action struct {
	Name       string
	Params     []Param
	ReturnType *Type // name is always 'ret'

	Multipart                          bool
	QueryParam, JSONParam, ReasonParam string

	URLParams []URLParam
	JSONBody  []JSONField

	WrappedResponseField string

	HTTPMethod   HTTPMethod
	EndpointExpr string
}

func (a *Action) AllowsAuditLogReason() bool {
	for _, param := range a.Params {
		if param.Type.Package == "api" && param.Type.Name == "AuditLogReason" {
			return true
		}
	}

	return false
}

func varFromParserParam(p parser.Param) Param {
	return Param{Name: p.Name, Type: typeFromParserType(p.Type)}
}

type Param struct {
	Name string
	Type Type
}

func typeFromParserType(t parser.Type) Type {
	return Type{
		Slice:    t.Slice,
		Variadic: t.Variadic,
		ElemPtr:  t.ElemPtr,
		Package:  t.Package,
		Name:     t.Name,
	}
}

func typeFromString(t string) Type {
	split := strings.SplitN(t, ".", 2)
	if len(split) == 1 && split[0] == "dismock" {
		return Type{Name: split[1]}
	} else if len(split) == 1 {
		return Type{Name: split[0]}
	}

	return Type{Package: split[0], Name: split[1]}
}

type Type struct {
	Slice    uint8
	Variadic bool
	ElemPtr  bool
	Package  string
	Name     string
}

func (t *Type) String() string {
	var prefix string

	if t.Slice > 0 {
		prefix = strings.Repeat("[]", int(t.Slice))
	} else if t.Variadic {
		prefix = "..."
	}

	if t.ElemPtr {
		prefix += "*"
	}

	if t.Package == "" {
		return prefix + t.Name
	}

	return prefix + t.Package + "." + t.Name
}

func urlParamFromConfigURLParam(up config.URLParam, paramType *Type, inferParam string) (URLParam, error) {
	newP := URLParam{
		GoName:    strcase.ToCamel(up.Name),
		Name:      up.Name,
		Omitempty: up.Omitemtpy,
		Param:     up.Param,
		Literal:   !strings.Contains(up.Value, "{{.Var}}"),
		Value:     up.Value,
	}

	if newP.Name == "" {
		newP.Name = strcase.ToSnake(up.Param)
		newP.GoName = strcase.ToCamel(up.Param)
	}

	if up.Type != "" {
		newP.Type = typeFromString(up.Type)
	} else {
		newP.Type = *paramType
	}

	if inferParam != "" {
		newP.Param = inferParam
	}

	if up.Value == "" {
		newP.Value = newP.Param
	}

	if newP.Literal {
		return newP, nil
	}

	t := template.New("")
	t, err := t.Parse(up.Value)
	if err != nil {
		return URLParam{}, err
	}

	var buf bytes.Buffer

	data := struct {
		Param string
		Var   string
	}{
		Param: up.Param,
		Var:   "_params." + newP.GoName,
	}

	if err := t.Execute(&buf, data); err != nil {
		return URLParam{}, err
	}

	newP.Value = strings.TrimSpace(buf.String())

	return newP, nil
}

type URLParam struct {
	GoName    string
	Name      string
	Type      Type
	Omitempty bool
	Param     string
	Literal   bool
	Value     string
}

func jsonFieldFromConfigJSONField(f config.JSONField, paramType *Type, inferParam string) (JSONField, error) {
	newF := JSONField{
		GoName:    strcase.ToCamel(f.Name),
		Name:      f.Name,
		Omitempty: f.Omitemtpy,
		Param:     f.Param,
		Literal:   !strings.Contains(f.Value, "{{.Var}}"),
		Value:     f.Value,
	}

	if newF.Name == "" {
		newF.Name = strcase.ToSnake(f.Param)
		newF.GoName = strcase.ToCamel(f.Param)
	}

	if f.Type != "" {
		newF.Type = typeFromString(f.Type)
	} else {
		newF.Type = *paramType
	}

	if inferParam != "" {
		newF.Param = inferParam
	}

	if f.Value == "" {
		newF.Value = newF.Param
	}

	if newF.Literal {
		return newF, nil
	}

	t := template.New("")
	t, err := t.Parse(f.Value)
	if err != nil {
		return JSONField{}, err
	}

	var buf bytes.Buffer

	data := struct {
		Param string
		Var   string
	}{
		Param: f.Param,
		Var:   "_body." + newF.GoName,
	}

	if err := t.Execute(&buf, data); err != nil {
		return JSONField{}, err
	}

	newF.Value = strings.TrimSpace(buf.String())

	return newF, nil
}

type JSONField struct {
	GoName    string
	Name      string
	Type      Type
	Omitempty bool
	Param     string
	Literal   bool
	Value     string
}

type HTTPMethod string

func (m HTTPMethod) AsHTTPVar() string {
	switch m {
	case http.MethodGet:
		return "http.MethodGet"
	case http.MethodHead:
		return "http.MethodHead"
	case http.MethodPost:
		return "http.MethodPost"
	case http.MethodPut:
		return "http.MethodPut"
	case http.MethodPatch:
		return "http.MethodPatch"
	case http.MethodDelete:
		return "http.MethodDelete"
	case http.MethodConnect:
		return "http.MethodConnect"
	case http.MethodOptions:
		return "http.MethodOptions"
	case http.MethodTrace:
		return "http.MethodTrace"
	default:
		return ""
	}
}

func endpointVarForRelative(evar string) (string, error) {
	switch evar {
	case "EndpointApplications":
		return "applications/", nil
	case "EndpointChannels":
		return "channels/", nil
	case "EndpointGuilds":
		return "guilds/", nil
	case "EndpointInteractions":
		return "interactions/", nil
	case "EndpointInvites":
		return "invites/", nil
	case "EndpointStageInstances":
		return "stage-instances/", nil
	case "EndpointWebhooks":
		return "webhooks/", nil
	case "EndpointUsers":
		return "users/", nil
	case "EndpointMe":
		return "users/@me", nil
	case "EndpointLogin":
		return "auth/login", nil
	case "EndpointTOTP":
		return "auth/mfa/totp", nil
	case "Endpoint":
		return "", nil
	default:
		return "", fmt.Errorf("unknown endpoint var '%s'", evar)
	}
}
