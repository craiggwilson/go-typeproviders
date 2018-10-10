package structbuilder

import "strings"

// Struct represents a struct.
type Struct struct {
	Name   string
	Fields []Field
	Tags   []string
}

// QuotedTags gets the tags quoted with a backtick.
func (s *Struct) QuotedTags() string {
	if len(s.Tags) > 0 {
		return "`" + strings.Join(s.Tags, " ") + "`"
	}

	return ""
}

// Field represents a field in a struct.
type Field struct {
	Name string
	Type string
	Tags []string
}

// QuotedTags gets the tags quoted with a backtick.
func (f *Field) QuotedTags() string {
	if len(f.Tags) > 0 {
		return "`" + strings.Join(f.Tags, " ") + "`"
	}

	return ""
}
