package database

import (
	"github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/lib/utils"
	"github.com/LCY2013/thinking-in-go/gedis/lib/wildcard"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
)

// execDel removes a key from db
func execDel(db *DB, args database.CmdLine) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}

	deleted := db.Removes(keys...)

	// aof handler
	if deleted > 0 {
		db.addAof(utils.ToMergeCmdLine("del", args...))
	}

	return reply.MakeIntReply(int64(deleted))
}

// execExists checks if a is existed in db
func execExists(db *DB, args database.CmdLine) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// execFlushDB removes all data in current db
func execFlushDB(db *DB, args database.CmdLine) resp.Reply {
	db.Flush()

	db.addAof(utils.ToMergeCmdLine("flushdb", args...))

	return &reply.OkReply{}
}

// execType returns the type of entity, including: string, list, hash, set and zset
func execType(db *DB, args database.CmdLine) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
		//case *list.LinkedList:
		//    return reply.MakeStatusReply("list")
		//case dict.Dict:
		//    return reply.MakeStatusReply("hash")
		//case *set.Set:
		//    return reply.MakeStatusReply("set")
		//case *sortedset.SortedSet:
		//    return reply.MakeStatusReply("zset")
	}
	return &reply.UnknownErrReply{}
}

// execRename a key
func execRename(db *DB, args database.CmdLine) resp.Reply {
	if len(args) != 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	src := string(args[0])
	dest := string(args[1])

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)

	db.addAof(utils.ToMergeCmdLine("rename", args...))

	return &reply.OkReply{}
}

// execRenameNx a key, only if the new key does not exist
func execRenameNx(db *DB, args database.CmdLine) resp.Reply {
	src := string(args[0])
	dest := string(args[1])

	_, ok := db.GetEntity(dest)
	if ok {
		return reply.MakeIntReply(0)
	}

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	db.Removes(src, dest) // clean src and dest with their ttl
	db.PutEntity(dest, entity)

	db.addAof(utils.ToMergeCmdLine("renamenx", args...))

	return reply.MakeIntReply(1)
}

// execKeys returns all keys matching the given pattern
func execKeys(db *DB, args database.CmdLine) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make(database.CmdLine, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}

func init() {
	RegisterCommand("Del", execDel, -2)
	RegisterCommand("Exists", execExists, -2)
	RegisterCommand("Keys", execKeys, 2)
	RegisterCommand("FlushDB", execFlushDB, -1)
	RegisterCommand("Type", execType, 2)
	RegisterCommand("Rename", execRename, 3)
	RegisterCommand("RenameNx", execRenameNx, 3)
}
