package core

import (
	"bytes"
	"errors"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"github.com/hello/haneda/sense"
	"log"
	"net/http"
)

type Bridge struct {
	client   *http.Client
	endpoint string
	routes   map[haneda.PreamblePbType]string
}

func NewBridge(endpoint string) *Bridge {
	client := http.DefaultClient
	// configure with timeouts

	routes := make(map[haneda.PreamblePbType]string)
	routes[haneda.Preamble_BATCHED_PERIODIC_DATA] = "/in/sense/batch"
	routes[haneda.Preamble_SENSE_LOG] = "/logs"
	return &Bridge{
		endpoint: endpoint,
		client:   client,
		routes:   routes,
	}
}

func (b *Bridge) Route(mp *sense.MessageParts) error {
	route, found := b.routes[mp.Header.GetType()]
	if !found {
		log.Println("not found", mp.Header.GetType())
	}

	fmt.Println(mp.Header.GetType().String(), route)
}

func headers(req *http.Request, s *SenseConn) *http.Request {
	req.Header.Add("X-Hello-Sense-Id", s.SenseId)
	req.Header.Add("X-Hello-Sense-MFW", s.MiddleFirmwareVersion)
	req.Header.Add("X-Hello-Sense-TFW", s.TopFirmwareVersion)
	req.Header.Add("Content-Type", "application/x-protobuf")
	return req
}

func (b *Bridge) Logs(message *api.SenseLog, s *SenseConn) error {
	log.Println("sending log message to:", b.endpoint)
	content, _ := proto.Marshal(message)

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

func (b *Bridge) PeriodicData(message *api.BatchedPeriodicData, s *SenseConn) error {
	log.Println("sending periodic data to:", b.endpoint)
	content, _ := proto.Marshal(message)

	signed, err := sign(content, s.PrivKey)
	if err != nil {
		return err
	}

	body := bytes.NewReader(signed)
	req, _ := http.NewRequest("POST", b.endpoint+"/in/sense/batch", body)
	req = headers(req, s)
	resp, err := b.client.Do(req)

	if err != nil {
		log.Println(err)
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Got non 200 from Suripu service")
	}
	return nil
}
