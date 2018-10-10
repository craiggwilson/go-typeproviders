package generate

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io/ioutil"
	"sort"
	"text/template"

	"github.com/craiggwilson/go-typeproviders/pkg/structbuilder"
)

// StructProvider is the interface that wraps the ProvideStructs method.
type StructProvider interface {
	// ProvideStructs builds all the structs that should be included. The filename
	// may refer to a file that will be used and, if it exists, should be parsed
	// for existing structures and amended.
	ProvideStructs(ctx context.Context, filename string) ([]*structbuilder.Struct, error)
}

// Generate uses the struct provider to generate and write code to the provided
// filename.
func Generate(ctx context.Context, p StructProvider, filename string, pkg string) error {
	structs, err := p.ProvideStructs(ctx, filename)
	if err != nil {
		return err
	}

	importPaths := uniqueImportPaths(structs)

	data := struct {
		Structs     []*structbuilder.Struct
		Package     string
		ImportPaths []string
	}{
		structs,
		pkg,
		importPaths,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, &data); err != nil {
		return err
	}

	//formatted := buf.Bytes()
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if filename != "" {
		return ioutil.WriteFile(filename, formatted, 0666)
	}

	fmt.Println(string(formatted))
	return nil
}

func uniqueImportPaths(structs []*structbuilder.Struct) []string {
	set := make(map[string]struct{})
	var results []string
	for _, s := range structs {
		for _, f := range s.Fields {
			if f.Type.ImportPath != "" {
				if _, ok := set[f.Type.ImportPath]; !ok {
					set[f.Type.ImportPath] = struct{}{}
					results = append(results, f.Type.ImportPath)
				}
			}
		}
	}

	sort.Strings(results)
	return results
}

var tmpl = template.Must(template.New("file").Parse(`/*
* CODE GENERATED AUTOMATICALLY WITH github.com/craiggwilson/go-typeproviders
* THIS FILE SHOULD NOT BE EDITED BY HAND
*/
{{define "struct"}}
type {{ .Name }} struct {
	{{range .Fields}}
	{{.Name}} {{.Type.Name}} {{.QuotedTags}}
	{{- end}}
}
{{end}}

package {{.Package}}

import (
	{{range .ImportPaths}}
		"{{.}}"
	{{- end}}
)
{{range .Structs}}
{{template "struct" .}}
{{end}}`))
