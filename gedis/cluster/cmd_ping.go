package cluster

import (
	databaseface "github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
)

func ping(cluster *ClusterDatabase, conn resp.Connection, args databaseface.CmdLine) resp.Reply {
	return cluster.db.Exec(conn, args)
}
