package database

import (
	"github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
)

// EchoDatabase echo database
type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e EchoDatabase) Exec(client resp.Connection, args database.CmdLine) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e EchoDatabase) AfterClientClose(conn resp.Connection) {
	logger.Info("EchoDatabase.AfterClientClose")
}

func (e EchoDatabase) Close() {
	logger.Info("EchoDatabase.Close")
}
