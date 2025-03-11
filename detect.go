package main

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

func DetectService(
	includedPackage []string,
	directoryPath string,
	handler func(fs *token.FileSet, pkg *doc.Package, t *doc.Type) error,
) error {
	fs := token.NewFileSet()
	nodes, err := parser.ParseDir(fs, directoryPath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return err
	}

	for packageName, node := range nodes {
		fmt.Printf("Inspecting %s \n", packageName)

		pkg := doc.New(node, directoryPath, doc.AllDecls)

		for _, d := range pkg.Types {
			t := d

			if strings.Contains(t.Doc, "bertaut_api:") {
				err := handler(fs, pkg, t)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}
