package api

import (
	"fmt"
	"github.com/hello/haneda/haneda"
)

type Msg struct {
	Name string `json:"name"`
}

// JsonMessage is the message published by redis
type JsonMessage struct {
	Name string `json:"name"`
	Id   uint64 `json:"id"`
}

// MessageParts
type MessageParts struct {
	Header *haneda.Preamble
	Body   []byte
	Sig    []byte
}

// Stat holds statistics
type Stat struct {
	Source string
	Count  int
}

func (s Stat) String() string {
	return fmt.Sprintf("source:%s, count:%d\n", s.Source, s.Count)
}

// AuthenticateFunc validates username password against db
type AuthenticateFunc func(username, password string) bool
