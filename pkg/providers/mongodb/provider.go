package mongodb

import (
	"github.com/craiggwilson/go-typeproviders/pkg/structbuilder"
)

// Config holds information required for configuration mongodb.
type Config struct {
	URI            string
	DatabaseName   string
	CollectionName string
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
func (p *StructProvider) ProvideStructs(filename string) ([]*structbuilder.StructBuilder, error) {
	// client, err := mongo.Connect(cfg.URI)
	// if err != nil {
	// 	return err
	// }

	structs := []*structbuilder.StructBuilder{
		{
			Name: "Temp",
		},
	}

	return structs, nil
}
