package main

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type wssrv struct {
}

type Thing struct {
	TimeStamp time.Time
}

func (s wssrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Upgrading %v to websocket", r.RemoteAddr)
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err.Error())
	}
	defer ws.Close(websocket.StatusNormalClosure, "goodbye")
	ctx := ws.CloseRead(r.Context())
	msgs := make(chan Thing)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msgs <- Thing{TimeStamp: time.Now()}:
				time.Sleep(3 * time.Second)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Closing connection: %v", ctx.Err())
			return

		case msg := <-msgs:
			wsjson.Write(r.Context(), ws, msg)
		}
	}
}

func main() {
	log.Printf("Starting...")
	go debugData()
	if err := run(); err != nil {
		log.Fatalf(err.Error())
	}
}

func run() error {
	s := &http.Server{
		Addr:         ":8080",
		Handler:      wssrv{},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	return s.ListenAndServe()
}

func debugData() {
	for {
		log.Printf("Currently running goroutines: %v", runtime.NumGoroutine())
		time.Sleep(5 * time.Second)
	}
}
