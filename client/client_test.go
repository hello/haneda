package sense

import (
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/core"
	"github.com/hello/haneda/sense"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

var (
	senseId = sense.SenseId("name")
	key     = []byte("qwertyuiop")
	ks      = &sense.FakeKeyStore{SenseId: string(senseId), Key: key}
	conf    = &core.HelloConfig{
		Redis: &core.RedisConfig{
			PubSub: "example",
		},
	}
)

type FakePeriodicDataBrigde struct {
	core.NoopBridge
	generator sense.MessageGenerator
}

func (b *FakePeriodicDataBrigde) Pair(message *api.MorpheusCommand, key []byte) ([]byte, error) {
	sr, err := b.generator.Do(uint64(1))
	return sr.Body, err
}

func TestConnectDisconnect(t *testing.T) {
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(ioutil.Discard, &core.NoopBridge{}, nil, done, messages, ks, conf)
	go simple.Start()
	defer simple.Shutdown()
	ts := httptest.NewServer(simple)
	defer ts.Close()
	wsurl, _ := url.Parse(ts.URL)
	wsurl.Scheme = "ws"
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := sense.New15(senseId, interrupt, done, key)

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
}

func TestReadTimeout(t *testing.T) {
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(ioutil.Discard, &core.NoopBridge{}, nil, done, messages, ks, conf)
	go simple.Start()
	defer simple.Shutdown()

	ts := httptest.NewServer(simple)
	defer ts.Close()
	wsurl, _ := url.Parse(ts.URL)
	wsurl.Scheme = "ws"
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := sense.New15(senseId, interrupt, done, key)
	defer client.Disconnect()
	err := client.Connect(wsurl, http.Header{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	_, err = client.Read(1 * time.Second)
	if err == nil {
		t.Fatal("client.Read: expected timeout error")
	}
}

func TestPair(t *testing.T) {
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)
	token := "1"
	bridge := &FakePeriodicDataBrigde{
		generator: &sense.PairingMessageGenerator{
			Token:        token,
			DeviceId:     senseId,
			MorpheusType: api.MorpheusCommand_MORPHEUS_COMMAND_PAIR_SENSE.Enum(),
		},
	}

	simple := core.NewSimpleHelloServer(ioutil.Discard, bridge, nil, done, messages, ks, conf)
	go simple.Start()
	defer simple.Shutdown()

	ts := httptest.NewServer(simple)
	defer ts.Close()
	wsurl, _ := url.Parse(ts.URL)
	wsurl.Scheme = "ws"

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	device := sense.New15(senseId, interrupt, done, key)
	s15 := sense.NewSenseOneFive("1", string(senseId), device)

	defer s15.Unplug()
	err := s15.Plug(wsurl, http.Header{})
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	pairErr := s15.Pair()
	if pairErr != nil {
		t.Errorf("pair: %v", pairErr)
	}
}
