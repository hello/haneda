package core

import (
	"fmt"
	"github.com/hello/haneda/sense"
)

type ConnectionRemover interface {
	Remove(senseId sense.SenseId)
}

type ConnectionAdder interface {
	Add(senseId sense.SenseId) chan *sense.MessageParts
}

type Hub struct {
	removeChan chan sense.SenseId
}

func (h *Hub) Remove(senseId sense.SenseId) {
	fmt.Println("removing", senseId)
}

func (h *Hub) Add(senseId sense.SenseId) chan *sense.MessageParts {
	fmt.Println("adding", senseId)
	return make(chan *sense.MessageParts, 1)
}
