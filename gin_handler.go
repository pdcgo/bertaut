package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
)

type ginHandler struct {
	ft     *ast.FuncType
	fnname string
}

func (gh *ginHandler) parameterAst() {
	for _, d := range gh.ft.Params.List {
		log.Println(d.Names)
	}

	if gh.ft.TypeParams != nil {
		for _, d := range gh.ft.TypeParams.List {
			log.Println("type", d)
		}
	} else {
		log.Println("type nil")
	}

}

func (gh *ginHandler) prepareDeclAst() ([]ast.Stmt, []ast.Expr) {
	hasil := []ast.Stmt{}
	variable := []ast.Expr{}
	c := 0
	resname := func() string {
		c += 1
		return fmt.Sprintf("result%d", c)
	}

	for _, d := range gh.ft.Results.List {
		switch t := d.Type.(type) {
		case *ast.Ident:
			switch t.Name {
			case "error":
				hasil = append(hasil, &ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{ast.NewIdent("err")},
								Type:  ast.NewIdent(t.Name),
							},
						},
					},
				})
				variable = append(variable, ast.NewIdent("err"))
				continue
			default:
				varname := ast.NewIdent(resname())
				hasil = append(hasil, &ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{varname},
								Type:  ast.NewIdent(t.Name),
							},
						},
					},
				})
				variable = append(variable, varname)
			}
		case *ast.StructType:
			fmt.Println("Struct {")
			for _, field := range t.Fields.List {
				for _, fieldName := range field.Names {
					fmt.Printf("  %s: %s\n", fieldName.Name, printExprTypeInner(field.Type))
				}
			}
			fmt.Println("}")
		case *ast.StarExpr:
			// fmt.Println("Pointer to:", printExprTypeInner(t.X)) // Example: "*User"
		case *ast.SelectorExpr:
			fmt.Println("Selector:", t.Sel.Name) // Example: "pkg.Type"
		default:
			fmt.Printf("Unknown Type: %T\n", d)
		}
	}

	return hasil, variable
}

func printExprTypeInner(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return printExprTypeInner(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

func (gh *ginHandler) Ast() *ast.FuncLit {
	gh.parameterAst()
	declaration, varnames := gh.prepareDeclAst()
	varres, _ := varnames[0].(*ast.Ident)
	resname := varres.Name

	if resname == "err" {
		resname = ""
	}

	declaration = append(declaration, // err = srv.CreateUser()
		&ast.AssignStmt{
			Lhs: varnames,
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("srv"),
						Sel: ast.NewIdent(gh.fnname),
					},
				},
			},
		},
		// if err != nil { ctx.AbortWithStatusJSON(...) }
		standardErrorStmt(),
		standardRes(resname),
	)

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
		Body: &ast.BlockStmt{
			List: declaration,
		},
	}

	return fn
}
