package buildsqlx

// Connection encloses DB struct
type Connection struct {
	driver string
}

// NewConnection returns pre-defined Connection structure
func NewConnection(driverName string) *Connection {
	return &Connection{driver: driverName}
}
