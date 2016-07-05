package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/garyburd/redigo/redis"
	"github.com/hello/haneda/core"
	"github.com/hello/haneda/sense"
	config "github.com/stvp/go-toml-config"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	configPath = flag.String("c", "server.conf", "Path to config file. Ex: kenko -c /etc/hello/kenko.conf")
)

var (
	serverExternalHost = config.String("server.external_host", ":8082")
	serverInternalHost = config.String("server.internal_host", ":8082")
	proxyEndpoint      = config.String("proxy.endpoint", "http://localhost:5555")
	redisHost          = config.String("redis.host", ":6379")
	pubSubKey          = config.String("redis.pubsub", "example")
	awsKey             = config.String("aws.key", "")
	awsSecret          = config.String("aws.secret", "")
	awsRegion          = config.String("aws.region", "")
	keyStoreTable      = config.String("aws.keystore_table", "")
)

func proxy(w http.ResponseWriter, r *http.Request) {

}

type PublishHanlder struct {
	pool  *redis.Pool
	topic string
}

func (h *PublishHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := h.pool.Get()
	msg := `{"name": "Sense1", "id": 1234567890}`
	conn.Do("PUBLISH", h.topic, msg)
	fmt.Fprintf(w, "%s", msg)
}

func webserver(topic string, pool *redis.Pool, messages chan *sense.MessageParts) {

	ph := &PublishHanlder{
		topic: topic,
		pool:  pool,
	}

	http.Handle("/publish", ph)
	err := http.ListenAndServe(*serverInternalHost, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/health")
	fmt.Fprintf(w, "%s\n", "ok")
}

func main() {
	flag.Parse()

	if err := config.Parse(*configPath); err != nil {
		log.Printf("[haneda-server] can't find configuration: %s\n", *configPath)
		log.Fatal(err)
	}

	log.Printf("[haneda-server] Configuration loaded from: %s\n", *configPath)
	msg := "[haneda-server] Configured to proxy requests to: %s.\n"
	log.Printf(msg, *proxyEndpoint)

	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", *redisHost)

		if err != nil {
			return nil, err
		}

		return c, err
	}, 10)

	defer redisPool.Close()
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 2)

	config := &aws.Config{}
	config.Region = awsRegion

	if *awsKey != "" && *awsSecret != "" {
		config.Credentials = credentials.NewStaticCredentials(*awsKey, *awsSecret, "")
	}

	ks := sense.NewDynamoDBKeyStore(*keyStoreTable, config)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	forwarder := core.NewHttpForwarder(*proxyEndpoint, client)

	simple := core.NewSimpleHelloServer(forwarder, *pubSubKey, redisPool, done, messages, ks)
	go simple.Start()
	go webserver(*pubSubKey, redisPool, messages)

	http.HandleFunc("/health", HealthHandler)
	http.Handle("/protobuf", simple)

	err := http.ListenAndServe(*serverExternalHost, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
