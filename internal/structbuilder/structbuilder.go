package structbuilder

import (
	"io"

	"github.com/alecthomas/template"
)

// StructBuilder helps build structs.
type StructBuilder struct {
	Name string
}

// Write writes the struct definition to the writer.
func (b *StructBuilder) Write(w io.Writer) error {
	if err := structTmpl.Execute(w, b); err != nil {
		return err
	}

	return nil
}

var structTmpl = template.Must(template.New("struct").Parse(`
type {{ .Name }} struct {
}
`))
