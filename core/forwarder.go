package core

import (
	"bytes"
	"errors"
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	LengthSigPlusIv = 48
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

type Forwarder interface {
	Do(message, privKey []byte, path string, expectedHttpStatusCode int) ([]byte, error)
}

type Bridge interface {
	PeriodicData(message *api.BatchedPeriodicData, privKey []byte) ([]byte, error)
	Logs(message *api.SenseLog, privKey []byte) error
}

type HttpForwarder struct {
	PeriodicDataFwd    Forwarder
	LogsFwd            Forwarder
	SenseStateFwd      Forwarder
	FileManifestFwd    Forwarder
	MorpheusCommandFwd Forwarder
	client             *http.Client
	endpoint           string
	routes             map[haneda.PreamblePbType]string
}

func (f *HttpForwarder) PeriodicData(message *api.BatchedPeriodicData, privKey []byte) ([]byte, error) {
	content, err := proto.Marshal(message)
	if err != nil {
		return []byte{}, err
	}
	resp, err := f.PeriodicDataFwd.Do(content, privKey, "/in/sense/batch", http.StatusOK)
	return resp, err
}

func (f *HttpForwarder) Logs(message *api.SenseLog, privKey []byte) error {
	content, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	_, err = f.LogsFwd.Do(content, privKey, "/logs", http.StatusNoContent)
	return err
}

type GenericForwarder struct {
	client   *http.Client
	endpoint string
	routes   map[haneda.PreamblePbType]string
}

func (f *GenericForwarder) Do(content, privKey []byte, path string, expectedHttpStatusCode int) ([]byte, error) {
	signed, err := sign(content, privKey)
	if err != nil {
		return []byte{}, err
	}

	body := bytes.NewReader(signed)
	req, _ := http.NewRequest("POST", f.endpoint+path, body)
	// req = headers(req, s)
	resp, err := f.client.Do(req)

	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	if resp.StatusCode != expectedHttpStatusCode {
		log.Printf("want=%d got=%d endpoint=%", expectedHttpStatusCode, resp.StatusCode, f.endpoint)
		return []byte{}, ErrUnexpectedStatusCode
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if len(data) > LengthSigPlusIv {
		return data[48:], err
	}
	return []byte{}, err
}

func NewHttpForwarder(endpoint string, client *http.Client) *HttpForwarder {

	genForwarder := &GenericForwarder{
		client:   client,
		endpoint: endpoint,
	}

	fwd := &HttpForwarder{
		endpoint:           endpoint,
		client:             client,
		routes:             make(map[haneda.PreamblePbType]string),
		LogsFwd:            genForwarder,
		PeriodicDataFwd:    genForwarder,
		MorpheusCommandFwd: genForwarder,
		SenseStateFwd:      genForwarder,
		FileManifestFwd:    genForwarder,
	}
	return fwd
}
