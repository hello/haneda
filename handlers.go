package haneda

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"net/http"
)

type WSHandler struct {
	stats      chan api.Stat
	ventilator *Ventilator
}

func NewWsHandler(stats chan api.Stat, vent *Ventilator) *WSHandler {
	return &WSHandler{
		stats:      stats,
		ventilator: vent,
	}
}

func (h WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// fmt.Printf("origin: %s\n", r.Header.Get("Origin"))

	// fmt.Printf("%v\n", r.Header)
	sense, err := extractBasicAuth(r, checkCreds)

	if err != nil {
		fmt.Println("Bad auth")
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}
	fmt.Println(sense)
	// fmt.Printf("Method: %s\n", r.Method)
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	h.ventilator.register(sense, conn)
	go echo(conn, sense, h.stats)
}
