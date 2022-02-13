package main

import (
	"log"
	"net/http"
	"time"
)

type Tracker struct {
}

func (tracker *Tracker) Event(data string) {
	time.Sleep(time.Millisecond)
	log.Panicln(data)
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
