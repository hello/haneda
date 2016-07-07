package core

import (
	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/graphite"
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
	logger   log.Logger
	metrics  *graphite.Emitter
	stats    chan *HelloStat
}

type HelloStat struct {
	ErrRead  *uint64
	OkRead   *uint64
	ErrParse *uint64
	ErrProxy *uint64
}

func hInt64(v int) *uint64 {
	p := new(uint64)
	*p = uint64(v)
	return p
}
