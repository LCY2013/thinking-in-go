package client

import (
	"github.com/LCY2013/thinking-in-go/gedis/interface/resp"
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		logger.Error(err)
	}

	replies := []resp.Reply{
		reply.MakeMultiBulkReply([][]byte{
			[]byte("set"),
			[]byte("hello"),
			[]byte("fufeng"),
		}),
		reply.MakeMultiBulkReply([][]byte{
			[]byte("get"),
			[]byte("hello"),
		}),
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
