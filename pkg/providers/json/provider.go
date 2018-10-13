package json

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/craiggwilson/go-typeproviders/pkg/internal/bsonutil"
	"github.com/craiggwilson/go-typeproviders/pkg/structbuilder"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Config holds information required for configuration mongodb.
type Config struct {
	StructName string
	Input      io.Reader
}

// NewStructProvider makes a StructProvider.
func NewStructProvider(cfg Config) *StructProvider {
	return &StructProvider{
		cfg: cfg,
	}
}

// StructProvider provides structs.
type StructProvider struct {
	cfg Config
}

// ProvideStructs implements the generators.StructProvider interface.
func (p *StructProvider) ProvideStructs(ctx context.Context) ([]*structbuilder.Struct, error) {
	jsonString, err := ioutil.ReadAll(p.cfg.Input)
	if err != nil {
		return nil, err
	}
	doc, err := bson.ParseExtJSONObject(string(jsonString))
	if err != nil && err != io.EOF {
		return nil, err
	}

	tb := bsonutil.NewTypeBuilder()
	tb.IncludeDocument(doc)

	result := bsonutil.BuildStruct(p.cfg.StructName, tb)
	return []*structbuilder.Struct{result}, nil
}
