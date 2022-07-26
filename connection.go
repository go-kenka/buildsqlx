package buildsqlx

import "sync"

var (
	conn *Connection
	pool *sync.Pool
	once sync.Once
)

// Connection encloses DB struct
type Connection struct {
	driver string
}

// NewConnection returns pre-defined Connection structure
func NewConnection(driverName string) *Connection {
	once.Do(func() {
		conn = &Connection{driver: driverName}
		pool = &sync.Pool{
			New: func() any {
				return newDB(conn)
			},
		}
	})

	return conn
}

// DB get a sql builder
func (c *Connection) DB() *DB {
	return pool.Get().(*DB)
}

// DB is an entity that composite builder and Conn types
type DB struct {
	Builder *builder
	Conn    *Connection
}

// newDB constructs default DB structure
func newDB(c *Connection) *DB {
	b := newBuilder()
	return &DB{Builder: b, Conn: c}
}
