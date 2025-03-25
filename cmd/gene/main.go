package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/doc"
	"log"
	"strings"

	"github.com/pdcgo/bertaut"
)

func PackageProcess(packageName string, pkg *doc.Package) error {
	file := bertaut.CreateFile(packageName, func(ctx context.Context) []ast.Decl {
		filedecl := []ast.Decl{}

		for _, d := range pkg.Types {
			t := d
			if strings.Contains(t.Doc, "bertaut_api:") {
				regdecl := bertaut.RegisterApi(ctx, t.Name, func(ctx context.Context) []ast.Stmt { // fungsi register api

					hasil := []ast.Stmt{}
					bertaut.IterateMethod(t, func(method, uri string) {
						log.Println(method)
						log.Println(uri)

						callexpr := bertaut.RegisterGin(ctx, method, uri, func(ctx context.Context) {

						})

						hasil = append(hasil, &ast.ExprStmt{X: callexpr})
					})

					// handler := bertaut.RegisterGin(
					// 	ctx,
					// 	"MethodPost",
					// 	"/v2/mock",
					// 	func(ctx context.Context) {

					// 	},
					// )

					return hasil
				})

				filedecl = append(filedecl, regdecl)
			}
		}

		return filedecl
	})

	bertaut.SaveAst(fmt.Sprintf("./%s/v2_api_gen.go", packageName), file)

	return nil
}

func main() {

	bertaut.DetectService(
		[]string{
			"mock_http",
		},
		"./mock_http",
		PackageProcess,
	)

	// err := DetectService(
	// 	[]string{
	// 		"mock_http",
	// 	},
	// 	"./mock_http",
	// 	TestService,
	// )
	// if err != nil {
	// 	panic(err)
	// }

}

// func createHandler(packageName string, writer io.Writer) error {
// 	file := &ast.File{
// 		Name: ast.NewIdent(packageName), // package main
// 		Decls: []ast.Decl{
// 			&ast.GenDecl{ // Import statement: import "github.com/gin-gonic/gin"
// 				Tok: token.IMPORT,
// 				Specs: []ast.Spec{
// 					&ast.ImportSpec{
// 						Path: &ast.BasicLit{
// 							Kind:  token.STRING,
// 							Value: "\"github.com/gin-gonic/gin\"",
// 						},
// 					},
// 					&ast.ImportSpec{
// 						Path: &ast.BasicLit{
// 							Kind:  token.STRING,
// 							Value: "\"net/http\"",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	handleCall := &ast.CallExpr{
// 		Fun: &ast.SelectorExpr{ // g.Handle
// 			X:   ast.NewIdent("g"),
// 			Sel: ast.NewIdent("Handle"),
// 		},
// 		Args: []ast.Expr{
// 			// http.MethodPost
// 			&ast.SelectorExpr{
// 				X:   ast.NewIdent("http"),
// 				Sel: ast.NewIdent("MethodPost"),
// 			},
// 			// "test"
// 			&ast.BasicLit{
// 				Kind:  token.STRING,
// 				Value: `"test"`,
// 			},
// 			// func(ctx *gin.Context) {}
// 			&ast.FuncLit{
// 				Type: &ast.FuncType{
// 					Params: &ast.FieldList{
// 						List: []*ast.Field{
// 							{
// 								Names: []*ast.Ident{ast.NewIdent("ctx")},
// 								Type: &ast.StarExpr{ // *gin.Context
// 									X: &ast.SelectorExpr{
// 										X:   ast.NewIdent("gin"),
// 										Sel: ast.NewIdent("Context"),
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Body: &ast.BlockStmt{List: []ast.Stmt{}}, // Empty function body
// 			},
// 		},
// 	}

// 	funcDecl := &ast.FuncDecl{
// 		Name: ast.NewIdent("UserRegister"), // Function name
// 		Type: &ast.FuncType{
// 			Params: &ast.FieldList{
// 				List: []*ast.Field{
// 					{
// 						Names: []*ast.Ident{ast.NewIdent("usr")}, // Parameter name: usr
// 						Type:  ast.NewIdent("UserService"),       // Parameter type: UserApi
// 					},
// 					{
// 						Names: []*ast.Ident{ast.NewIdent("g")}, // Parameter name: g
// 						Type: &ast.StarExpr{ // Pointer type: *gin.RouterGroup
// 							X: &ast.SelectorExpr{
// 								X:   ast.NewIdent("gin"),         // Package name: gin
// 								Sel: ast.NewIdent("RouterGroup"), // Struct name: RouterGroup
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		Body: &ast.BlockStmt{List: []ast.Stmt{
// 			&ast.ExprStmt{X: handleCall},
// 		}},
// 	}
// 	file.Decls = append(file.Decls, funcDecl)

// 	// Write formatted AST to file
// 	fs := token.NewFileSet()
// 	if err := format.Node(writer, fs, file); err != nil {
// 		return err
// 	}

// 	return nil
// }
