package main

import (
	"io"
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
	ctx := r.Context()

	go func() {
		for ctx.Err() == nil {
			_, r, err := ws.Reader(ctx)
			if err != nil {
				log.Println(err)
				return
			}
			data, err := io.ReadAll(r)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(string(data))
		}
	}()

	// Loop until the context is canceled
	for ctx.Err() == nil {
		wsjson.Write(ctx, ws, Thing{TimeStamp: time.Now()})
		time.Sleep(3 * time.Second)
	}
	log.Printf("Closing connection %v, reason: %v", r.RemoteAddr, ctx.Err())
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
