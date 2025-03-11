package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
)

var identifier = "bertaut_api:"

func TestService(fs *token.FileSet, pkg *doc.Package, t *doc.Type) error {
	var err error
	reg := NewRegisterFile(fs, pkg, t)
	log.Println(t.Methods, t.Consts)

	for _, fn := range t.Methods {
		log.Println(fn.Doc, "asdasd")
		// log.Println(fn.Name)
	}

	err = reg.
		Initialize().
		RegisterFunc(func() error {

			return nil
		}).
		MemberFunc().
		Write()

	return err
}

type RegErr struct {
	Msg  string
	Path string
}

// Error implements error.
func (r *RegErr) Error() string {
	data, _ := json.Marshal(r)
	return string(data)
}

type RegisterFile struct {
	fs  *token.FileSet
	t   *doc.Type
	pkg *doc.Package
	ast *ast.File
	err error

	base string
}

func (r *RegisterFile) Initialize() *RegisterFile {
	if r.err != nil {
		return r
	}

	base, err := r.baseRoute()
	if err != nil {
		return r.setErr(err)
	}

	r.base = base
	return r
}

func (r *RegisterFile) MemberFunc() *RegisterFile {

	if iface, ok := r.t.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType); ok {
		for _, d := range iface.Methods.List {
			log.Println(d.Doc.List[0], "asdasd")
		}
	}

	return r
}

func (r *RegisterFile) RegisterFunc(handle func() error) *RegisterFile {
	if r.err != nil {
		return r
	}

	funcname := "Register" + r.t.Name + "Api"

	handleCall := &ast.CallExpr{
		Fun: &ast.SelectorExpr{ // g.Handle
			X:   ast.NewIdent("g"),
			Sel: ast.NewIdent("Handle"),
		},
		Args: []ast.Expr{
			// http.MethodPost
			&ast.SelectorExpr{
				X:   ast.NewIdent("http"),
				Sel: ast.NewIdent("MethodPost"),
			},
			// "test"
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: `"test"`,
			},
			// func(ctx *gin.Context) {}
			&ast.FuncLit{
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("ctx")},
								Type: &ast.StarExpr{ // *gin.Context
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("gin"),
										Sel: ast.NewIdent("Context"),
									},
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{}}, // Empty function body
			},
		},
	}

	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(funcname), // Function name
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("srv")}, // Parameter name: usr
						Type:  ast.NewIdent(r.t.Name),            // Parameter type: UserApi
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
		Body: &ast.BlockStmt{List: []ast.Stmt{
			&ast.ExprStmt{X: handleCall},
		}},
	}

	r.ast.Decls = append(r.ast.Decls, funcDecl)
	return r

}

func (r *RegisterFile) Filename() string {
	pos := r.fs.Position(r.t.Decl.Pos())
	return strings.Replace(pos.Filename, ".go", "_api_gen.go", 1)
}

func (r *RegisterFile) Write() error {
	if r.err != nil {
		return r.err
	}

	filename := r.Filename()
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fs := token.NewFileSet()
	if err := format.Node(file, fs, r.ast); err != nil {
		fmt.Println("Error formatting AST:", err)
		return err
	}

	return nil
}

func (r *RegisterFile) genError(msg string) error {
	pos := r.fs.Position(r.t.Decl.Pos())
	return &RegErr{
		Msg:  msg,
		Path: pos.Filename,
	}
}

func (r *RegisterFile) baseRoute() (string, error) {
	raws := strings.Split(r.t.Doc, "\n")
	for _, raw := range raws {
		if !strings.Contains(raw, identifier) {
			continue
		}

		bases := strings.Split(raw, ":")
		if len(bases) < 2 {
			return "", r.genError("raw")
		}
		base := bases[1]
		re := regexp.MustCompile(`\s+`) // Matches spaces, tabs, newlines
		base = re.ReplaceAllString(base, "")

		return base, nil

	}

	return "", nil
}

func (r *RegisterFile) setErr(err error) *RegisterFile {
	if r.err != nil {
		return r
	}

	if err != nil {
		r.err = err
	}

	return r
}

func NewRegisterFile(fs *token.FileSet, pkg *doc.Package, t *doc.Type) *RegisterFile {

	return &RegisterFile{
		fs:  fs,
		t:   t,
		pkg: pkg,
		ast: &ast.File{
			Name: ast.NewIdent(pkg.Name),
			Decls: []ast.Decl{
				&ast.GenDecl{
					Tok: token.IMPORT,
					Specs: []ast.Spec{
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"github.com/gin-gonic/gin\"",
							},
						},
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"net/http\"",
							},
						},
					},
				},
			},
		},
	}
}
