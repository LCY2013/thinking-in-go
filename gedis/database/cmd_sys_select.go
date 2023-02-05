package database

import (
	"github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	"strconv"
)

func execSelect(c resp.Connection, mdb *StandaloneDatabase, args database.CmdLine) resp.Reply {
	if len(args) != 2 {
		return reply.MakeArgNumErrReply("select")
	}

	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(mdb.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}

func init() {
	RegisterSysCommand("select", Auth, 1)
}
