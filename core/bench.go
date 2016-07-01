package core

import (
	"fmt"
	"github.com/hello/haneda/sense"
	"log"
	"net/http"
)

type BenchServer struct {
	Messages       chan *sense.MessageParts
	SignedMessages chan []byte
	Key            []byte
	Bridge         Bridge
	Remover        ConnectionRemover
}

func (s *BenchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving http req")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("here", err)
		return
	}

	senseId := sense.SenseId("fake")

	senseConn := &SenseConn{
		SenseId:               senseId,
		Conn:                  conn,
		TopFirmwareVersion:    sense.TopFirmwareVersion("top"),       // get from headers
		MiddleFirmwareVersion: sense.MiddleFirmwareVersion("middle"), // get from headers
		PrivKey:               s.Key,
		out:                   s.Messages,
		internalMsgs:          s.SignedMessages,
		bridge:                s.Bridge,
		remover:               s.Remover,
	}

	go senseConn.Serve()
}

func (s *BenchServer) Start() {
	log.Println("Bench server started")
}
