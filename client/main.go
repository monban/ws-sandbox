package main

import (
	"context"
	"encoding/json"
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
		if err := readStuff(ctx, ws); err != nil {
			log.Println(err)
		}
		time.Sleep(sleepTimeout)
	}
}

func readStuff(ctx context.Context, ws *websocket.Conn) error {
	var thing Thing
	for {
		ctx_to, _ := context.WithTimeout(ctx, 5*time.Second)
		mtype, r, err := ws.Reader(ctx_to)
		if err != nil {
			return err
		}

		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &thing); err != nil {
			return err
		}

		log.Printf("%d %v\n", mtype, thing)
	}
}
