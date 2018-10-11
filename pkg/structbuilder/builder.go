package structbuilder

import "fmt"

// Type is the interface for building types.
type Type interface {
	Name() string
	Count() uint
	_type()
}

func (*StructBuilder) _type()    {}
func (*ArrayBuilder) _type()     {}
func (*PrimitiveBuilder) _type() {}

// NewStructBuilder creates a new struct builder.
func NewStructBuilder(name string) *StructBuilder {
	return &StructBuilder{
		name:  name,
		count: 1,
	}
}

// StructBuilder helps building structs.
type StructBuilder struct {
	name   string
	count  uint
	fields []*FieldBuilder
}

// Name implements the Builder interface.
func (sb *StructBuilder) Name() string {
	return sb.name
}

// Count implements the Builder interface.
func (sb *StructBuilder) Count() uint {
	return sb.count
}

// Fields returns the fields in the struct.
func (sb *StructBuilder) Fields() []*FieldBuilder {
	return sb.fields
}

// Include includes the specified field/type in the struct.
func (sb *StructBuilder) Include(name string, t Type) {
	for _, f := range sb.fields {
		if f.name == name {
			f.count++
			f.include(t)
			return
		}
	}

	f := FieldBuilder{
		name:  name,
		count: 1,
	}

	f.include(t)

	sb.fields = append(sb.fields, &f)
}

// Merge merges two StructBuilders.
func (sb *StructBuilder) Merge(other *StructBuilder) {
	sb.count += other.count
	for _, of := range other.fields {
		found := false
		for _, sbf := range sb.fields {
			if sbf.name == of.name {
				sbf.merge(of)
				found = true
				break
			}
		}

		if !found {
			sb.fields = append(sb.fields, of)
		}
	}
}

// FieldBuilder helps build structs.
type FieldBuilder struct {
	name  string
	count uint
	types []Type
}

// Name returns the name of the field.
func (fb *FieldBuilder) Name() string {
	return fb.name
}

// Count returns the number of times the field was included.
func (fb *FieldBuilder) Count() uint {
	return fb.count
}

// Types returns the fields types.
func (fb *FieldBuilder) Types() []Type {
	return fb.types
}

func (fb *FieldBuilder) include(t Type) {
	for _, fbt := range fb.types {
		if fbt.Name() == t.Name() {
			merge(fbt, t)
			return
		}
	}

	fb.types = append(fb.types, t)
}

func (fb *FieldBuilder) merge(other *FieldBuilder) {
	for _, ot := range other.types {
		found := false
		for _, fbt := range fb.types {
			if fbt.Name() == ot.Name() {
				merge(fbt, ot)
				found = true
				break
			}
		}
		if !found {
			fb.types = append(fb.types, ot)
		}
	}
}

// NewPrimitiveBuilder creates a new primitive type.
func NewPrimitiveBuilder(name string) *PrimitiveBuilder {
	return &PrimitiveBuilder{
		name:  name,
		count: 1,
	}
}

// PrimitiveBuilder is a type for primitives.
type PrimitiveBuilder struct {
	name  string
	count uint
}

// Name implements the Type interface.
func (pb *PrimitiveBuilder) Name() string {
	return pb.name
}

// Count implements the Type interface.
func (pb *PrimitiveBuilder) Count() uint {
	return pb.count
}

// Merge merges two PrimitiveBuilders.
func (pb *PrimitiveBuilder) Merge(other *PrimitiveBuilder) {
	pb.count += other.count
}

// NewArrayBuilder creates a new array builder.
func NewArrayBuilder(name string) *ArrayBuilder {
	return &ArrayBuilder{
		name:  name,
		count: 1,
	}
}

// ArrayBuilder is a type for arrays.
type ArrayBuilder struct {
	name  string
	count uint
	types []Type
}

// Name implements the Type interface.
func (ab *ArrayBuilder) Name() string {
	return ab.name
}

// Count implements the Type interface.
func (ab *ArrayBuilder) Count() uint {
	return ab.count
}

// Include increments the count of the specified name and type by 1.
func (ab *ArrayBuilder) Include(t Type) {
	for _, abt := range ab.types {
		if abt.Name() == t.Name() {
			merge(abt, t)
			return
		}
	}

	ab.types = append(ab.types, t)
}

// Merge merges two ArrayBuilders.
func (ab *ArrayBuilder) Merge(other *ArrayBuilder) {
	for _, ot := range other.types {
		found := false
		for _, abt := range ab.types {
			if abt.Name() == ot.Name() {
				merge(abt, ot)
				found = true
			}
		}
		if !found {
			ab.types = append(ab.types, ot)
		}
	}
}

// Types returns the fields types.
func (ab *ArrayBuilder) Types() []Type {
	return ab.types
}

// Merge merges the two types. The two types must be the same.
func merge(a, b Type) {
	aPrim, aOK := a.(*PrimitiveBuilder)
	bPrim, bOK := b.(*PrimitiveBuilder)
	if aOK && bOK {
		aPrim.Merge(bPrim)
		return
	} else if aOK != bOK {
		panic(fmt.Sprintf("cannot merge %T and %T", a, b))
	}

	aArray, aOK := a.(*ArrayBuilder)
	bArray, bOK := b.(*ArrayBuilder)
	if aOK && bOK {
		aArray.Merge(bArray)
		return
	} else if aOK != bOK {
		panic(fmt.Sprintf("cannot merge %T and %T", a, b))
	}

	aStruct, aOK := a.(*StructBuilder)
	bStruct, bOK := b.(*StructBuilder)
	if aOK && bOK {
		aStruct.Merge(bStruct)
		return
	} else if aOK != bOK {
		panic(fmt.Sprintf("cannot merge %T and %T", a, b))
	}

	panic(fmt.Sprintf("%T and/or %T is unsupported", a, b))
}
