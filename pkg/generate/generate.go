package generate

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io/ioutil"
	"sort"
	"strings"
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
func Generate(ctx context.Context, p StructProvider, filename string, pkg string, embedStructs bool) error {
	structs, err := p.ProvideStructs(ctx, filename)
	if err != nil {
		return err
	}

	if !embedStructs {
		var results []*structbuilder.Struct
		for _, s := range structs {
			results = append(results, s.UnembedStructs()...)
		}

		structs = results
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

var tmpl = template.Must(template.New("file").Funcs(template.FuncMap{
	"brackets": func(count int) string {
		return strings.Repeat("[]", count)
	},
	"canBeNull": func(canBeNull bool) string {
		if canBeNull {
			return "*"
		}

		return ""
	},
	"quotedTags": func(tags []string) string {
		if len(tags) > 0 {
			return "`" + strings.Join(tags, " ") + "`"
		}

		return ""
	},
}).Parse(`/* CODE GENERATED AUTOMATICALLY WITH github.com/craiggwilson/go-typeproviders */
{{define "embeddedStruct" -}}
struct {
	{{range .Fields}}
	{{.Name}} {{brackets .Type.ArrayCount }} {{canBeNull .Type.CanBeNull }} {{if .Type.EmbeddedStruct }}{{template "embeddedStruct" .Type.EmbeddedStruct }} {{else}} {{.Type.Name}} {{end}} {{quotedTags .Tags }}
	{{end}}
}
{{- end}}

{{define "struct" -}}
type {{ .Name }} {{template "embeddedStruct" .}}
{{end}}

package {{.Package}}

{{if (len .ImportPaths) eq 0}}
import (
	{{range .ImportPaths}}
		"{{.}}"
	{{- end}}
)
{{end}}

{{range .Structs}}
{{template "struct" .}}
{{end}}`))
