package database

import (
	"github.com/LCY2013/thinking-in-go/gedis/aof"
	"github.com/LCY2013/thinking-in-go/gedis/config"
	"github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"strings"
)

// StandaloneDatabase is a set of multiple database set
type StandaloneDatabase struct {
	dbSet []*DB
	// handle aof persistence
	aofHandler *aof.AofHandler
}

// NewStandaloneDatabase creates a redis database,
func NewStandaloneDatabase() *StandaloneDatabase {
	ndb := &StandaloneDatabase{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	ndb.dbSet = make([]*DB, config.Properties.Databases)

	for idx := range ndb.dbSet {
		singleDB := makeDB()
		singleDB.index = idx
		ndb.dbSet[idx] = singleDB
	}

	// init aof append
	if config.Properties.AppendOnly {
		handler, err := aof.NewAOFHandler(ndb)
		if err != nil {
			panic(err)
		}
		ndb.aofHandler = handler
		for _, db := range ndb.dbSet {
			// avoid closure
			singleDB := db
			singleDB.addAof = func(line database.CmdLine) {
				ndb.aofHandler.AddAof(singleDB.index, line)
			}
		}
	}

	return ndb
}

// Exec executes command
// parameter `cmdLine` contains command and its arguments, for example: "set key value"
func (mdb *StandaloneDatabase) Exec(c resp.Connection, cmdLine database.CmdLine) (result resp.Reply) {
	cmdName := strings.ToLower(string(cmdLine[0]))

	// system command
	sysCmd, ok := cmdSysTable[cmdName]
	if ok {
		return sysCmd.executor(c, cmdLine)
	}

	// normal command
	dbIndex := c.GetDBIndex()
	selectedDB := mdb.dbSet[dbIndex]
	return selectedDB.Exec(c, cmdLine)
}

// Close graceful shutdown database
func (mdb *StandaloneDatabase) Close() {

}

func (mdb *StandaloneDatabase) AfterClientClose(c resp.Connection) {
}
