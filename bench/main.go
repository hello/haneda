package main

import (
	"flag"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/core"
	"github.com/hello/haneda/haneda"
	"github.com/hello/haneda/sense"
	config "github.com/stvp/go-toml-config"
	"log"
	"net/http"
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
	auth sense.MessageSigner
}

func (b *BenchClient) genRandomMessage(i int) ([]byte, error) {
	var mp *sense.MessageParts
	var err error

	if i%2 == 0 {
		mp, err = periodic(uint64(i))
	} else {
		mp, err = logs(uint64(i))
	}
	if err != nil {
		panic(err)
	}
	return b.auth.Sign(mp)

}

func logs(messageId uint64) (*sense.MessageParts, error) {
	pb := &haneda.Preamble{}
	pb.Type = haneda.Preamble_SENSE_LOG.Enum()
	pb.Id = proto.Uint64(messageId)

	sLog := &api.SenseLog{}
	combined := fmt.Sprintf("Log #%d", messageId)
	sLog.Text = &combined
	n := sense.SenseId("bench-client")
	sLog.DeviceId = proto.String(string(n))
	body, err := proto.Marshal(sLog)
	if err != nil {
		return nil, err
	}
	mp := &sense.MessageParts{
		Header:  pb,
		Body:    body,
		SenseId: n,
	}
	return mp, nil
}

func periodic(messageId uint64) (*sense.MessageParts, error) {
	header := &haneda.Preamble{}
	header.Type = haneda.Preamble_BATCHED_PERIODIC_DATA.Enum()
	header.Id = proto.Uint64(messageId)

	batched := &api.BatchedPeriodicData{}
	periodic := &api.PeriodicData{}
	periodic.Temperature = proto.Int32(3500)

	n := string("bench-client")
	batched.DeviceId = &n
	batched.FirmwareVersion = proto.Int32(888)
	batched.Data = append(batched.Data, periodic)

	body, pbErr := proto.Marshal(batched)
	if pbErr != nil {
		return nil, pbErr
	}

	mp := &sense.MessageParts{
		Header:  header,
		Body:    body,
		SenseId: sense.SenseId(n),
	}
	return mp, nil
}

func (c *BenchClient) Start(endpoint string, in chan []byte, tickDuration time.Duration) {
	wc, _, err := websocket.DefaultDialer.Dial(endpoint, http.Header{})
	if err != nil {
		fmt.Println(err)
		return
	}
	done := make(chan bool, 0)
	go func(c *websocket.Conn, done chan bool) {
		fmt.Println("starting reading")
		for {
			_, content, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				done <- false
			}
			fmt.Println("len:", len(content))
		}
	}(wc, done)

	tick := time.NewTicker(tickDuration)
	timeout := time.NewTimer(10 * time.Second)
	i := 1
outer:
	for {
		select {
		case b := <-done:
			fmt.Println("Done. interrupting", b)
			break outer
		case <-tick.C:
			m, err := c.genRandomMessage(i)
			i++
			if err != nil {
				fmt.Println(err)
				break outer
			}
			in <- m
		case <-timeout.C:
			fmt.Println("Timeout.")
			break outer
		}
	}
	fmt.Println("Done")
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

	messages := make(chan *sense.MessageParts, 2)
	signedMessages := make(chan []byte, 2)

	bench := &core.BenchServer{
		Messages:       messages,
		Bridge:         &core.NoopBridge{},
		SignedMessages: signedMessages,
	}

	go bench.Start()

	bc := &BenchClient{
		auth: sense.NewAuth([]byte("1234567891234567"), sense.SenseId("whatever")),
	}
	wsPath := "/bench"
	http.Handle(wsPath, bench)

	go func() {
		err := http.ListenAndServe(*serverExternalHost, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
	if !*serverOnly {
		time.Sleep(2 * time.Second)
		bc.Start("ws://"+*serverExternalHost+wsPath, signedMessages, 100*time.Millisecond)
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
