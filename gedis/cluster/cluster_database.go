// Package cluster provides a server side cluster which is transparent to client. You can connect to any node in the cluster to access all data in the cluster
package cluster

import (
	"context"
	"github.com/LCY2013/thinking-in-go/gedis/config"
	"github.com/LCY2013/thinking-in-go/gedis/database"
	databaseface "github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/lib/consistenthash"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	pool "github.com/jolestar/go-commons-pool/v2"
	"strings"
)

// CmdFunc represents the handler of a redis command
type CmdFunc func(cluster *ClusterDatabase, conn resp.Connection, args databaseface.CmdLine) resp.Reply

var (
	router = makeRouter()
)

// ClusterDatabase represents a node of godis cluster
// it holds part of data and coordinates other nodes to finish transactions
type ClusterDatabase struct {
	self string

	nodes      []string
	peerPicker *consistenthash.NodeMap
	peerConn   map[string]*pool.ObjectPool
	db         databaseface.Database
}

// NewClusterDatabase creates and starts a node of cluster
func NewClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self: config.Properties.Self,

		db:         database.NewStandaloneDatabase(),
		peerPicker: consistenthash.NewNodeMap(nil),
		peerConn:   make(map[string]*pool.ObjectPool),
	}

	// init peer connections
	/*for _, peer := range config.Properties.Peers {
		cluster.peerConn[peer] = pool.NewObjectPoolWithDefaultConfig(context.TODO(), &connectionFactory{
			peer: peer,
		})
	}*/

	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		if _, ok := cluster.peerConn[peer]; ok {
			continue
		}

		cluster.peerConn[peer] = pool.NewObjectPoolWithDefaultConfig(context.TODO(), &connectionFactory{
			peer: peer,
		})
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, config.Properties.Self)
	cluster.nodes = nodes

	// init consistent hash
	cluster.peerPicker.AddNode(nodes...)

	return cluster
}

func (c *ClusterDatabase) Exec(client resp.Connection, args databaseface.CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command '" + cmdName + "', or not supported in cluster mode")
	}
	return cmdFunc(c, client, args)
}

// AfterClientClose does some clean after client close connection
func (c *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	c.db.AfterClientClose(conn)
}

// Close stops current node of cluster
func (c *ClusterDatabase) Close() {
	c.db.Close()
}
