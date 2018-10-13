package structbuilder

import "strings"

// Struct represents a struct.
type Struct struct {
	Name   string
	Fields []*Field
	Tags   []string
}

// QuotedTags gets the tags quoted with a backtick.
func (s *Struct) QuotedTags() string {
	if len(s.Tags) > 0 {
		return "`" + strings.Join(s.Tags, " ") + "`"
	}

	return ""
}

// UnembedStructs unembeds all the structs of the children recursively.
func (s *Struct) UnembedStructs() []*Struct {
	results := []*Struct{s}
	for _, f := range s.Fields {
		if f.Type.EmbeddedStruct != nil {
			results = append(results, f.Type.EmbeddedStruct.UnembedStructs()...)
			f.Type.EmbeddedStruct = nil
		}
	}

	return results
}

// Field represents a field in a struct.
type Field struct {
	Name string
	Type *FieldType
	Tags []string
}

// FieldType represents the type of the field including its import path.
type FieldType struct {
	ImportPath string
	Name       string
	ArrayCount int
	CanBeNull  bool

	EmbeddedStruct *Struct
}
