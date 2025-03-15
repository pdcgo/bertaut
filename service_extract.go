package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/token"
	"os"
	"regexp"
	"strings"
)

var identifier = "bertaut_api:"

func TestService(fs *token.FileSet, pkg *doc.Package, t *doc.Type) error {
	var err error
	reg := NewRegisterFile(fs, pkg, t)

	fmt.Printf("generating %s\n", t.Name)

	err = reg.
		Initialize().
		RegisterFunc(func() ([]*ast.CallExpr, error) {
			var err error
			var hasil []*ast.CallExpr

			if iface, ok := t.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType); ok {
				reg := ginCallHandle{
					iface:     iface,
					ifaceName: t.Name,
					baseuri:   reg.base,
				}

				hasil, err = reg.registerGinHandlers()
			}

			return hasil, err

		}).
		MemberFunc().
		Write()

	return err
}

type ParserCtx struct {
	fs  *token.FileSet
	t   *doc.Type
	pkg *doc.Package
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
	ctx *ParserCtx
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

	return r
}

func (r *RegisterFile) RegisterFunc(handle func() ([]*ast.CallExpr, error)) *RegisterFile {
	if r.err != nil {
		return r
	}

	stmts, err := handle()
	if err != nil {
		return r.setErr(err)
	}

	body := []ast.Stmt{}
	body = append(body, &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{ast.NewIdent("err")},
					Type:  ast.NewIdent("error"),
				},
			},
		},
	})

	for _, stmt := range stmts {
		body = append(body, &ast.ExprStmt{X: stmt})
		body = append(body, documentationCall("/users"))
	}

	funcname := "Register" + r.ctx.t.Name + "Api"

	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(funcname), // Function name
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("srv")}, // Parameter name: usr
						Type:  ast.NewIdent(r.ctx.t.Name),        // Parameter type: UserApi
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
					{
						Names: []*ast.Ident{{Name: "doc"}},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{{Name: "method"}},
										Type:  &ast.Ident{Name: "string"},
									},
									{
										Names: []*ast.Ident{{Name: "path"}},
										Type:  &ast.Ident{Name: "string"},
									},
									{
										Names: []*ast.Ident{{Name: "query"}},
										Type:  &ast.Ident{Name: "any"},
									},
									{
										Names: []*ast.Ident{{Name: "payload"}},
										Type:  &ast.Ident{Name: "any"},
									},
								},
							},
							Results: &ast.FieldList{
								List: []*ast.Field{
									{
										Type: &ast.Ident{Name: "error"},
									},
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{List: body},
	}

	r.ast.Decls = append(r.ast.Decls, funcDecl)
	return r

}

func (r *RegisterFile) Filename() string {
	pos := r.ctx.fs.Position(r.ctx.t.Decl.Pos())
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

	if err := formatImportFile(filename); err != nil {
		return err
	}

	return nil
}

func (r *RegisterFile) genError(msg string) error {
	pos := r.ctx.fs.Position(r.ctx.t.Decl.Pos())
	return &RegErr{
		Msg:  msg,
		Path: pos.Filename,
	}
}

func (r *RegisterFile) baseRoute() (string, error) {
	raws := strings.Split(r.ctx.t.Doc, "\n")
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
		ctx: &ParserCtx{
			fs:  fs,
			t:   t,
			pkg: pkg,
		},
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
