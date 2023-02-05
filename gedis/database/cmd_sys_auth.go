package database

import (
	"github.com/LCY2013/thinking-in-go/gedis/config"
	"github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
)

// Auth validate client's password
func Auth(c resp.Connection, args database.CmdLine) resp.Reply {
	if len(args) != 1 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'auth' command")
	}
	if config.Properties.RequirePass == "" {
		return reply.MakeErrReply("ERR Client sent AUTH, but no password is set")
	}
	passwd := string(args[0])
	c.SetPassword(passwd)

	if !isAuthenticated(c) {
		return reply.MakeErrReply("NOAUTH Authentication required")
	}

	if config.Properties.RequirePass != passwd {
		return reply.MakeErrReply("ERR invalid password")
	}
	return &reply.OkReply{}
}

func isAuthenticated(c resp.Connection) bool {
	if config.Properties.RequirePass == "" {
		return true
	}
	return c.GetPassword() == config.Properties.RequirePass
}

func init() {
	RegisterSysCommand("auth", Auth, 1)
}
