package main

import (
	"github.com/hello/haneda/core"
	"github.com/hello/haneda/sense"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/signal"
	"testing"
)

type FakeKeyStore struct {
	senseId string
}

func (k *FakeKeyStore) Get(senseId string) ([]byte, error) {
	return []byte("1234"), nil
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

	bridge := &core.NoopBridge{}
	simple := core.NewSimpleHelloServer(bridge, "example", nil, done, messages, ks)

	ts := httptest.NewServer(simple)
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
	simple.Shutdown()
}

func TestWrite(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(&core.NoopBridge{}, "example", nil, done, messages, ks)

	ts := httptest.NewServer(simple)
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
