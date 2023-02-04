package database

import "github.com/LCY2013/thinking-in-go/gedis/interface/resp"

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// Database is the interface for redis style storage engine
type Database interface {
	Exec(client resp.Connection, args CmdLine) resp.Reply
	AfterClientClose(conn resp.Connection)
	Close()
}
