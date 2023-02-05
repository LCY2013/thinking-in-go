package cluster

import (
	databaseface "github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
)

func FlushDB(cluster *ClusterDatabase, conn resp.Connection, args databaseface.CmdLine) resp.Reply {
	replies := cluster.broadcast(conn, args)

	for _, v := range replies {
		if reply.IsErrorReply(v) {
			errReply := v.(reply.ErrorReply)
			return reply.MakeErrReply("error occurs: " + errReply.Error())
		}
	}

	return &reply.OkReply{}
}
