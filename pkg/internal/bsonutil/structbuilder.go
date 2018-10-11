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

// BuildStruct builds multiple []*structbuilder.Structs from a structbuilder.StructBuilder.
func BuildStructs(name string, sb *structbuilder.StructBuilder, embedStructs bool) ([]*structbuilder.Struct, error) {
	var results []*structbuilder.Struct

	type buildItem struct {
		b    *structbuilder.StructBuilder
		path string
	}

	structsToBuild := []buildItem{{sb, name}}

	for len(structsToBuild) > 0 {
		current := structsToBuild[0]
		structsToBuild = structsToBuild[1:]

		sb = current.b
		path := current.path

		s := structbuilder.Struct{
			Name: naming.Struct(path),
		}

		for _, f := range sb.Fields() {
			nestedStructPath := path + "_" + f.Name()
			t, typeName, importPath := selectType(nestedStructPath, f.Types())
			fieldName := naming.ExportedField(f.Name())
			if strings.HasPrefix(typeName, "[]") {
				fieldName = naming.Pluralize(fieldName)
			}

			if st, ok := t.(*structbuilder.StructBuilder); ok {
				structsToBuild = append(structsToBuild, buildItem{
					b:    st,
					path: nestedStructPath,
				})
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

		results = append(results, &s)
	}

	return results, nil
}

// selectType takes the current path into the type selection and a list of types.
// It returns the dominate type as well as the typeName and importPath to use.
func selectType(path string, types []structbuilder.Type) (structbuilder.Type, string, string) {
	t := chooseType(types)
	typeName, importPath := typeNameAndImportPath(t)
	switch tt := t.(type) {
	case *structbuilder.ArrayBuilder:
		t, typeName, importPath = selectType(path, tt.Types())
		typeName = "[]" + typeName
	case *structbuilder.StructBuilder:
		typeName = naming.Struct(path)
		importPath = ""
	}

	return t, typeName, importPath
}

func chooseType(types []structbuilder.Type) structbuilder.Type {
	return types[0]
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
