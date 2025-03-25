package bertaut

import (
	"fmt"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strings"
)

func DetectService(
	includedPackage []string,
	directoryPath string,
	handler func(packageName string, pkg *doc.Package) error,
) error {
	fs := token.NewFileSet()
	nodes, err := parser.ParseDir(fs, directoryPath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return err
	}

	for packageName, node := range nodes {
		found := false

		fmt.Printf("Inspecting %s \n", packageName)
		pkg := doc.New(node, directoryPath, doc.AllDecls)

	CC:
		for _, d := range pkg.Types {
			t := d

			if strings.Contains(t.Doc, "bertaut_api:") {
				found = true
				break CC
			}
		}
		if found {
			err = handler(packageName, pkg)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func SaveAst(filename string, node any) error {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fs := token.NewFileSet()
	if err := format.Node(file, fs, node); err != nil {
		fmt.Println("Error formatting AST:", err)
		return err
	}

	if err := formatImportFile(filename); err != nil {
		return err
	}

	return nil
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
