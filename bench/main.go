package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/core"
	"github.com/hello/haneda/sense"
	config "github.com/stvp/go-toml-config"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	configPath = flag.String("c", "server.conf", "Path to config file. Ex: kenko -c /etc/hello/kenko.conf")
	serverOnly = flag.Bool("server", false, "server=1 enables client")
)

var (
	serverExternalHost = config.String("server.external_host", ":8082")
	proxyEndpoint      = config.String("proxy.endpoint", "http://localhost:5555")
)

type BenchClient struct {
	auth   sense.MessageSigner
	funcs  []genFunc
	logger log.Logger
}

type genFunc func(msgId uint64) (*sense.MessageParts, error)

func (b *BenchClient) genRandomMessage(i int) ([]byte, error) {
	var mp *sense.MessageParts
	var err error

	r := rand.Int31n(int32(len(b.funcs)))

	f := b.funcs[r]
	mp, err = f(uint64(i))

	if err != nil {
		panic(err)
	}
	return b.auth.Sign(mp)
}

func (b *BenchClient) Start(endpoint string, in chan []byte, tickDuration time.Duration) {
	wc, _, err := websocket.DefaultDialer.Dial(endpoint, http.Header{})
	if err != nil {
		b.logger.Log("err", err)
		return
	}
	done := make(chan bool, 0)
	go func(c *websocket.Conn, done chan bool) {
		b.logger.Log("action", "start_reading")
		i := 0
		totalBytes := 0
		for {
			_, content, err := c.ReadMessage()
			if err != nil {
				b.logger.Log("error", err)
				done <- false
			}
			totalBytes += len(content)
			if i%100 == 0 {
				b.logger.Log("iteration", i)
				b.logger.Log("total", totalBytes)
			}
			i++
		}
	}(wc, done)

	tick := time.NewTicker(tickDuration)
	timeout := time.NewTimer(100 * time.Second)
	i := 1
outer:
	for {
		select {
		case <-done:
			b.logger.Log("interrupting", true)
			break outer
		case <-tick.C:
			m, err := b.genRandomMessage(i)
			i++
			if err != nil {
				b.logger.Log("error", err)
				break outer
			}
			in <- m
		case <-timeout.C:
			b.logger.Log("timeout", true)
			break outer
		}
	}
	b.logger.Log("done", true)
}

func main() {
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	if err := config.Parse(*configPath); err != nil {
		logger.Log("error", err)
	}

	logger.Log("proxy_endpoint", *proxyEndpoint)

	messages := make(chan *sense.MessageParts, 2)
	signedMessages := make(chan []byte, 2)

	bench := &core.BenchServer{
		Messages:       messages,
		SignedMessages: signedMessages,
		Logger:         logger,
	}

	go bench.Start()

	privKey, _ := hex.DecodeString("AD332E8DFE33490AAF35CA2824ECADC0")

	manifestGenerator := &sense.FileManifestGenerator{}

	bc := &BenchClient{
		auth:   sense.NewSenseAuthHmacSha1(privKey, sense.SenseId("whatever")),
		funcs:  []genFunc{sense.GenSyncResp, sense.GenMesseji, manifestGenerator.Generate},
		logger: logger,
	}

	wsPath := "/protobuf"
	http.Handle(wsPath, bench)

	go func() {
		err := http.ListenAndServe(*serverExternalHost, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	if !*serverOnly {
		time.Sleep(2 * time.Second)
		bc.Start("ws://"+*serverExternalHost+wsPath, signedMessages, 1*time.Second)
	} else {
		fmt.Println("block forever, server mode")
		i := 0
		for {
			i++

			signed, err := bc.genRandomMessage(i)
			i++
			if err != nil {
				fmt.Println(err)
				break
			}
			signedMessages <- signed
			fmt.Println("Generated message ", i)
		}
	}
}
