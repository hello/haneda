package core

import (
	"github.com/hello/haneda/api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

type mockTransport struct {
	statusCode int
}

func newMockTransport(statusCode int) http.RoundTripper {
	return &mockTransport{statusCode: statusCode}
}

// Implement http.RoundTripper
func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Create mocked http.Response
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: t.statusCode,
	}
	response.Header.Set("Content-Type", "application/json")

	responseBody :=
		`{
    "Accept-Encoding": [
        "mock"
    ],
    "User-Agent": [
        "mock"
    ],
    "X-Ip-Country": [
        "Japan(Mock)"
    ],
    "X-Real-Ip": [
        "192.168.1.1"
    ]
}`
	response.Body = ioutil.NopCloser(strings.NewReader(responseBody))
	return response, nil
}

var (
	client = http.DefaultClient
)

func TestBridgeLogs(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	client.Transport = newMockTransport(200)
	bridge := NewHttpBridgeWithClient("http://example.com/", client)
	conn := &SenseConn{PrivKey: []byte("1234567891234567")}

	logMessage := &api.SenseLog{}
	err := bridge.Logs(logMessage, conn)
	if err == nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	client.Transport = newMockTransport(200)

	err = bridge.Logs(logMessage, conn)
	if err == nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
}

func TestBridgeLogsBadKey(t *testing.T) {
	client.Transport = newMockTransport(200)
	bridge := NewHttpBridgeWithClient("http://example.com/", client)
	conn := &SenseConn{PrivKey: []byte("bad key")}

	logMessage := &api.SenseLog{}
	err := bridge.Logs(logMessage, conn)
	if err == nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
}
