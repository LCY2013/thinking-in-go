package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Tracker struct {
	wg sync.WaitGroup
}

func (tracker *Tracker) Event(data string) {
	tracker.wg.Add(1)
	go func() {
		defer tracker.wg.Done()
		time.Sleep(time.Millisecond)
		log.Panicln(data)
	}()
}

func (tracker *Tracker) Shutdown() {
	tracker.wg.Wait()
}

type App struct {
	tracker Tracker
}

func (a *App) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	go a.tracker.Event("this event")
}

func main() {

}
