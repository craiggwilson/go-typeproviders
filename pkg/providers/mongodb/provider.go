package mongodb

import (
	"context"
	"fmt"

	"github.com/craiggwilson/go-typeproviders/pkg/internal/bsonutil"
	"github.com/craiggwilson/go-typeproviders/pkg/structbuilder"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Config holds information required for configuration mongodb.
type Config struct {
	URI            string
	DatabaseName   string
	CollectionName string
	SampleSize     uint
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
	client, err := mongo.Connect(ctx, p.cfg.URI)
	if err != nil {
		return nil, err
	}

	db := client.Database(p.cfg.DatabaseName)

	if p.cfg.CollectionName == "" {
		return p.provideFromDatabase(ctx, db)
	}

	coll := db.Collection(p.cfg.CollectionName)
	return p.provideFromCollection(ctx, coll)
}

func (p *StructProvider) provideFromDatabase(ctx context.Context, db *mongo.Database) ([]*structbuilder.Struct, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *StructProvider) provideFromCollection(ctx context.Context, coll *mongo.Collection) ([]*structbuilder.Struct, error) {
	pipeline := bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements(
				"$sample",
				bson.EC.Int64("size", int64(p.cfg.SampleSize)),
			),
		),
	)
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var sb *structbuilder.StructBuilder

	for cursor.Next(ctx) {
		doc := bson.NewDocument()
		err := cursor.Decode(doc)
		if err != nil {
			return nil, err
		}

		tb := bsonutil.DocumentToStructBuilder(doc)
		if sb == nil {
			sb = tb
		} else {
			sb.Merge(tb)
		}
	}

	s, err := bsonutil.BuildStruct(coll.Name(), sb)
	if err != nil {
		return nil, err
	}

	return []*structbuilder.Struct{s}, nil
}
