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

type Writer interface {
	Write(context.Context, websocket.MessageType, []byte) error
}

func main() {
	ctx := context.Background()
	var err error
	var ws *websocket.Conn
	log.Println("Starting...")
	// Start an infinite loop which sets up a socket and begins communicating on
	// it. If the connection is lost, start over and make a new one.
	for {
		ws, _, err = websocket.Dial(ctx, "ws://localhost:8080", nil)
		if err != nil {
			log.Printf("Socket failed to connect: %v", err)
			time.Sleep(sleepTimeout)
			continue
		}
		log.Println("Socket connected")
		in := readStuff(ctx, ws)
		err = run(ctx, ws, in)
		if err != nil {
			log.Printf("Socket closed: %v", err)
			ws.Close(websocket.StatusInternalError, err.Error())
		}
	}
}

func run(ctx context.Context, ws Writer, in <-chan string) error {
	for {
		select {
		case data := <-in:
			log.Println(data)
		case <-time.After(1 * time.Second):
			err := ws.Write(ctx, 1, []byte("hello"))
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func readStuff(ctx context.Context, ws *websocket.Conn) <-chan string {
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
