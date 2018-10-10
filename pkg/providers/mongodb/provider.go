package mongodb

import (
	"context"
	"fmt"

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
				bson.EC.Int64("size", 1),
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

		tb := mapStruct(doc)
		if sb == nil {
			sb = tb
		} else {
			sb.Merge(tb)
		}
	}

	s, err := buildStruct(coll.Name(), sb)
	if err != nil {
		return nil, err
	}

	return []*structbuilder.Struct{s}, nil
}

func mapStruct(doc *bson.Document) *structbuilder.StructBuilder {
	t := structbuilder.NewStructBuilder("struct")

	iter := doc.Iterator()
	for iter.Next() {
		e := iter.Element()
		vt := mapValue(e.Value())
		t.Include(e.Key(), vt)
	}

	return t
}

func mapArray(arr *bson.Array) *structbuilder.ArrayBuilder {
	t := structbuilder.NewArrayBuilder("array")
	iter, _ := arr.Iterator()
	for iter.Next() {
		vt := mapValue(iter.Value())
		t.Include(vt)
	}

	return t
}

func mapValue(v *bson.Value) structbuilder.Type {
	switch v.Type() {
	case bson.TypeArray:
		return mapArray(v.MutableArray())
	case bson.TypeEmbeddedDocument:
		return mapStruct(v.MutableDocument())
	default:
		return structbuilder.NewPrimitiveBuilder(mapPrimitiveTypeName(v.Type()))
	}
}

func mapPrimitiveTypeName(t bson.Type) string {
	switch t {
	case bson.TypeBinary:
		return "[]byte"
	case bson.TypeBoolean:
		return "bool"
	case bson.TypeDateTime:
		return "time.Time"
	case bson.TypeDecimal128:
		return "decimal128.Decimal128"
	case bson.TypeDouble:
		return "float64"
	case bson.TypeInt32:
		return "int32"
	case bson.TypeInt64:
		return "int64"
	case bson.TypeObjectID:
		return "objectid.ObjectID"
	case bson.TypeString:
		return "string"
	case bson.TypeTimestamp:
		return "time.Time"
	default:
		return "*bson.Value"
	}
}

func buildStruct(name string, b *structbuilder.StructBuilder) (*structbuilder.Struct, error) {
	s := structbuilder.Struct{
		Name: name,
	}

	for _, f := range b.Fields() {
		s.Fields = append(s.Fields, structbuilder.Field{
			Name: f.Name,
			Type: f.Types[0].Name(),
		})
	}

	return &s, nil
}
