package structbuilder

// StructBuilder helps build structs.
type StructBuilder struct {
	Name   string
	Fields []StructFieldBuilder
}

// StructFieldBuilder helps build struct fields.
type StructFieldBuilder struct {
	Name string
	Type string
	Tags string
}

// QuotedTags gets the tags quoted with a backtick.
func (sfb *StructFieldBuilder) QuotedTags() string {
	if len(sfb.Tags) > 0 {
		return "`" + sfb.Tags + "`"
	}

	return ""
}
