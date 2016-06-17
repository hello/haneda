package haneda

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type WSHandler struct {
	ventilator  *Ventilator
	helloServer *HelloServer
}

func NewWsHandler(vent *Ventilator, server *HelloServer) *WSHandler {
	return &WSHandler{
		ventilator:  vent,
		helloServer: server,
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
	senseConn := &SenseConn{
		SenseId:               sense,
		Conn:                  conn,
		TopFirmwareVersion:    "top",         // get from headers
		MiddleFirmwareVersion: "middle",      // get from headers
		PrivKey:               []byte("abc"), // get from keystore
	}

	h.ventilator.Add(senseConn)
	go h.helloServer.Spin(senseConn)
}
