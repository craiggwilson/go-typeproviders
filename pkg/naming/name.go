package naming

import (
	"github.com/craiggwilson/go-typeproviders/pkg/internal/inflect"
)

func init() {
	acronyms := []string{
		"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP",
		"HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA",
		"SMTP", "SSH", "TCP", "TLS", "TTL", "UDP", "UI", "UID", "UUID", "URI",
		"URL", "UTF8", "VM", "XML", "XSRF", "XSS",
	}

	for _, acronym := range acronyms {
		inflect.AddAcronym(acronym)
	}
}

// Struct returns a proper name for a struct
func Struct(name string) string {
	return inflect.Camelize(inflect.Singularize(name))
}

// ExportedField returns a proper name for a field.
func ExportedField(name string) string {
	return inflect.Camelize(name)
}

// Pluralize returns a plural form of the name.
func Pluralize(name string) string {
	return inflect.Pluralize(name)
}

// Singularize returns a singular form of the name.
func Singularize(name string) string {
	return inflect.Singularize(name)
}
