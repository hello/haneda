package main

import (
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/core"
	"github.com/hello/haneda/sense"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

type FakeKeyStore struct {
	senseId string
}

func (k *FakeKeyStore) Get(senseId string) ([]byte, error) {
	return []byte("1234"), nil
}

type NoopBridge struct {
	sync.Mutex
	calls map[string]int
}

func (b *NoopBridge) register(path string) {
	b.Lock()
	defer b.Unlock()
	val, found := b.calls[path]
	if !found {
		val = 0
	}
	b.calls[path] = val + 1
}

func (b *NoopBridge) check(path string) int {
	b.Lock()
	defer b.Unlock()
	val, found := b.calls[path]
	if !found {
		return 0
	}
	return val
}

func (b *NoopBridge) PeriodicData(message *api.BatchedPeriodicData, privKey []byte) ([]byte, error) {
	return []byte{}, nil
}

func (b *NoopBridge) Logs(message *api.SenseLog, privKey []byte) error {
	return nil
}

var (
	senseId = sense.SenseId("name")
	ks      = &FakeKeyStore{
		senseId: string(senseId),
	}
)

func TestConnectToWebSocketServer(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	bridge := &NoopBridge{}
	simple := core.NewSimpleHelloServer(bridge, "example", nil, done, messages, ks)
	wsHandler := core.NewSimpleWsHandler(simple)

	ts := httptest.NewServer(wsHandler)
	defer ts.Close()
	wsurl, _ := url.Parse(ts.URL)
	wsurl.Scheme = "ws"
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := sense.New15(senseId, interrupt, done, []byte("1234"))

	err := client.Connect(wsurl, http.Header{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	err = client.Disconnect()
	if err != nil {
		t.Errorf("Error disconnecting: %v", err)
		t.FailNow()
	}

	time.Sleep(5 * time.Millisecond)

	if bridge.check("logs") != 0 {
		t.Errorf("%s", "Should not have called path")
	}
	simple.Shutdown()
}

func TestWrite(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(&NoopBridge{}, "example", nil, done, messages, ks)
	wsHandler := core.NewSimpleWsHandler(simple)

	ts := httptest.NewServer(wsHandler)
	defer ts.Close()
	wsurl, _ := url.Parse(ts.URL)
	wsurl.Scheme = "ws"
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := sense.New15(senseId, interrupt, done, []byte("1234"))
	err := client.Connect(wsurl, http.Header{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	client.Write([]byte{})
	if err != nil {
		t.Errorf("Error Writing: %v", err)
		t.FailNow()
	}
	err = client.Disconnect()
	if err != nil {
		t.Errorf("Error disconnecting: %v", err)
		t.FailNow()
	}
	simple.Shutdown()
}
