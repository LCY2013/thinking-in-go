package handler

import (
	"context"
	"github.com/LCY2013/thinking-in-go/gedis/database"
	databaseface "github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"github.com/LCY2013/thinking-in-go/gedis/lib/sync/wait"
	"github.com/LCY2013/thinking-in-go/gedis/resp/connection"
	"github.com/LCY2013/thinking-in-go/gedis/resp/parser"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
 * A tcp.RespHandler implements redis protocol
 */

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

// RespHandler implements tcp.Handler and serves as a redis handler
type RespHandler struct {
	activeConn sync.Map              // *client -> placeholder
	db         databaseface.Database // database interface
	closing    atomic.Bool           // refusing new client and new request
	wait       wait.Wait
}

// MakeHandler creates a RespHandler instance
func MakeHandler() *RespHandler {
	return &RespHandler{
		//db: database.NewEchoDatabase(),
		db: database.NewDatabase(),
	}
}

// closeClient closes the connection for redis-cli
func (r *RespHandler) closeClient(conn *connection.Connection) {
	_ = conn.Close()
	r.db.AfterClientClose(conn)
	r.activeConn.Delete(conn)
}

// Handle receives and executes redis commands
func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Load() {
		// closing handler refuse new connection
		_ = conn.Close()
		return
	}

	clientConn := connection.NewConn(conn)
	r.activeConn.Store(clientConn, struct{}{})

	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// connection closed
				r.closeClient(clientConn)
				logger.Info("connection closed: " + clientConn.RemoteAddr().String())
				return
			}
			// protocol err
			errReply := reply.MakeErrReply(payload.Err.Error())
			err := clientConn.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(clientConn)
				logger.Info("connection closed: " + clientConn.RemoteAddr().String())
				return
			}
			continue
		}
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}

		switch mbr := payload.Data.(type) {
		case *reply.BulkReply:
			result := r.db.Exec(clientConn, [][]byte{
				mbr.Arg,
			})
			if result != nil {
				_ = clientConn.Write(result.ToBytes())
			} else {
				_ = clientConn.Write(unknownErrReplyBytes)
			}
			continue
		case *reply.StatusReply:
			result := r.db.Exec(clientConn, [][]byte{
				[]byte(mbr.Status),
			})
			if result != nil {
				_ = clientConn.Write(result.ToBytes())
			} else {
				_ = clientConn.Write(unknownErrReplyBytes)
			}
			continue
		case *reply.MultiBulkReply:
			result := r.db.Exec(clientConn, mbr.Args)
			if result != nil {
				_ = clientConn.Write(result.ToBytes())
			} else {
				_ = clientConn.Write(unknownErrReplyBytes)
			}
			continue
		}

		_ = clientConn.Write(unknownErrReplyBytes)
		/*mbr, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := r.db.Exec(clientConn, mbr.Args)
		if result != nil {
			_ = clientConn.Write(result.ToBytes())
		} else {
			_ = clientConn.Write(unknownErrReplyBytes)
		}*/
	}

	// close the connection
	_ = conn.Close()
}

// Close stops handler
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down...")
	r.closing.Store(true)
	var (
		err error
	)

	// close all connection for this handler (redis-cli)
	r.wait.AsyncDo(func() {
		r.activeConn.Range(func(key, value any) bool {
			if conn, ok := key.(*connection.Connection); ok {
				err = conn.Close()
			}
			return true
		})
	})

	// wait for closed timeout or close finished
	r.wait.WaitWithTimeout(15 * time.Second)

	return err
}
