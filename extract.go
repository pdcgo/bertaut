package bertaut

import (
	"go/ast"
	"go/doc"
	"net/url"
	"regexp"
	"strings"
)

var identifier = "bertaut_api:"

func IterateMethod(t *doc.Type, handle func(method, uri string)) {
	if iface, ok := t.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType); ok {

		baseuri := extractBaseUri(t)

		for _, met := range iface.Methods.List {
			httpMethod := extractMethod(met.Doc.List)
			method := mapHttpMethod(httpMethod)
			uri := extractUri(baseuri, met)
			if method == "" {
				continue
			}

			handle(method, uri)

		}
	}
}

func extractBaseUri(t *doc.Type) string {
	raws := strings.Split(t.Doc, "\n")
	for _, raw := range raws {
		if !strings.Contains(raw, identifier) {
			continue
		}

		bases := strings.Split(raw, ":")
		if len(bases) < 2 {
			return ""
		}
		base := bases[1]
		re := regexp.MustCompile(`\s+`) // Matches spaces, tabs, newlines
		base = re.ReplaceAllString(base, "")

		return base

	}

	return ""

}

func extractUri(baseuri string, met *ast.Field) string {
	uripath, _ := url.JoinPath(baseuri, convertCamelToSnake(met.Names[0].Name))
	return uripath
}

func mapHttpMethod(method string) string {
	methodfunc := "MethodGet"
	switch method {
	case "post":
		methodfunc = "MethodPost"
	case "get":
		methodfunc = "MethodGet"
	case "put":
		methodfunc = "MethodPut"
	case "delete":
		methodfunc = "MethodDelete"
	}

	return methodfunc
}

func extractMethod(coms []*ast.Comment) string {
	for _, com := range coms {
		if !strings.HasPrefix(com.Text, `// method:`) {
			continue
		}

		method := strings.ReplaceAll(com.Text, `// method: `, "")
		method = strings.TrimSpace(method)
		return method
	}

	return ""
}

func convertCamelToSnake(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`) // Match lowercase + uppercase transition
	snake := re.ReplaceAllString(s, `${1}_${2}`)  // Insert underscore between matches
	return strings.ToLower(snake)                 // Convert to lowercase
}
