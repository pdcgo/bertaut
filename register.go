package bertaut

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
)

func RegisterApi(ctx context.Context, serviceName string, handle func(ctx context.Context) []ast.Stmt) *ast.FuncDecl {
	funcname := "Register" + serviceName + "Api"

	body := handle(ctx)

	return &ast.FuncDecl{
		Name: ast.NewIdent(funcname),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("srv")}, // Parameter name: usr
						Type:  ast.NewIdent(serviceName),         // Parameter type: UserApi
					},
					{
						Names: []*ast.Ident{ast.NewIdent("g")}, // Parameter name: g
						Type: &ast.StarExpr{ // Pointer type: *gin.RouterGroup
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("gin"),         // Package name: gin
								Sel: ast.NewIdent("RouterGroup"), // Struct name: RouterGroup
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func RegisterGin(
	ctx context.Context,
	method string,
	uri string,
	handle func(ctx context.Context)) *ast.CallExpr {

	fn := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("gin"),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{},
	}

	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("g"),
			Sel: ast.NewIdent("Handle"),
		},
		Args: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("http"),
				Sel: ast.NewIdent(method),
			},
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, uri),
			},
			fn,
		},
	}
}
