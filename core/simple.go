package core

import (
	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/graphite"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"github.com/hello/haneda/sense"
	"io"
	"net/http"
	"time"
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
	s.logger.Log("action", "connection-attempt", "sense_id", connectedSense, "ip_address", r.RemoteAddr)
	if connectedSense == "" {
		s.logger.Log("error", "missing_header", "ip_address", r.RemoteAddr)
		http.Error(w, "Missing header with Sense ID", 400)
		return
	}

	key, err := s.keystore.Get(connectedSense)
	if err != nil {
		s.logger.Log("action", "key_store_connect", "error", err, "sense_id", connectedSense)
		http.Error(w, "Server error", 500)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Log("error", err, "reason", "failed_upgrading_ws")
		return
	}

	senseId := sense.SenseId(connectedSense)
	c := s.adder.Add(senseId)

	auth := sense.NewSenseAuthHmacSha1(key, senseId)

	logger := log.NewContext(s.logger).With("ts", log.DefaultTimestampUTC, "sense_id", senseId, "ip_address", r.RemoteAddr)

	senseConn := &SenseConn{
		SenseId:               senseId,
		Conn:                  conn,
		TopFirmwareVersion:    sense.TopFirmwareVersion("top"),       // get from headers
		MiddleFirmwareVersion: sense.MiddleFirmwareVersion("middle"), // get from headers
		PrivKey:               key,
		out:                   c,
		internalMsgs:          make(chan []byte, 2),
		bridge:                s.bridge,
		remover:               s.remover,
		signer:                auth,
		parser:                auth,
		logger:                logger,
	}

	go senseConn.Serve(s.stats)
}

func (s *SimpleHelloServer) Shutdown() {
	s.Lock()
	defer s.Unlock()
	s.pairs = make(map[sense.SenseId]chan *sense.MessageParts)
	// h.pool.Close()
	close(s.done)
	close(s.messages)
}

func NewSimpleHelloServer(w io.Writer, bridge Bridge, pool *redis.Pool, done chan bool, messages chan *sense.MessageParts, ks sense.KeyStore, helloConf *HelloConfig) *SimpleHelloServer {

	logger := log.NewLogfmtLogger(w)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "app", "simple-hello-server")
	hubLogger := log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "app", "hub")

	logger.Log("making", "stats", "cap", 10)
	stats := make(chan *HelloStat, 10)

	hub := NewHub(hubLogger, stats)

	s := &SimpleHelloServer{
		bridge:   bridge,
		done:     done,
		pool:     pool,
		topic:    helloConf.Redis.PubSub,
		messages: messages,
		pairs:    make(map[sense.SenseId]chan *sense.MessageParts),
		keystore: ks,
		adder:    hub,
		remover:  hub,
		sender:   hub,
		logger:   logger,
		stats:    stats,
	}

	if helloConf.Graphite != nil {
		logger.Log("metrics", "enabled")
		graphiteLogger := log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "app", "graphite")
		s.metrics = graphite.NewEmitter("tcp", helloConf.Graphite.Host, helloConf.Graphite.Prefix+".", 5*time.Second, graphiteLogger)
	}
	logger.Log("simple", "starting")
	return s
}

func (s *SimpleHelloServer) Start() {
	s.logger.Log("server", "start")
	var errRead metrics.Counter
	var okRead metrics.Counter
	var errParse metrics.Counter
	var errProxy metrics.Counter

	var currentConnections metrics.Gauge
	var connDuration metrics.Histogram

	if s.metrics != nil {
		s.logger.Log("action", "running")
		errRead = s.metrics.NewCounter("err_read")
		okRead = s.metrics.NewCounter("ok_read")
		errParse = s.metrics.NewCounter("err_parse")
		errProxy = s.metrics.NewCounter("err_proxy")
		currentConnections = s.metrics.NewGauge("curr_conns")
		// 50, 90, 95, 99 represent the %ile we care about
		// sigfigs = 3, I have no clue what it does
		// max value is a year
		h, err := s.metrics.NewHistogram("conn_duration", 0, 3600*24*365, 3, 50, 90, 95, 99)
		if err != nil {
			panic(err)
		}
		connDuration = h
	}

	for {
		select {
		case m := <-s.messages:
			s.sender.Send(m)
		case stat := <-s.stats:
			if s.metrics != nil {
				if stat.ErrRead != nil {
					s.logger.Log("errRead", *stat.ErrRead)
					errRead.Add(*stat.ErrRead)
				}
				if stat.OkRead != nil {
					s.logger.Log("okRead", *stat.OkRead)
					okRead.Add(*stat.OkRead)
				}

				if stat.ErrParse != nil {
					s.logger.Log("errParse", *stat.ErrParse)
					errParse.Add(*stat.ErrParse)
				}

				if stat.ErrProxy != nil {
					s.logger.Log("errProxy", *stat.ErrProxy)
					errProxy.Add(*stat.ErrProxy)
				}

				if stat.CurrConns != nil {
					s.logger.Log("currConns", *stat.CurrConns)
					currentConnections.Set(*stat.CurrConns)
				}

				if stat.ConnDuration != nil {
					s.logger.Log("connDuration", *stat.ConnDuration)
					connDuration.Observe(*stat.ConnDuration)
				}
			}
		}
	}
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
	case haneda.Preamble_MORPHEUS_COMMAND:
		m := &api.MorpheusCommand{}
		err := proto.Unmarshal(message.Body, m)
		if err != nil {
			return empty, err
		}
		return bridge.Pair(m, s.PrivKey)
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
