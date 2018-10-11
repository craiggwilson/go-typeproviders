package bsonutil

import (
	"fmt"
	"strings"

	"github.com/craiggwilson/go-typeproviders/pkg/naming"
	"github.com/craiggwilson/go-typeproviders/pkg/structbuilder"
	"github.com/mongodb/mongo-go-driver/bson"
)

// DocumentToStructBuilder maps a document to a structbuilder.StructBuilder.
func DocumentToStructBuilder(doc *bson.Document) *structbuilder.StructBuilder {
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
		return DocumentToStructBuilder(v.MutableDocument())
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

// BuildStruct builds a structbuilder.Struct from a structbuilder.StructBuilder.
func BuildStruct(name string, sb *structbuilder.StructBuilder) (*structbuilder.Struct, error) {
	s := structbuilder.Struct{
		Name: naming.Struct(name),
	}

	for _, f := range sb.Fields() {
		typeName, importPath := selectType(f.Types())
		fieldName := naming.ExportedField(f.Name())
		if strings.HasPrefix(typeName, "[]") {
			fieldName = naming.Pluralize(fieldName)
		}

		s.Fields = append(s.Fields, structbuilder.Field{
			Name: fieldName,
			Tags: []string{
				fmt.Sprintf(`"bson:%s"`, f.Name()),
				fmt.Sprintf(`"json:%s"`, f.Name()),
			},
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

func selectType(types []structbuilder.Type) (string, string) {
	t := chooseType(types)
	typeName, importPath := typeNameAndImportPath(t)
	switch tt := t.(type) {
	case *structbuilder.ArrayBuilder:
		typeName, importPath = selectType(tt.Types())
		typeName = "[]" + typeName
	}

	return typeName, importPath
}

func typeNameAndImportPath(t structbuilder.Type) (string, string) {
	parts := strings.SplitN(t.Name(), " ", 2)
	typeName := parts[0]
	importPath := ""
	if len(parts) == 2 {
		importPath = parts[1]
	}

	return typeName, importPath
}
