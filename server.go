package haneda

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"net/http"
	"strings"
)

func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func checkCreds(username, password string) bool {
	if password == "foo" {
		return true
	}
	return false
}

type authenticateFunc func(username, password string) bool

func extractBasicAuth(r *http.Request, f authenticateFunc) (string, error) {
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

func parse(content []byte) (*api.MessageParts, error) {
	bbuf := bytes.NewReader(content)
	var llen uint64
	err := binary.Read(bbuf, binary.LittleEndian, &llen)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	payload := make([]byte, llen)
	_, err = bbuf.Read(payload)
	sig := make([]byte, 20)
	_, err = bbuf.Read(sig)

	m := &api.MessageParts{
		Header: make([]byte, 0),
		Body:   payload,
		Sig:    sig,
	}

	return m, nil
}

func serialize(id uint64) []byte {
	pb := &haneda.Preamble{}
	pb.Id = &id
	m, _ := proto.Marshal(pb)
	return m
}

func echo(conn *websocket.Conn, senseId string, stats chan api.Stat) {
	i := 0
	defer conn.Close()
	for {

		_, content, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading json.", err)
			break
		}

		mp, err := parse(content)
		if err != nil {
			break
		}

		match := CheckMAC(mp.Body, mp.Sig, []byte("abc"))
		if !match {
			fmt.Println("don't match!!!")
			fmt.Printf("%v\n", mp.Body)
			fmt.Println("len(content)", len(content))
			fmt.Printf("%x\n", mp.Sig)
		}

		pb := &haneda.Preamble{}
		protoErr := proto.Unmarshal(mp.Body, pb)
		if protoErr != nil {
			fmt.Println(protoErr)
			break
		}

		// fmt.Printf("Got message: %s %d\n", pb.GetType().String(), pb.GetId())
		newId := uint64(88888)
		m := serialize(newId)
		if err = conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
			fmt.Println(err)
			break
		}
		i++
	}
	s := api.Stat{
		Source: conn.RemoteAddr().String(),
		Count:  i,
	}
	fmt.Println(s)
	stats <- s
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/health")
	fmt.Fprintf(w, "%s\n", "ok")
}

func DisplayStats(s chan api.Stat) {
	running := 0
	for m := range s {
		running += m.Count
		fmt.Printf("stats: %s = %d\n", m.Source, m.Count)
		fmt.Printf("\t total so far: %d\n", running)
	}
}
