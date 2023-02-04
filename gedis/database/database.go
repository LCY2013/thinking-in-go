package database

import (
	"github.com/LCY2013/thinking-in-go/gedis/config"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	"strconv"
	"strings"
)

// Database is a set of multiple database set
type Database struct {
	dbSet []*DB
}

// NewDatabase creates a redis database,
func NewDatabase() *Database {
	ndb := &Database{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	ndb.dbSet = make([]*DB, config.Properties.Databases)

	for idx := range ndb.dbSet {
		singleDB := makeDB()
		singleDB.index = idx
		ndb.dbSet[idx] = singleDB
	}

	return ndb
}

// Exec executes command
// parameter `cmdLine` contains command and its arguments, for example: "set key value"
func (mdb *Database) Exec(c resp.Connection, cmdLine [][]byte) (result resp.Reply) {
	cmdName := strings.ToLower(string(cmdLine[0]))
	if cmdName == "select" {
		if len(cmdLine) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(c, mdb, cmdLine[1:])
	}
	// normal commands
	dbIndex := c.GetDBIndex()
	selectedDB := mdb.dbSet[dbIndex]
	return selectedDB.Exec(c, cmdLine)
}

// Close graceful shutdown database
func (mdb *Database) Close() {

}

func (mdb *Database) AfterClientClose(c resp.Connection) {
}

func execSelect(c resp.Connection, mdb *Database, args [][]byte) resp.Reply {
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
