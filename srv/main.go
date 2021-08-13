package main

import (
  "net/http"
  "nhooyr.io/websocket"
  "time"
  "log"
  //"context"
)

type wssrv struct {
}

func (s wssrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Printf("Upgrading %v to websocket", r.RemoteAddr)
  ws, err := websocket.Accept(w, r, nil)
  if err != nil {
    log.Println(err.Error())
  }
  ws.CloseRead(r.Context())
}

func main() {
  log.Printf("Starting...")
  if err := run(); err != nil {
    log.Fatalf(err.Error())
  }
}

func run() error {
  s := &http.Server{
    Addr: ":8080",
    Handler: wssrv{ },
    ReadTimeout:  time.Second * 10,
    WriteTimeout: time.Second * 10,
  }
  return s.ListenAndServe()
}
