package core

import (
	"bytes"
	"errors"
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"github.com/hello/haneda/sense"
	"io/ioutil"
	"log"
	"net/http"
)

type Bridge interface {
	PeriodicData(message *api.BatchedPeriodicData, s *SenseConn) ([]byte, error)
	Logs(message *api.SenseLog, s *SenseConn) error
}

type HttpBridge struct {
	client   *http.Client
	endpoint string
	routes   map[haneda.PreamblePbType]string
}

func NewHttpBridge(endpoint string) Bridge {
	client := http.DefaultClient
	// configure with timeouts

	routes := make(map[haneda.PreamblePbType]string)
	routes[haneda.Preamble_BATCHED_PERIODIC_DATA] = "/in/sense/batch"
	routes[haneda.Preamble_SENSE_LOG] = "/logs"
	return HttpBridge{
		endpoint: endpoint,
		client:   client,
		routes:   routes,
	}
}

func NewHttpBridgeWithClient(endpoint string, client *http.Client) Bridge {
	return HttpBridge{
		endpoint: endpoint,
		client:   client,
		routes:   make(map[haneda.PreamblePbType]string),
	}
}

func (b *HttpBridge) Route(mp *sense.MessageParts) error {
	route, found := b.routes[mp.Header.GetType()]
	if !found {
		log.Println("not found", mp.Header.GetType())
	}

	log.Println(mp.Header.GetType().String(), route)
	return nil
}

func headers(req *http.Request, s *SenseConn) *http.Request {
	req.Header.Add("X-Hello-Sense-Id", s.SenseId)
	req.Header.Add("X-Hello-Sense-MFW", s.MiddleFirmwareVersion)
	req.Header.Add("X-Hello-Sense-TFW", s.TopFirmwareVersion)
	req.Header.Add("Content-Type", "application/x-protobuf")
	return req
}

func (b HttpBridge) Logs(message *api.SenseLog, s *SenseConn) error {
	log.Println("sending log message to:", b.endpoint)
	content, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	signed, err := sign(content, s.PrivKey)
	if err != nil {
		return err
	}

	body := bytes.NewReader(signed)
	req, _ := http.NewRequest("POST", b.endpoint+"/logs", body)
	req = headers(req, s)
	resp, err := b.client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}

	if resp.StatusCode != 204 {
		return errors.New("Got non 204 from Suripu service")
	}
	return nil

}

func (b HttpBridge) PeriodicData(message *api.BatchedPeriodicData, s *SenseConn) ([]byte, error) {
	log.Println("sending periodic data to:", b.endpoint)
	content, _ := proto.Marshal(message)

	empty := make([]byte, 0)

	signed, err := sign(content, s.PrivKey)
	if err != nil {
		return empty, err
	}

	body := bytes.NewReader(signed)
	req, _ := http.NewRequest("POST", b.endpoint+"/in/sense/batch", body)
	req = headers(req, s)
	resp, err := b.client.Do(req)

	if err != nil {
		log.Println(err)
		return empty, err
	}

	if resp.StatusCode != 200 {
		return empty, errors.New("Got non 200 from Suripu service")
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	log.Println("len(data):", len(data))
	return data[48:], err
}
