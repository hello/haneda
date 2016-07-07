package core

import (
	"bytes"
	"errors"
	"github.com/go-kit/kit/log"
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	LengthSigPlusIv = 48
)

var (
	ErrUnexpectedStatusCode  = errors.New("unexpected status code")
	ErrEndpointNotConfigured = errors.New("endpoint not configured")
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
	route, configured := f.routes[haneda.Preamble_BATCHED_PERIODIC_DATA]
	if !configured {
		return []byte{}, ErrEndpointNotConfigured
	}
	resp, err := f.PeriodicDataFwd.Do(content, privKey, route, http.StatusOK)
	return resp, err
}

func (f *HttpForwarder) Logs(message *api.SenseLog, privKey []byte) error {
	content, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	route, configured := f.routes[haneda.Preamble_SENSE_LOG]
	if !configured {
		return ErrEndpointNotConfigured
	}

	_, err = f.LogsFwd.Do(content, privKey, route, http.StatusNoContent)
	return err
}

type GenericForwarder struct {
	client   *http.Client
	endpoint string
	routes   map[haneda.PreamblePbType]string
	logger   log.Logger
}

func (f *GenericForwarder) Do(content, privKey []byte, path string, expectedHttpStatusCode int) ([]byte, error) {

	auth := &SuripuAuth{key: privKey}

	signed, err := auth.sign(content)
	if err != nil {
		f.logger.Log("error", err)
		return []byte{}, err
	}

	body := bytes.NewReader(signed)
	req, _ := http.NewRequest("POST", f.endpoint+path, body)
	// req = headers(req, s)
	resp, err := f.client.Do(req)

	if err != nil {
		f.logger.Log("error", err)
		return []byte{}, err
	}

	if resp.StatusCode != expectedHttpStatusCode {
		f.logger.Log("want", expectedHttpStatusCode, "got", resp.StatusCode, "endpoint", f.endpoint)
		return []byte{}, ErrUnexpectedStatusCode
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if len(data) > LengthSigPlusIv {
		return data[LengthSigPlusIv:], err
	}
	f.logger.Log("msg", "response_size_too_short", "response_size", len(data))
	return []byte{}, err
}

func NewDefaultHttpForwarder(endpoint string) *HttpForwarder {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return NewHttpForwarder(endpoint, client)
}

func NewHttpForwarder(endpoint string, client *http.Client) *HttpForwarder {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	genForwarder := &GenericForwarder{
		client:   client,
		endpoint: endpoint,
		logger:   logger,
	}

	routes := make(map[haneda.PreamblePbType]string)
	routes[haneda.Preamble_BATCHED_PERIODIC_DATA] = "/in/sense/batch"
	routes[haneda.Preamble_SENSE_LOG] = "/logs"

	fwd := &HttpForwarder{
		endpoint:           endpoint,
		client:             client,
		routes:             routes,
		LogsFwd:            genForwarder,
		PeriodicDataFwd:    genForwarder,
		MorpheusCommandFwd: genForwarder,
		SenseStateFwd:      genForwarder,
		FileManifestFwd:    genForwarder,
	}
	return fwd
}
