package main

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/core"
	"github.com/hello/haneda/haneda"
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
	key     = []byte("abcdefghjklmnop")
	ks      = &sense.FakeKeyStore{SenseId: string(senseId), Key: key}
	conf    = &core.HelloConfig{
		Redis: &core.RedisConfig{
			PubSub: "example",
		},
	}
)

type FakePeriodicDataBrigde struct {
	core.NoopBridge
}

func (b *FakePeriodicDataBrigde) PeriodicData(message *api.BatchedPeriodicData, key []byte) ([]byte, error) {
	sr, err := sense.GenSyncResp(uint64(1))
	return sr.Body, err
}

func TestConnectToWebSocketServer(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(&FakePeriodicDataBrigde{}, nil, done, messages, ks, conf)

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
	simple.Shutdown()
}

func TestWrite(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(&FakePeriodicDataBrigde{}, nil, done, messages, ks, conf)

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

type genFunc func(msgId uint64) (*sense.MessageParts, error)

func TestPingPong(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	done := make(chan bool, 0)
	messages := make(chan *sense.MessageParts, 0)

	simple := core.NewSimpleHelloServer(&FakePeriodicDataBrigde{}, nil, done, messages, ks, conf)

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

	generators := make(map[haneda.PreamblePbType]genFunc)

	generators[haneda.Preamble_BATCHED_PERIODIC_DATA] = sense.GenPeriodic
	generators[haneda.Preamble_SENSE_LOG] = sense.GenLogs

	tests := []struct {
		messageId        uint64
		headerType       haneda.PreamblePbType
		returnHeaderType *haneda.PreamblePbType
		timeout          time.Duration
	}{
		{
			messageId:        uint64(1),
			headerType:       haneda.Preamble_BATCHED_PERIODIC_DATA,
			returnHeaderType: haneda.Preamble_SYNC_RESPONSE.Enum(),
			timeout:          time.Second * 1,
		},
		{
			messageId:        uint64(1),
			headerType:       haneda.Preamble_SENSE_LOG,
			returnHeaderType: nil,
			timeout:          time.Second * 10,
		},
	}

	auth := sense.NewSenseAuthHmacSha1WithContext(key, senseId, "test")

	for _, test := range tests {

		gen, found := generators[test.headerType]
		if !found {

		}
		mp, _ := gen(test.messageId)

		signed, _ := auth.Sign(mp)
		err := client.Write(signed)

		if err != nil {
			t.Fatalf("%v", err)
		}

		msg, err := client.Read(test.timeout)
		if err != nil {
			t.Errorf("Error Reading: %v", err)
			t.FailNow()
		}

		checkAck("one", t, auth, test.messageId, msg)

		if test.returnHeaderType != nil {
			msg, err = client.Read(test.timeout)
			if err != nil {
				t.Errorf("Error Reading: %v", err)
				t.FailNow()
			}
			check("two", t, auth, test.messageId, msg, test.returnHeaderType)
		}
	}
	simple.Shutdown()
}

func checkAck(name string, t *testing.T, parser sense.MessageParser, id uint64, msg []byte) {
	rmp, err := parser.Parse(msg)
	if err != nil {
		t.Errorf("%s, Error Parsing: %v", name, err)
		t.FailNow()
	}
	ack := &haneda.Ack{}
	proto.Unmarshal(rmp.Body, ack)

	if ack.GetMessageId() != id {
		t.Errorf("name=%s got=%d want=%d", name, rmp.Header.GetId(), id)
		t.FailNow()
	}
}

func check(name string, t *testing.T, parser sense.MessageParser, id uint64, msg []byte, want *haneda.PreamblePbType) {
	rmp, err := parser.Parse(msg)
	if err != nil {
		t.Errorf("%s, Error Parsing: %v", name, err)
		t.FailNow()
	}

	if rmp.Header.GetType() != *want {
		t.Errorf("name=%s got=%s want=%s", name, rmp.Header.GetType().Enum(), want.Enum())
		t.FailNow()
	}

}
