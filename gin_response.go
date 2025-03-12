package main

import (
	"go/ast"
	"go/token"
)

func standardRes(varname string) *ast.ExprStmt {
	var res ast.Expr
	if varname == "" {
		res = &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("gin"),
				Sel: ast.NewIdent("H"),
			},
			Elts: []ast.Expr{
				&ast.KeyValueExpr{
					Key:   &ast.BasicLit{Kind: token.STRING, Value: `"message"`},
					Value: &ast.BasicLit{Kind: token.STRING, Value: `""`},
				},
			},
		}
	} else {
		res = ast.NewIdent(varname)
	}

	callExpr := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("ctx"),
				Sel: ast.NewIdent("JSON"),
			},
			Args: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("http"),
					Sel: ast.NewIdent("StatusOK"),
				},
				res,
			},
		},
	}

	return callExpr
}

func standardErrorStmt() *ast.IfStmt {
	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent("err"),
			Op: token.NEQ,
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("ctx"),
							Sel: ast.NewIdent("AbortWithStatusJSON"),
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("http"),
								Sel: ast.NewIdent("StatusInternalServerError"),
							},
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("gin"),
									Sel: ast.NewIdent("H"),
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   &ast.BasicLit{Kind: token.STRING, Value: `"message"`},
										Value: &ast.CallExpr{Fun: &ast.SelectorExpr{X: ast.NewIdent("err"), Sel: ast.NewIdent("Error")}},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
