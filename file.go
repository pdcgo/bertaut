package bertaut

import (
	"context"
	"go/ast"
	"go/token"
)

type KeyCtx string

const (
	IMPORT     KeyCtx = "import"
	FUNCPARAMS KeyCtx = "funcparams"
)

func CreateFile(pkgName string, handle func(ctx context.Context) []ast.Decl) *ast.File {
	ctx := context.TODO()

	importdec := &ast.GenDecl{
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
	}

	ctx = context.WithValue(ctx, IMPORT, importdec)
	declaration := []ast.Decl{
		importdec,
	}

	declaration = append(declaration, handle(ctx)...)

	return &ast.File{
		Name:  ast.NewIdent(pkgName),
		Decls: declaration,
	}

}
