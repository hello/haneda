package sense

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
	"time"
)

var (
	senseId = sense.SenseId("name")
	ks      = &sense.FakeKeyStore{SenseId: string(senseId), Key: []byte("abc")}
	conf    = &core.HelloConfig{
		Redis: &core.RedisConfig{
			PubSub: "example",
		},
	}
)

func TestConnectDisconnect(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(&core.NoopBridge{}, nil, done, messages, ks, conf)

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

func TestReadTimeout(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(&core.NoopBridge{}, nil, done, messages, ks, conf)
	defer simple.Shutdown()

	ts := httptest.NewServer(simple)
	defer ts.Close()
	wsurl, _ := url.Parse(ts.URL)
	wsurl.Scheme = "ws"
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := sense.New15(senseId, interrupt, done, []byte("1234"))
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
