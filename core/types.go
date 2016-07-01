package core

import (
	"github.com/garyburd/redigo/redis"
	"github.com/hello/haneda/sense"
	"sync"
)

type AuthenticateFunc func(username, password string) bool

type Listener interface {
	Listen(topic string)
}

type Sender interface {
	Send(message *sense.MessageParts)
}

type SimpleWsHandler struct {
	server *SimpleHelloServer
}

type SimpleHelloServer struct {
	sync.Mutex
	pairs    map[sense.SenseId]chan *sense.MessageParts
	bridge   Bridge
	done     chan bool
	pool     *redis.Pool
	topic    string
	messages chan *sense.MessageParts
	keystore sense.KeyStore
	adder    ConnectionAdder
	remover  ConnectionRemover
}
