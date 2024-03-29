package core

import (
	"github.com/go-kit/kit/log"
	"github.com/hello/haneda/sense"
	"io"
	"io/ioutil"
	"net/http"
)

type BenchServer struct {
	Messages       chan *sense.MessageParts
	SignedMessages chan []byte
	Key            []byte
	Bridge         Bridge
	Remover        ConnectionRemover
	Logger         log.Logger
}

func loggers(w io.Writer) *Loggers {
	logger := log.NewLogfmtLogger(w)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "app", "bench")
	return &Loggers{
		Debug: logger,
		Info:  logger,
		Warn:  logger,
		Error: logger,
	}
}
func (s *BenchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.Logger.Log("Serving http req")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Log("error", err)
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
		loggers:               loggers(ioutil.Discard),
	}
	stats := make(chan *HelloStat, 10)
	go senseConn.Serve(stats)
}

func (s *BenchServer) Start() {
	s.Logger.Log("Bench server started")
}
