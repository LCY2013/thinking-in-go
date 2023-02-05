package main

import (
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	"net"
	"os"
	"testing"
)

var (
	conn net.Conn
	err  error
)

func init() {
	conn, err = net.Dial("tcp", "localhost:6379")
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}

func TestPing(t *testing.T) {
	replies := []resp.Reply{
		reply.MakeStatusReply("ping"),
	}

	for _, replie := range replies {
		_, err = conn.Write(replie.ToBytes())
		if err != nil {
			logger.Error(err)
		}

		msg := make([]byte, 1024)
		n, _ := conn.Read(msg)
		logger.Info(string(msg[:n]))
	}
}

func TestSet(t *testing.T) {
	replies := []resp.Reply{
		reply.MakeMultiBulkReply([][]byte{
			[]byte("set"),
			[]byte("hello"),
			[]byte("fufeng"),
		}),
		reply.MakeMultiBulkReply([][]byte{
			[]byte("set"),
			[]byte("key"),
			[]byte("value"),
		}),
	}

	for _, replie := range replies {
		_, err = conn.Write(replie.ToBytes())
		if err != nil {
			logger.Error(err)
		}

		msg := make([]byte, 1024)
		n, _ := conn.Read(msg)
		logger.Info(string(msg[:n]))
	}
}

func TestGet(t *testing.T) {
	replies := []resp.Reply{
		reply.MakeMultiBulkReply([][]byte{
			[]byte("get"),
			[]byte("hello"),
		}),
		reply.MakeMultiBulkReply([][]byte{
			[]byte("get"),
			[]byte("key"),
		}),
	}

	for _, replie := range replies {
		_, err = conn.Write(replie.ToBytes())
		if err != nil {
			logger.Error(err)
		}

		msg := make([]byte, 1024)
		n, _ := conn.Read(msg)
		logger.Info(string(msg[:n]))
	}
}
