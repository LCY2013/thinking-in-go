package connection

import (
	"github.com/LCY2013/thinking-in-go/gedis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

// Connection represents a connection with a redis-cli
type Connection struct {
	conn net.Conn
	// waiting until reply finished
	waitingReply wait.Wait
	// lock while handler sending response
	mu sync.Mutex
	// select db
	selectDB int
	// password
	password string
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

// RemoteAddr returns the remote network address
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// Close disconnect with the client
func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(15 * time.Second)
	_ = c.conn.Close()
	return nil
}

func (c *Connection) Write(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return nil
	}
	var (
		err error
	)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.waitingReply.SyncDo(func() {
		_, err = c.conn.Write(buf)
	})

	return err
}

func (c *Connection) GetDBIndex() int {
	return c.selectDB
}

func (c *Connection) SelectDB(db int) {
	c.selectDB = db
}

func (c *Connection) SetPassword(password string) {
	c.password = password
}

func (c *Connection) GetPassword() string {
	return c.password
}
