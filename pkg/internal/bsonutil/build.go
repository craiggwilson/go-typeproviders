package bsonutil

import (
	"fmt"
	"strings"

	"github.com/craiggwilson/go-typeproviders/pkg/naming"
	"github.com/craiggwilson/go-typeproviders/pkg/structbuilder"
	"github.com/mongodb/mongo-go-driver/bson"
)

func BuildStructs(name string, tb *TypeBuilder, embedStructs bool) ([]*structbuilder.Struct, error) {
	var results []*structbuilder.Struct

	type buildItem struct {
		tb   *TypeBuilder
		path string
	}

	structsToBuild := []buildItem{{tb, name}}

	for len(structsToBuild) > 0 {
		current := structsToBuild[0]
		structsToBuild = structsToBuild[1:]

		tb = current.tb
		path := current.path

		s := structbuilder.Struct{
			Name: naming.Struct(path),
		}

		for name, ftb := range tb.Fields {
			nestedStructPath := path + "_" + name
			fieldName := naming.ExportedField(name)
			_, fieldType := selectType(nestedStructPath, tb.Count, ftb)
			if strings.HasPrefix(fieldType.Name, "[]") {
				fieldName = naming.Pluralize(fieldName)
			}
			s.Fields = append(s.Fields, structbuilder.Field{
				Name: fieldName,
				Tags: []string{
					fmt.Sprintf(`"bson:%s"`, name),
					fmt.Sprintf(`"json:%s"`, name),
				},
				Type: fieldType,
			})
		}

		results = append(results, &s)
	}

	return results, nil
}

func selectType(path string, seenCount uint, tb *TypeBuilder) (*structbuilder.Struct, structbuilder.FieldType) {
	s, t := selectTypeExt(path, seenCount, tb)
	arrayPrefix := strings.Repeat("[]", int(t.arrayCount))
	t.FieldType.Name = arrayPrefix + t.FieldType.Name
	return s, t.FieldType
}

type fieldTypeExt struct {
	structbuilder.FieldType
	arrayCount uint
}

func selectTypeExt(path string, seenCount uint, tb *TypeBuilder) (*structbuilder.Struct, fieldTypeExt) {
	canBeNull := false
	totalTypeCount := uint(0)

	var resultStruct *structbuilder.Struct
	var fieldTypes []fieldTypeExt

	if tb.Fields != nil {
		// we found a document
		totalTypeCount += tb.Count
	}
	if tb.Array != nil {
		// we found an array
		totalTypeCount += tb.Array.Count
		nestedStruct, elementFieldType := selectTypeExt(path, tb.Array.Count, tb.Array)
		fieldTypes = append(fieldTypes, fieldTypeExt{
			FieldType:  elementFieldType.FieldType,
			arrayCount: 1 + elementFieldType.arrayCount,
		})

		resultStruct = nestedStruct
	}
	for key, primitiveCount := range tb.Primitives {
		if key == "null" {
			canBeNull = true
			continue
		}
		totalTypeCount += primitiveCount

		typeName, importPath := typeNameAndImportPath(key)
		fieldTypes = append(fieldTypes, fieldTypeExt{
			FieldType: structbuilder.FieldType{
				Name:       typeName,
				ImportPath: importPath,
			},
		})
	}

	fmt.Println("TOTALS", path, seenCount, totalTypeCount)
	if seenCount > totalTypeCount {
		canBeNull = true
	}

	switch len(fieldTypes) {
	case 0:
		typeName, importPath := typeNameAndImportPath(mapPrimitiveTypeName(bson.TypeUndefined))
		return nil, fieldTypeExt{
			FieldType: structbuilder.FieldType{
				Name:       typeName,
				ImportPath: importPath,
			},
		}
	case 1:
		if canBeNull {
			fieldTypes[0].Name = "*" + fieldTypes[0].Name
		}
		return resultStruct, fieldTypes[0]
	default:
		return nil, fieldTypeExt{
			FieldType: structbuilder.FieldType{
				Name: "string",
			},
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
