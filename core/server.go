package core

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

var (
	ErrBadHeader = errors.New("Bad Headers")
)

func checkCreds(username, password string) bool {
	if password == "foo" {
		return true
	}
	return false
}

func extractBasicAuth(r *http.Request, f AuthenticateFunc) (string, error) {
	if len(r.Header["Authorization"]) > 0 {

		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			return "", ErrBadHeader
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !f(pair[0], pair[1]) {

			return "", ErrBadHeader
		}
		return pair[0], nil
	}
	return "", ErrBadHeader
}

func basicAuth(auth string) (string, string) {
	payload, _ := base64.StdEncoding.DecodeString(auth)
	pair := strings.SplitN(string(payload), ":", 2)
	return pair[0], pair[1]
}
