package main

import (
	"flag"
	"fmt"
	toml "github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
	"github.com/hello/haneda/core"
	"github.com/hello/haneda/sense"
	"net/http"
	"os"
)

var (
	configPath = flag.String("c", "server.conf", "Path to config file. Ex: kenko -c /etc/hello/kenko.conf")
)

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

func webserver(host, topic string, pool *redis.Pool, messages chan *sense.MessageParts) {

	ph := &PublishHanlder{
		topic: topic,
		pool:  pool,
	}

	http.Handle("/publish", ph)
	err := http.ListenAndServe(host, nil)
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
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC, "app", "haneda-server")

	conf := &core.HelloConfig{}

	if _, err := toml.DecodeFile(*configPath, conf); err != nil {
		logger.Log("configuration-not-found", *configPath)
		return
	}
	logger.Log("action", "configuration-loaded", "path", *configPath)
	logger.Log("proxy", conf.Proxy.Endpoint)
	logger.Log("host", conf.Server.ExternalHost)

	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", conf.Redis.Host)

		if err != nil {
			return nil, err
		}

		return c, err
	}, 10)

	defer redisPool.Close()

	// done channel is used to stop the server gracefully
	done := make(chan bool, 0)

	// messages chan is used to push messages outside of the sense connection lifecycle
	messages := make(chan *sense.MessageParts, 2)

	awsConfig := &aws.Config{}
	awsConfig.Region = &conf.Aws.Region

	if conf.Aws.Key != "" && conf.Aws.Secret != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(conf.Aws.Key, conf.Aws.Secret, "")
	}

	ks := sense.NewDynamoDBKeyStore(conf.Aws.KeyStoreTableName, awsConfig)

	forwarder := core.NewDefaultHttpForwarder(conf.Proxy.Endpoint)

	simple := core.NewSimpleHelloServer(os.Stderr, forwarder, redisPool, done, messages, ks, conf)

	go simple.Start()

	// todo: implement messeji http handler API
	go webserver(conf.Server.InternalHost, conf.Redis.PubSub, redisPool, messages)

	http.HandleFunc("/health", HealthHandler)
	http.Handle("/protobuf", simple)

	err := http.ListenAndServe(conf.Server.ExternalHost, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
