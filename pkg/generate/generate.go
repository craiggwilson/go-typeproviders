package generate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/craiggwilson/go-typeproviders/pkg/structbuilder"
)

// StructProvider is the interface that wraps the ProvideStructs method.
type StructProvider interface {
	// ProvideStructs builds all the structs that should be included. The filename
	// may refer to a file that will be used and, if it exists, should be parsed
	// for existing structures and amended.
	ProvideStructs(filename string) ([]*structbuilder.StructBuilder, error)
}

// Generate uses the struct provider to generate and write code to the provided
// filename.
func Generate(p StructProvider, filename string, pkg string) error {
	structs, err := p.ProvideStructs(filename)
	if err != nil {
		return err
	}

	data := struct {
		Structs []*structbuilder.StructBuilder
		Package string
	}{
		Structs: structs,
		Package: pkg,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, &data); err != nil {
		return err
	}

	if filename != "" {
		return ioutil.WriteFile(filename, buf.Bytes(), 0666)
	}

	fmt.Println(buf.String())
	return nil
}

var tmpl = template.Must(template.New("file").Parse(`/*
* CODE GENERATED AUTOMATICALLY WITH github.com/craiggwilson/go-typeproviders
* THIS FILE SHOULD NOT BE EDITED BY HAND
*/
{{define "struct"}}
type {{ .Name }} struct {
	{{range .Fields}}
		{{.Name}} {{.Type}} {{.QuotedTags}}
	{{end}}
}
{{end}}

package {{.Package}}
{{range .Structs}}
{{template "struct" .}}
{{end}}`))
