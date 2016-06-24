package core

import (
	"bytes"
	"errors"
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/api"
	"log"
	"net/http"
)

type Bridge struct {
	client   *http.Client
	endpoint string
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

func (b *Bridge) PeriodicData(message *api.BatchedPeriodicData) error {
	// make http request
	// fmt.Println("make http request to suripu")
	return nil
}

func NewBridge(endpoint string) *Bridge {
	client := http.DefaultClient
	// configure with timeouts
	return &Bridge{
		endpoint: endpoint,
		client:   client,
	}
}
