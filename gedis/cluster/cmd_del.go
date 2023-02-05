package cluster

import (
	databaseface "github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
)

func del(cluster *ClusterDatabase, conn resp.Connection, args databaseface.CmdLine) resp.Reply {
	broadcastReply := cluster.broadcast(conn, args)
	var deleted int64
	for peer, replyResp := range broadcastReply {
		if reply.IsErrorReply(replyResp) {
			return replyResp
		}
		switch rr := replyResp.(type) {
		case *reply.IntReply:
			deleted += rr.Code
			continue
		}
		return reply.MakeErrReply("error: " + peer)
	}
	return reply.MakeIntReply(deleted)
}
