package core

import (
	"github.com/go-kit/kit/log"
	"github.com/hello/haneda/sense"
)

type ConnectionRemover interface {
	Remove(senseId sense.SenseId)
}

type ConnectionAdder interface {
	Add(senseId sense.SenseId) chan *sense.MessageParts
}

type Hub struct {
	removeChan     chan sense.SenseId
	chanBufferSize int
	logger         log.Logger
}

func (h *Hub) Remove(senseId sense.SenseId) {
	h.logger.Log("action", "removing", "sense_id", senseId)
}

func (h *Hub) Add(senseId sense.SenseId) chan *sense.MessageParts {
	h.logger.Log("action", "adding", "sense_id", senseId)
	return make(chan *sense.MessageParts, h.chanBufferSize)
}
