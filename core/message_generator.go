package core

import (
	"github.com/hello/haneda/sense"
)

type MessageGenerator interface {
	Generate() chan<- *sense.MessageParts
}

type RandomMessageGenerator struct {
	messages chan *sense.MessageParts
}

func (g *RandomMessageGenerator) Generate() chan<- *sense.MessageParts {
	return g.messages
}
