package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"net/url"
	"strings"
)

type ginCallHandle struct {
	iface     *ast.InterfaceType
	ifaceName string
	baseuri   string
}

func (g *ginCallHandle) registerGinHandlers() ([]*ast.CallExpr, error) {
	hasil := []*ast.CallExpr{}
	for _, met := range g.iface.Methods.List {

		if g.extractMethod(met.Doc.List) == "" {
			continue
		}

		metast := g.genGinHandler(met)
		hasil = append(hasil, metast)
	}

	return hasil, nil
}

func (g *ginCallHandle) genGinHandler(met *ast.Field) *ast.CallExpr {
	methodfunc := "MethodGet"
	method := g.extractMethod(met.Doc.List)
	switch method {
	case "post":
		methodfunc = "MethodPost"
	case "get":
		methodfunc = "MethodGet"
	case "put":
		methodfunc = "MethodPut"
	case "delete":
		methodfunc = "MethodDelete"
	}

	uripath, _ := url.JoinPath(g.baseuri, convertCamelToSnake(met.Names[0].Name))

	// generating handler api
	ft, _ := met.Type.(*ast.FuncType)
	handler := ginHandler{
		ft:     ft,
		fnname: met.Names[0].Name,
	}

	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("g"),
			Sel: ast.NewIdent("Handle"),
		},
		Args: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("http"),
				Sel: ast.NewIdent(methodfunc),
			},
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, uripath),
			},
			handler.Ast(),
		},
	}
}

func (g *ginCallHandle) extractMethod(coms []*ast.Comment) string {
	for _, com := range coms {
		if !strings.HasPrefix(com.Text, `// method:`) {
			continue
		}

		method := strings.ReplaceAll(com.Text, `// method: `, "")
		method = strings.TrimSpace(method)
		return method
	}

	return ""
}
