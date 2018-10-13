package bsonutil

import (
	"github.com/mongodb/mongo-go-driver/bson"
)

// NewTypeBuilder makes a TypeBuilder.
func NewTypeBuilder() *TypeBuilder {
	return &TypeBuilder{}
}

// TypeBuilder is used to build up a type.
type TypeBuilder struct {
	Fields     map[string]*TypeBuilder
	Array      *TypeBuilder
	Primitives map[string]uint

	CanBeNull bool
	Count     uint
}

func (tb *TypeBuilder) IncludeDocument(doc *bson.Document) {
	tb.Count++

	iter := doc.Iterator()
	for iter.Next() {
		e := iter.Element()
		tb.includeField(e.Key(), e.Value())
	}
}

func (tb *TypeBuilder) includeField(name string, v *bson.Value) {
	if tb.Fields == nil {
		tb.Fields = make(map[string]*TypeBuilder)
	}

	var ftb *TypeBuilder
	var ok bool
	if ftb, ok = tb.Fields[name]; !ok {
		ftb = NewTypeBuilder()
		tb.Fields[name] = ftb
	}

	ftb.includeValue(v)
}

func (tb *TypeBuilder) includeValue(v *bson.Value) {
	tb.Count++
	switch v.Type() {
	case bson.TypeArray:
		if tb.Array == nil {
			tb.Array = NewTypeBuilder()
		}
		tb.Array.includeArray(v.MutableArray())
	case bson.TypeEmbeddedDocument:
		tb.IncludeDocument(v.MutableDocument())
	default:
		tb.includePrimitive(v)
	}
}

func (tb *TypeBuilder) includeArray(arr *bson.Array) {
	tb.Count++
	iter, _ := arr.Iterator()
	for iter.Next() {
		tb.includeValue(iter.Value())
	}
}

func (tb *TypeBuilder) includePrimitive(v *bson.Value) {
	name := mapPrimitiveTypeName(v.Type())
	if tb.Primitives == nil {
		tb.Primitives = make(map[string]uint)
	}

	tb.Primitives[name]++
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
	case bson.TypeNull:
		return "null"
	default:
		return "*bson.Value github.com/mongodb/mongo-go-driver/bson"
	}
}

// // Merge merges two StructBuilders.
// func (tb *TypeBuilder) Merge(other *TypeBuilder) {
// 	tb.Count += other.Count
// 	tb.CanBeNull = tb.CanBeNull || other.CanBeNull

// 	if tb.Array != nil && other.Array != nil {
// 		tb.Array.Merge(other.Array)
// 	} else if other.Array != nil {
// 		tb.Array = other.Array
// 	}

// 	if tb.Primitives != nil && other.Primitives != nil {
// 		for key := range tb.Primitives {
// 			if op, ok := other.Primitives[key]; ok {
// 				tb.Primitives[key] += op
// 			}
// 		}

// 		for key := range other.Primitives {
// 			if _, ok := tb.Primitives[key]; !ok {
// 				tb.Primitives[key] = other.Primitives[key]
// 			}
// 		}
// 	} else if other.Primitives != nil {
// 		tb.Primitives = other.Primitives
// 	}

// 	if tb.Fields != nil && other.Fields != nil {
// 		for key := range tb.Fields {
// 			if of, ok := other.Fields[key]; ok {
// 				tb.Fields[key].Merge(of)
// 			}
// 		}

// 		for key := range other.Fields {
// 			if _, ok := tb.Fields[key]; !ok {
// 				tb.Fields[key] = other.Fields[key]
// 			}
// 		}
// 	}
// }
