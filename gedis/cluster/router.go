package cluster

import (
	"github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
)

func makeRouter() map[string]CmdFunc {
	return map[string]CmdFunc{
		"ping": ping,
		"del":  del,

		"exists":   defaultFunc,
		"type":     defaultFunc,
		"rename":   Rename,
		"renamenx": Rename,

		"set":    defaultFunc,
		"setnx":  defaultFunc,
		"get":    defaultFunc,
		"getset": defaultFunc,

		"flushdb": FlushDB,
	}
}

// defaultFunc relay command to responsible peer, and return its reply to client
func defaultFunc(cluster *ClusterDatabase, conn resp.Connection, args database.CmdLine) resp.Reply {
	key := string(args[1])
	peer := cluster.peerPicker.PickNode(key)
	return cluster.relay(peer, conn, args)
}
