package mongodb

// Config holds information required for configuration mongodb.
type Config struct {
	URI            string
	DatabaseName   string
	CollectionName string

	Package string
}
