package main

import (
  "nhooyr.io/websocket"
  "nhooyr.io/websocket/wsjson"
  "context"
  "log"
)

func main() {
  ctx := context.Background()
  ws, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
  if err != nil {
    log.Fatal(err)
    return
  }
  var data interface{}
  for {
    if err := wsjson.Read(ctx, ws, &data); err != nil {
      log.Fatal(err)
    }
  }
}
