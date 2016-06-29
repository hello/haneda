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
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{} // use default options

type SimpleWsHandler struct {
	server *SimpleHelloServer
}

func NewSimpleWsHandler(server *SimpleHelloServer) *SimpleWsHandler {
	return &SimpleWsHandler{
		server: server,
	}
}

func (h *SimpleWsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// sense, err := extractBasicAuth(r, checkCreds)
	// if err != nil {
	// 	log.Println("Bad auth")
	// 	http.Error(w, "authorization failed", http.StatusUnauthorized)
	// 	return
	// }

	sense := r.Header.Get("X-Hello-Sense-Id")
	log.Println(fmt.Sprintf("sense_id=%s ip_address=%s\n", sense, r.RemoteAddr))
	if sense == "" {
		http.Error(w, "Missing header with Sense ID", 400)
		return
	}

	key, err := h.server.keystore.Get(sense)
	if err != nil {
		log.Println("Couldn't connect to keystore for Sense", sense, err)
		http.Error(w, "Server error", 500)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("here", err)
		return
	}
	senseConn := &SenseConn{
		SenseId:               sense,
		Conn:                  conn,
		TopFirmwareVersion:    "top",    // get from headers
		MiddleFirmwareVersion: "middle", // get from headers
		PrivKey:               key,
	}
	c := h.server.register(sense)
	go h.server.Spin(senseConn, c)
}

type SimpleHelloServer struct {
	sync.Mutex
	pairs    map[string]chan *sense.MessageParts
	bridge   Bridge
	done     chan bool
	pool     *redis.Pool
	topic    string
	messages chan *sense.MessageParts
	keystore sense.KeyStore
}

func NewSimpleHelloServer(bridge Bridge, topic string, pool *redis.Pool, done chan bool, messages chan *sense.MessageParts, ks sense.KeyStore) *SimpleHelloServer {
	return &SimpleHelloServer{
		bridge:   bridge,
		done:     done,
		pool:     pool,
		topic:    topic,
		messages: messages,
		pairs:    make(map[string]chan *sense.MessageParts),
		keystore: ks,
	}
}

func (h *SimpleHelloServer) IsConnected(senseId string) bool {
	h.Lock()
	defer h.Unlock()
	_, found := h.pairs[senseId]
	return found
}

func (h *SimpleHelloServer) Start() {
	log.Println("SimpleHelloServer running")
	for {
		select {
		case m := <-h.messages:
			c, found := h.pairs[m.SenseId]
			if found {
				c <- m
			}
		}
	}
}

func (h *SimpleHelloServer) register(senseId string) chan *sense.MessageParts {
	h.Lock()
	defer h.Unlock()
	c := make(chan *sense.MessageParts)
	h.pairs[senseId] = c
	log.Println("added", senseId)
	return c
}

func (h *SimpleHelloServer) remove(senseId string) {
	h.Lock()
	defer h.Unlock()
	delete(h.pairs, senseId)
	log.Println("removed", senseId)
}

func write(s *SenseConn, sub chan *sense.MessageParts, self chan []byte) {
outer:
	for {
		select {
		case m := <-sub:
			log.Println("Sending Message to:", m.SenseId)
		case m := <-self: // assuming already fully assembled messages
			if err := s.Conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
				fmt.Println(err)
				break outer
			}
		}
	}
	log.Println("Writing thread stopped for sense", s.SenseId)
}

func dispatch(bridge Bridge, message *sense.MessageParts, s *SenseConn) ([]byte, error) {
	// bridge.Route(message)
	empty := make([]byte, 0)
	switch message.Header.GetType() {
	case haneda.Preamble_SENSE_LOG:
		m := &api.SenseLog{}
		proto.Unmarshal(message.Body, m)
		return empty, bridge.Logs(m, s)
	case haneda.Preamble_BATCHED_PERIODIC_DATA:
		m := &api.BatchedPeriodicData{}
		proto.Unmarshal(message.Body, m)
		return bridge.PeriodicData(m, s)
	default:
		// fmt.Println("Unknown", messageParts.Header.GetType().String())
	}
	return empty, nil
}

func (h *SimpleHelloServer) Spin(s *SenseConn, sub chan *sense.MessageParts) {
	if s.Conn == nil {
		fmt.Println("Can't spin ws thread, conn is nil:", s.SenseId)
		return
	}

	auth := sense.NewAuth(s.PrivKey, s.SenseId)

	i := 0
	defer s.Conn.Close()
	self := make(chan []byte, 0)

	go write(s, sub, self)
	for {
		_, content, err := s.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading.", err)
			break
		}

		mp, err := auth.Parse(content)
		if err != nil {
			log.Println(s.SenseId, err)
			break
		}

		// prepare ack message
		ack := &haneda.Ack{}
		ack.MessageId = proto.Uint64(mp.Header.GetId())
		ack.Status = haneda.Ack_SUCCESS.Enum()

		outbox := make([]*sense.MessageParts, 0)
		resp, err := dispatch(h.bridge, mp, s)
		if err != nil {
			// Override status since bridge responded with error
			ack.Status = haneda.Ack_CLIENT_REQUEST_ERROR.Enum()
			fmt.Println(mp.Header.GetId(), "Ack_CLIENT_REQUEST_ERROR")
			fmt.Println(err)
		}

		body, _ := proto.Marshal(ack)

		header := &haneda.Preamble{}
		header.Type = haneda.Preamble_ACK.Enum()

		out := &sense.MessageParts{
			Header: header,
			Body:   body,
		}

		outbox = append(outbox, out)

		// response from server might be happy
		if len(resp) > 0 {
			switch mp.Header.GetType() {
			case haneda.Preamble_BATCHED_PERIODIC_DATA:
				syncHeader := &haneda.Preamble{}
				syncHeader.Type = haneda.Preamble_SYNC_RESPONSE.Enum()
				syncHeader.Id = proto.Uint64(uint64(time.Now().UnixNano()))

				// syncResp := &api.SyncResponse{}
				// syncResp.RingTimeAck = proto.String("From proxy")

				// syncBody, _ := proto.Marshal(syncResp)
				out2 := &sense.MessageParts{
					Header: syncHeader,
					Body:   resp,
				}
				outbox = append(outbox, out2)
			default:
				log.Println("No response needed")
			}
		}

	outer:
		for _, mp := range outbox {
			serialized, err := auth.Sign(mp)
			if err != nil {
				log.Println(err)
				break outer
			}
			self <- serialized
		}
		i++
	}
	h.remove(s.SenseId)
	log.Println(s.SenseId, "Processed:", i)
}

// Listen blocks and wait for messages to be published on the redis channel
func (h *SimpleHelloServer) Listen() {

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
}
