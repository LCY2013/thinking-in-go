package resp

// Connection represents a connection with redis client
type Connection interface {
	Write(buf []byte) error
	// GetDBIndex used for multi database
	GetDBIndex() int
	SelectDB(int)
	// SetPassword password
	SetPassword(password string)
	GetPassword() string
}
