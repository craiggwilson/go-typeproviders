package bsonutil

import (
	"fmt"
	"strings"

	"github.com/craiggwilson/go-typeproviders/pkg/naming"
	"github.com/craiggwilson/go-typeproviders/pkg/structbuilder"
	"github.com/mongodb/mongo-go-driver/bson"
)

func BuildStruct(name string, tb *TypeBuilder) *structbuilder.Struct {
	s := structbuilder.Struct{
		Name: naming.Struct(name),
	}

	for fieldName, ftb := range tb.Fields {
		exportedFieldName := naming.ExportedField(fieldName)
		path := s.Name + exportedFieldName
		fieldType := selectType(path, tb.Count, ftb)
		if fieldType.ArrayCount > 0 {
			exportedFieldName = naming.Pluralize(exportedFieldName)
		}
		s.Fields = append(s.Fields, &structbuilder.Field{
			Name: exportedFieldName,
			Tags: []string{
				fmt.Sprintf(`"bson:%s"`, fieldName),
				fmt.Sprintf(`"json:%s"`, fieldName),
			},
			Type: &fieldType,
		})
	}

	structbuilder.SortFieldsByName(s.Fields)

	return &s
}

func selectType(path string, seenCount uint, tb *TypeBuilder) structbuilder.FieldType {
	canBeNull := false
	totalTypeCount := uint(0)

	var fieldTypes []structbuilder.FieldType

	if tb.Fields != nil {
		// we found a document
		totalTypeCount += tb.Count

		rs := BuildStruct(path, tb)
		fieldTypes = append(fieldTypes, structbuilder.FieldType{
			Name:           rs.Name,
			EmbeddedStruct: rs,
		})
	}
	if tb.Array != nil {
		// we found an array
		totalTypeCount += tb.Array.Count
		elementFieldType := selectType(path, tb.Count, tb.Array)
		elementFieldType.ArrayCount++
		fieldTypes = append(fieldTypes, elementFieldType)
	}
	for key, primitiveCount := range tb.Primitives {
		if key == "null" {
			canBeNull = true
			continue
		}
		totalTypeCount += primitiveCount

		typeName, importPath := typeNameAndImportPath(key)
		fieldTypes = append(fieldTypes, structbuilder.FieldType{
			Name:       typeName,
			ImportPath: importPath,
		})
	}

	if seenCount > totalTypeCount {
		canBeNull = true
	}

	switch len(fieldTypes) {
	case 0:
		typeName, importPath := typeNameAndImportPath(mapPrimitiveTypeName(bson.TypeUndefined))
		return structbuilder.FieldType{
			Name:       typeName,
			ImportPath: importPath,
			CanBeNull:  true,
		}
	case 1:
		fieldTypes[0].CanBeNull = canBeNull
		return fieldTypes[0]
	default:
		return structbuilder.FieldType{
			Name: "blah",
		}
	}
}

func typeNameAndImportPath(name string) (string, string) {
	parts := strings.SplitN(name, " ", 2)
	typeName := parts[0]
	importPath := ""
	if len(parts) == 2 {
		importPath = parts[1]
	}

	return typeName, importPath
}
