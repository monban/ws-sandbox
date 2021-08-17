package main

import (
	"context"
	//"encoding/json"
	"io"
	"log"
	"time"

	"nhooyr.io/websocket"
)

const sleepTimeout = 3 * time.Second

type Thing struct {
	TimeStamp time.Time
}

func main() {
	ctx := context.Background()
	log.Println("Starting...")
	run(ctx)
}

func run(ctx context.Context) {
	for {
		ws, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
		if err != nil {
			log.Printf("Socket failed to connect: %v", err)
			time.Sleep(sleepTimeout)
			continue
		}
		log.Println("Socket connected")
		in := readStuff(ctx, ws)
		for {
			select {
			case data := <- in:
				log.Println(data)
			case <-time.After(1*time.Second):
				ws.Write(ctx, 1, []byte("hello"))
			case <-ctx.Done():
				log.Printf("Socket closing: %v", ctx.Err())
				return
			}
		}
	}
}

func readStuff(ctx context.Context, ws *websocket.Conn) (<-chan string) {
	//var thing Thing
	c := make(chan string)
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
			c <- string(data)
		}
	}()
	return c
}
