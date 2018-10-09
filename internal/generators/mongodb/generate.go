package mongodb

import (
	"io"

	"github.com/alecthomas/template"
	"github.com/craiggwilson/typeproviders/internal/structbuilder"
)

// NewGenerator makes a Generator.
func NewGenerator(cfg Config) *Generator {
	return &Generator{
		Cfg: cfg,
	}
}

// Generator handles generating struct for the provided config.
type Generator struct {
	Cfg Config

	Structs []*structbuilder.StructBuilder
}

// Write generates the structs and writes them to the writer.
func (g *Generator) Write(w io.Writer) error {
	if err := g.buildStructs(); err != nil {
		return err
	}

	if err := tmpl.Execute(w, g); err != nil {
		return err
	}

	return nil
}

func (g *Generator) buildStructs() error {
	// client, err := mongo.Connect(cfg.URI)
	// if err != nil {
	// 	return err
	// }

	g.Structs = append(g.Structs, &structbuilder.StructBuilder{Name: "Temp"})

	return nil
}

var tmpl = template.Must(template.New("collection").Parse(`/*
* CODE GENERATED AUTOMATICALLY WITH github.com/craiggwilson/go-typeproviders
* THIS FILE SHOULD NOT BE EDITED BY HAND
*/

package {{.Cfg.Package}}
`))
