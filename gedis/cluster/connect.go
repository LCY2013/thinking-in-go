package cluster

import (
	"context"
	"errors"
	"github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/lib/utils"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	"github.com/LCY2013/thinking-in-go/gedis/tcp/client"
	"strconv"
)

// getPeerClient get client for peer
func (c *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	factory, ok := c.peerConn[peer]
	if !ok {
		return nil, errors.New("connection factory not found")
	}
	raw, err := factory.BorrowObject(context.TODO())
	if err != nil {
		return nil, err
	}
	client, ok := raw.(*client.Client)
	if !ok {
		return nil, errors.New("connection factory make wrong type")
	}
	return client, nil
}

// returnPeerClient
func (c *ClusterDatabase) returnPeerClient(peer string, peerClient *client.Client) error {
	factory, ok := c.peerConn[peer]
	if !ok {
		return errors.New("connection factory not found")
	}
	return factory.ReturnObject(context.TODO(), peerClient)
}

// relay relays command to peer
// select db by c.GetDBIndex()
// cannot call Prepare, Commit, execRollback of self node
func (c *ClusterDatabase) relay(peer string, conn resp.Connection, args database.CmdLine) resp.Reply {
	if peer == c.self {
		// self do
		return c.db.Exec(conn, args)
	}

	// remote peer do
	peerClient, err := c.getPeerClient(peer)
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}

	defer func() {
		_ = c.returnPeerClient(peer, peerClient)
	}()

	// select db
	peerClient.Send(utils.ToCmdLine("select", strconv.Itoa(conn.GetDBIndex())))

	return peerClient.Send(args)
}

// broadcast broadcasts command to all node in cluster
func (c *ClusterDatabase) broadcast(conn resp.Connection, args database.CmdLine) map[string]resp.Reply {
	reply := make(map[string]resp.Reply, len(c.nodes))

	for _, peer := range c.nodes {
		reply[peer] = c.relay(peer, conn, args)
	}

	return reply
}
