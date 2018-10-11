package bson

import (
	"context"
	"io"

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
func (p *StructProvider) ProvideStructs(ctx context.Context, filename string) ([]*structbuilder.Struct, error) {
	var sb *structbuilder.StructBuilder
	for {
		doc := bson.NewDocument()
		_, err := doc.ReadFrom(p.cfg.Input)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		tb := bsonutil.DocumentToStructBuilder(doc)
		if sb == nil {
			sb = tb
		} else {
			sb.Merge(tb)
		}
	}

	results, err := bsonutil.BuildStructs(p.cfg.StructName, sb, false)
	if err != nil {
		return nil, err
	}

	return results, nil
}
