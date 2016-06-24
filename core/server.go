package core

import (
	"encoding/base64"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

func checkCreds(username, password string) bool {
	if password == "foo" {
		return true
	}
	return false
}

type SenseConn struct {
	TopFirmwareVersion    string
	MiddleFirmwareVersion string
	SenseId               string
	Conn                  *websocket.Conn
	PrivKey               []byte
}

type AuthenticateFunc func(username, password string) bool

func extractBasicAuth(r *http.Request, f AuthenticateFunc) (string, error) {
	if len(r.Header["Authorization"]) > 0 {

		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			return "", errors.New("Bad headers")
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !f(pair[0], pair[1]) {

			return "", errors.New("Bad headers")
		}
		return pair[0], nil
	}
	return "", errors.New("not found")
}

func basicAuth(auth string) (string, string) {
	payload, _ := base64.StdEncoding.DecodeString(auth)
	pair := strings.SplitN(string(payload), ":", 2)
	return pair[0], pair[1]
}
