package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func fieldToVar(listField []*ast.Field) ([]ast.Stmt, []ast.Expr) {
	hasil := []ast.Stmt{}
	variable := []ast.Expr{}

	c := 0
	resname := func() string {
		c += 1
		return fmt.Sprintf("param%d", c)
	}

	for _, d := range listField {
		switch t := d.Type.(type) {
		case *ast.StructType:
			fmt.Println("Struct {")
			for _, field := range t.Fields.List {
				for _, fieldName := range field.Names {
					fmt.Printf("  %s: %s\n", fieldName.Name, printExprTypeInner(field.Type))
				}
			}
			fmt.Println("}")
		case *ast.StarExpr: // untuk pointer
			var name string
			if len(d.Names) == 0 {
				name = resname()
			} else {
				name = d.Names[0].Name
			}

			varname := &ast.UnaryExpr{ // Represents the `&` operator
				Op: token.AND,              // `&` operator
				X:  &ast.Ident{Name: name}, // Argument: param1
			}

			hasil = append(hasil, &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{ast.NewIdent(name)},
							Type:  t.X,
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
		case *ast.Ident:
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
		default:
			fmt.Printf("Unknown Type: %T\n", d.Type)
		}
	}

	return hasil, variable

}

func getName(t any) string {
	switch d := t.(type) {
	case *ast.UnaryExpr:
		return getName(d.X)
	case *ast.Ident:
		return d.Name
	default:
		return ""

	}
}

func formatImportFile(fpath string) error {
	cmd := exec.Command("goimports", "-w", fpath)
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	log.Println("formatting", fpath, string(out))
	return nil
}

func convertCamelToSnake(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`) // Match lowercase + uppercase transition
	snake := re.ReplaceAllString(s, `${1}_${2}`)  // Insert underscore between matches
	return strings.ToLower(snake)                 // Convert to lowercase
}

func ginBindQuery(arg ast.Expr) *ast.AssignStmt {
	assignStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{Name: "err"}, // Left-hand side (LHS): err
		},
		Tok: token.ASSIGN, // Assignment operator: =
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "ctx"},       // Receiver: ctx
					Sel: &ast.Ident{Name: "BindQuery"}, // Method: BindQuery
				},
				Args: []ast.Expr{
					arg,
				},
			},
		},
	}

	return assignStmt
}

func ginBindJson(arg ast.Expr) *ast.AssignStmt {
	assignStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{Name: "err"}, // Left-hand side (LHS): err
		},
		Tok: token.ASSIGN, // Assignment operator: =
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "ctx"},      // Receiver: ctx
					Sel: &ast.Ident{Name: "BindJSON"}, // Method: BindQuery
				},
				Args: []ast.Expr{
					arg,
				},
			},
		},
	}

	return assignStmt
}
