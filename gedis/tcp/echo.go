package tcp

import (
	"bufio"
	"context"
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"github.com/LCY2013/thinking-in-go/gedis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
A echo server to test whether the server is functioning normally
*/

var (
	BYE = map[string]struct{}{
		"quit\r\n": struct{}{},
		"exit\r\n": struct{}{},
	}
	BYE_STRING = []byte("bye...\r\n")
)

// EchoHandler echos received line to client, using for test
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Bool
}

// MakeEchoHandler creates EchoHandler
func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

// EchoClient is client for EchoHandler, using for test
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close conn connection
func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(15 * time.Second)
	return c.Conn.Close()
}

// Handle echos received line to client
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Load() {
		// closing handler refuse new connection
		_ = conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}

	h.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		// may occurs: client EOF, client timeout, server early close
		readString, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.SyncDo(func() {
			b := []byte(readString)
			content := string(b)
			if _, ok := BYE[content]; ok {
				_, _ = client.Conn.Write(BYE_STRING)
				h.activeConn.Delete(client)
				_ = client.Conn.Close()
				return
			}
			_, _ = client.Conn.Write(b)
		})
	}
}

// Close stops echo handler
func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	h.closing.Store(true)

	// close all connections
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
