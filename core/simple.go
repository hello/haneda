package core

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"github.com/hello/haneda/sense"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func (s *SimpleHelloServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// sense, err := extractBasicAuth(r, checkCreds)
	// if err != nil {
	// 	log.Println("Bad auth")
	// 	http.Error(w, "authorization failed", http.StatusUnauthorized)
	// 	return
	// }

	connectedSense := r.Header.Get("X-Hello-Sense-Id")
	log.Println(fmt.Sprintf("sense_id=%s ip_address=%s\n", connectedSense, r.RemoteAddr))
	if connectedSense == "" {
		http.Error(w, "Missing header with Sense ID", 400)
		return
	}

	key, err := s.keystore.Get(connectedSense)
	if err != nil {
		log.Println("Couldn't connect to keystore for Sense", connectedSense, err)
		http.Error(w, "Server error", 500)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("here", err)
		return
	}

	senseId := sense.SenseId(connectedSense)
	c := s.adder.Add(senseId)

	auth := sense.NewSenseAuthHmacSha1(key, senseId)

	senseConn := &SenseConn{
		SenseId:               senseId,
		Conn:                  conn,
		TopFirmwareVersion:    sense.TopFirmwareVersion("top"),       // get from headers
		MiddleFirmwareVersion: sense.MiddleFirmwareVersion("middle"), // get from headers
		PrivKey:               key,
		out:                   c,
		internalMsgs:          make(chan []byte, 0),
		bridge:                s.bridge,
		remover:               s.remover,
		signer:                auth,
		parser:                auth,
	}

	go senseConn.Serve()
}

func (s *SimpleHelloServer) Shutdown() {
	s.Lock()
	defer s.Unlock()
	s.pairs = make(map[sense.SenseId]chan *sense.MessageParts)
	// h.pool.Close()
	close(s.done)
	close(s.messages)
}

func NewSimpleHelloServer(bridge Bridge, topic string, pool *redis.Pool, done chan bool, messages chan *sense.MessageParts, ks sense.KeyStore) *SimpleHelloServer {
	hub := &Hub{
		removeChan: make(chan sense.SenseId, 2),
	}

	return &SimpleHelloServer{
		bridge:   bridge,
		done:     done,
		pool:     pool,
		topic:    topic,
		messages: messages,
		pairs:    make(map[sense.SenseId]chan *sense.MessageParts),
		keystore: ks,
		adder:    hub,
		remover:  hub,
	}
}

func (s *SimpleHelloServer) Start() {
	log.Println("SimpleHelloServer running")
	for {
		select {
		case m := <-s.messages:
			c, found := s.pairs[m.SenseId]
			if found {
				c <- m
			} else {
				fmt.Println("%s", m.SenseId)
			}
		}
	}
}

func (s *SimpleHelloServer) Register(senseId sense.SenseId) chan *sense.MessageParts {
	s.Lock()
	defer s.Unlock()
	c := make(chan *sense.MessageParts)
	s.pairs[senseId] = c
	log.Println("added", senseId)
	return c
}

func (s *SimpleHelloServer) Remove(senseId sense.SenseId) {
	s.Lock()
	defer s.Unlock()
	delete(s.pairs, senseId)
	log.Println("removed", senseId)
}

func dispatch(bridge Bridge, message *sense.MessageParts, s *SenseConn) ([]byte, error) {
	// bridge.Route(message)
	empty := make([]byte, 0)
	switch message.Header.GetType() {
	case haneda.Preamble_SENSE_LOG:
		m := &api.SenseLog{}
		proto.Unmarshal(message.Body, m)
		return empty, bridge.Logs(m, s.PrivKey)
	case haneda.Preamble_BATCHED_PERIODIC_DATA:
		m := &api.BatchedPeriodicData{}
		proto.Unmarshal(message.Body, m)
		return bridge.PeriodicData(m, s.PrivKey)
	default:
		// fmt.Println("Unknown", messageParts.Header.GetType().String())
	}
	return empty, nil
}

// Listen blocks and wait for messages to be published on the redis channel

// for {
// 	// Get a connection from a pool
// 	c := h.pool.Get()
// 	psc := redis.PubSubConn{c}

// 	// Set up subscriptions
// 	psc.Subscribe(h.topic)

// 	// While not a permanent error on the connection.
// 	for c.Err() == nil {
// 		switch val := psc.Receive().(type) {
// 		case redis.Message:
// 			parsed, err := parse(val.Data)
// 			if err != nil {
// 				log.Println("Failed to parse message", err)
// 				continue
// 			}
// 			v.Receive(parsed)
// 		case redis.Subscription:
// 			fmt.Printf("%s: %s %d\n", val.Channel, val.Kind, val.Count)
// 		case error:
// 			fmt.Printf("%v", val)
// 		}
// 	}
// 	c.Close()
// }
