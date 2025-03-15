package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

type ginHandler struct {
	ft     *ast.FuncType
	fnname string
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
			varname := ast.NewIdent(resname())
			hasil = append(hasil, &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{varname},
							Type:  t,
						},
					},
				},
			})
			variable = append(variable, varname)
		case *ast.SelectorExpr:
			varname := ast.NewIdent(resname())
			hasil = append(hasil, &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{varname},
							Type:  ast.NewIdent(t.Sel.Name),
						},
					},
				},
			})
			variable = append(variable, varname)

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

	declaration, varnames := gh.prepareDeclAst()
	varres, _ := varnames[0].(*ast.Ident)
	resname := varres.Name

	if resname == "err" {
		resname = ""
	}

	decsrvparam, srvparam := fieldToVar(gh.ft.Params.List) // parameter untuk serivce
	declaration = append(declaration, decsrvparam...)

	for _, d := range srvparam {
		name := getName(d)
		switch name {
		case "query":
			declaration = append(declaration,
				ginBindQuery(d),
				standardErrorStmt(),
			)
		case "payload":
			declaration = append(declaration,
				ginBindJson(d),
				standardErrorStmt(),
			)

		case "identity":
			declaration = append(declaration,
				buildFromCtx("identity"),
				standardErrorStmt(),
			)

		default:
			declaration = append(declaration,
				buildFromCtx(name),
				standardErrorStmt(),
			)
		}

	}

	declaration = append(declaration,
		// ginBindQuery(srvparam[0]),
		&ast.AssignStmt{ // err = srv.CreateUser()
			Lhs: varnames,
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("srv"),
						Sel: ast.NewIdent(gh.fnname),
					},
					Args: srvparam,
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

func buildFromCtx(callname string) *ast.AssignStmt {
	assignStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{Name: "err"}, // Left-hand side (LHS): err
		},
		Tok: token.ASSIGN, // Assignment operator: =
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: callname},       // Receiver: identity
					Sel: &ast.Ident{Name: "BuildFromCtx"}, // Method: BuildFromCtx
				},
				Args: []ast.Expr{
					&ast.Ident{Name: "ctx"}, // Argument: ctx
				},
			},
		},
	}

	return assignStmt
}

func documentationCall(uri string) *ast.IfStmt {
	ifStmt := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  &ast.Ident{Name: "doc"},
			Op: token.NEQ, // "!=" operator
			Y:  &ast.Ident{Name: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// err = doc(http.MethodGet, g.BasePath()+"/users/item", nil, nil)
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "err"}},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: "doc"},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   &ast.Ident{Name: "http"},
									Sel: &ast.Ident{Name: "MethodGet"},
								},
								&ast.BinaryExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "g"},
											Sel: &ast.Ident{Name: "BasePath"},
										},
									},
									Op: token.ADD,
									Y:  &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, uri)},
								},
								&ast.Ident{Name: "nil"},
								&ast.Ident{Name: "nil"},
							},
						},
					},
				},
				// if err != nil { panic(err) }
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ, // "!=" operator
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun:  &ast.Ident{Name: "panic"},
									Args: []ast.Expr{&ast.Ident{Name: "err"}},
								},
							},
						},
					},
				},
			},
		},
	}

	return ifStmt
}
