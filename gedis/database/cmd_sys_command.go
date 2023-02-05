package database

import "strings"

var cmdSysTable = make(map[string]*commandSys)

type commandSys struct {
	executor ExecSysFunc
	arity    int // allow number of args, arity < 0 means len(args) >= -arity
}

// RegisterSysCommand registers a new command
// arity means allowed number of cmdArgs, arity < 0 means len(args) >= -arity.
// for example: the arity of `get` is 2, `mget` is -2
func RegisterSysCommand(name string, executor ExecSysFunc, arity int) {
	name = strings.ToLower(name)
	cmdSysTable[name] = &commandSys{
		executor: executor,
		arity:    arity,
	}
}
