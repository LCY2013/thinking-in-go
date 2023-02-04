// Package database is a memory database with redis compatible interface
package database

import (
	"github.com/LCY2013/thinking-in-go/gedis/datastruct/dict"
	"github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	"strings"
)

// DB stores data and execute user's commands
type DB struct {
	index int
	// key -> DataEntity
	data dict.Dict
}

// ExecFunc is interface for command executor
// args don't include cmd line
type ExecFunc func(db *DB, args CmdLine) resp.Reply

// CmdLine is alias for [][]byte, represents a command line
type CmdLine [][]byte

// makeDB create DB instance
func makeDB() *DB {
	return &DB{
		data: dict.MakeConcurrent(128),
	}
}

// Exec executes command within one database
func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command '" + cmdName + "'")
	}

	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}

	fn := cmd.executor

	return fn(db, cmdLine[1:])
}

func validateArity(arity int, cmdArgs CmdLine) bool {
	if arity >= 0 {
		return arity == len(cmdArgs)
	}

	return len(cmdArgs) >= -arity
}

/* ---- data Access ----- */

// GetEntity returns DataEntity bind to given key
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {

	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

// PutEntity a DataEntity into DB
func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfExists edit an existing DataEntity
func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

// PutIfAbsent insert an DataEntity only if the key not exists
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

// Remove the given key from db
func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

// Removes the given keys from db
func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

// Flush clean database
func (db *DB) Flush() {
	db.data.Clear()

}
