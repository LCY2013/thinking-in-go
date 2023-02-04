package resp

// Connection represents a connection with redis client
type Connection interface {
	Write(buf []byte) error
	// GetDBIndex used for multi database
	GetDBIndex() int
	SelectDB(int)
}

// DataEntity stores data bound to a key, including a string, list, hash, set and so on
type DataEntity[T any] struct {
	Data T
}
