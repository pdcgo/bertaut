package main

import "go/ast"

func AstNil() ast.Expr {
	return &ast.Ident{Name: "nil"}
}
