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
	Fields     []*FieldBuilder
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
	for _, fb := range tb.Fields {
		if fb.Name == name {
			fb.includeValue(v)
			return
		}
	}

	fb := NewFieldBuilder(name)
	fb.includeValue(v)
	tb.Fields = append(tb.Fields, fb)
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

// NewFieldBuilder makes a FieldBuilder.
func NewFieldBuilder(name string) *FieldBuilder {
	return &FieldBuilder{
		Name:        name,
		TypeBuilder: NewTypeBuilder(),
	}
}

// FieldBuilder builds up fields.
type FieldBuilder struct {
	*TypeBuilder
	Name string
}
