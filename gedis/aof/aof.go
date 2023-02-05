package aof

import (
	"github.com/LCY2013/thinking-in-go/gedis/config"
	databaseface "github.com/LCY2013/thinking-in-go/gedis/interface/database"
	"github.com/LCY2013/thinking-in-go/gedis/lib/async/run"
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"github.com/LCY2013/thinking-in-go/gedis/lib/utils"
	"github.com/LCY2013/thinking-in-go/gedis/resp/connection"
	"github.com/LCY2013/thinking-in-go/gedis/resp/parser"
	"github.com/LCY2013/thinking-in-go/gedis/resp/reply"
	"io"
	"os"
	"strconv"
)

const (
	aofQueueSize = 1 << 16
)

// payload represents
type payload struct {
	cmdLine databaseface.CmdLine
	dbIndex int
}

// AofHandler receive msgs from channel and write to AOF file
type AofHandler struct {
	db          databaseface.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFileName string
	currentDB   int
}

// NewAOFHandler creates a new aof.AofHandler
func NewAOFHandler(db databaseface.Database) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFileName = config.Properties.AppendFilename
	handler.db = db

	// load aof to memory
	handler.loadAof()

	aofFile, err := os.OpenFile(handler.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	handler.aofChan = make(chan *payload, aofQueueSize)

	run.GO(func() {
		handler.handleAof()
	})

	return handler, nil
}

// handleAof listen aof channel and write into file
func (h *AofHandler) handleAof() {
	var (
		err error
	)
	// serialized execution
	h.currentDB = 0
	for p := range h.aofChan {
		if p.dbIndex != h.currentDB {
			// select db
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err = h.aofFile.Write(data)
			if err != nil {
				logger.Warn(err)
				continue // skip this command
			}
			h.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err = h.aofFile.Write(data)
		if err != nil {
			logger.Warn(err)
			continue // skip this command
		}
	}
}

// AddAof send command to aof goroutine through channel
func (h *AofHandler) AddAof(dbIndex int, cmdLine databaseface.CmdLine) {
	if config.Properties.AppendOnly && h.aofChan != nil {
		h.aofChan <- &payload{
			cmdLine: cmdLine,
			dbIndex: dbIndex,
		}
	}
}

// LoadAof read aof file
func (h *AofHandler) loadAof() {
	file, err := os.Open(h.aofFileName)
	if err != nil {
		logger.Warn(err)
		return
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Warn(err)
		}
	}(file)

	// only used for save dbIndex
	fakeConn := &connection.Connection{}
	ch := parser.ParseStream(file)
	for p := range ch {
		if p.Err != nil && p.Err == io.EOF {
			break
		}

		if p.Err != nil {
			logger.Error("parse error: " + p.Err.Error())
			continue
		}

		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}

		switch mbr := p.Data.(type) {
		case *reply.BulkReply:
			result := h.db.Exec(fakeConn, [][]byte{
				mbr.Arg,
			})
			if reply.IsErrorReply(result) {
				logger.Error("exec err", err)
			}
			continue
		case *reply.StatusReply:
			result := h.db.Exec(fakeConn, [][]byte{
				[]byte(mbr.Status),
			})
			if reply.IsErrorReply(result) {
				logger.Error("exec err", err)
			}
			continue
		case *reply.MultiBulkReply:
			result := h.db.Exec(fakeConn, mbr.Args)
			if reply.IsErrorReply(result) {
				logger.Error("exec err", err)
			}
			continue
		}
	}
}
