package tcp

import (
	"context"
	"fmt"
	"github.com/LCY2013/thinking-in-go/gedis/interface/tcp"
	"github.com/LCY2013/thinking-in-go/gedis/lib/async/run"
	"github.com/LCY2013/thinking-in-go/gedis/lib/logger"
	"github.com/LCY2013/thinking-in-go/gedis/lib/shutdown"
	"github.com/sourcegraph/conc"
	"net"
	"os"
	"os/signal"
	"time"
)

/*
A TCP Connection Server
*/

// Config is the Tcp Server configuration
type Config struct {
	Address string
}

// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, shutdown.Signals...)

	run.GO(func() {
		<-sigChan

		select {
		case <-sigChan:
			logger.Warn("kill tcp server now")
			os.Exit(1)
		case <-time.After(15 * time.Second):
		}

		closeChan <- struct{}{}
	})

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan chan struct{}) {
	// close during unexpected error
	defer func() {
		_ = listener.Close() // listener.Accept() will return err immediately
		_ = handler.Close()  // close connections
	}()

	// watch signal for system
	run.GO(func() {
		<-closeChan

		_ = listener.Close() // listener.Accept() will return err immediately
		_ = handler.Close()  // close connections
	})

	// listen port
	var wait conc.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}

		// handle
		logger.Info(fmt.Sprintf("accept new connection: %s", conn.RemoteAddr().String()))
		wait.Go(func() {
			handler.Handle(context.Background(), conn)
		})
	}
	wait.Wait()
}
