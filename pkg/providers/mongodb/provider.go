package mongodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/craiggwilson/go-typeproviders/pkg/naming"

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
		return "time.Time time"
	case bson.TypeDecimal128:
		return "decimal128.Decimal128 github.com/mongodb/mongo-go-driver/bson/decimal128"
	case bson.TypeDouble:
		return "float64"
	case bson.TypeInt32:
		return "int32"
	case bson.TypeInt64:
		return "int64"
	case bson.TypeObjectID:
		return "objectid.ObjectID github.com/mongodb/mongo-go-driver/bson/objectid"
	case bson.TypeString:
		return "string"
	case bson.TypeTimestamp:
		return "time.Time time"
	default:
		return "*bson.Value github.com/mongodb/mongo-go-driver/bson"
	}
}

func buildStruct(name string, sb *structbuilder.StructBuilder) (*structbuilder.Struct, error) {
	s := structbuilder.Struct{
		Name: naming.Struct(name),
	}

	for _, f := range sb.Fields() {
		t := chooseType(f.Types)
		parts := strings.SplitN(t.Name(), " ", 2)
		typeName := parts[0]
		importPath := ""
		if len(parts) == 2 {
			importPath = parts[1]
		}

		fieldName := naming.ExportedField(f.Name)
		if typeName == "array" {
			fieldName = naming.Pluralize(fieldName)
			typeName = "[]" + typeName
		}

		s.Fields = append(s.Fields, structbuilder.Field{
			Name: fieldName,
			Type: structbuilder.FieldType{
				Name:       typeName,
				ImportPath: importPath,
			},
		})
	}

	return &s, nil
}

func chooseType(types []structbuilder.Type) structbuilder.Type {
	return types[0]
}
