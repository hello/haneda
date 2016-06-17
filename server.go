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

type SenseConn struct {
	TopFirmwareVersion    string
	MiddleFirmwareVersion string
	SenseId               string
	Conn                  *websocket.Conn
	PrivKey               []byte
}

type ApiMessage struct {
	Type MessageType
	Body []byte
}

type MessageType int

const (
	L MessageType = iota
	P
	R
)

type LogMessageWrapper struct {
	Id      uint64
	Content *api.SenseLog
}

type PeriodicMessage struct {
	Id      int32  `json:"id"`
	Temp    int32  `json:"temp"`
	SenseId string `json:"sense_id"`
}

type RegisterMessage struct {
	Id      int32  `json:"id"`
	SenseId string `json:"sense_id"`
}

type AckMessageWrapper struct {
	Id      uint64
	SenseId string
	// TODO: add protobuf here
}

type HelloServer struct {
	Logs          chan *LogMessageWrapper
	Periodic      chan *PeriodicMessage
	Registrations chan *RegisterMessage
	Ventilator    *Ventilator
}

func NewHelloServer(v *Ventilator) *HelloServer {
	return &HelloServer{
		Logs:          make(chan *LogMessageWrapper, 2),
		Periodic:      make(chan *PeriodicMessage, 2),
		Registrations: make(chan *RegisterMessage, 2),
		Ventilator:    v,
	}
}

func (h *HelloServer) parse(messageParts *api.MessageParts) {
	/*

		switch preamble.GetType() {
		case haneda.Preamble_BATCHED_PERIODIC_DATA:
			periodic := &api.BatchedPeriodicData{}
			pbErr := proto.Unmarshal(body, periodic)
			if pbErr != nil {
				fmt.Println("here", pbErr)
				return nil, pbErr
			}
		case haneda.Preamble_SENSE_LOG:
			slog := &api.SenseLog{}
			pbErr := proto.Unmarshal(body, slog)
			if pbErr != nil {
				fmt.Println("log", pbErr)
				return nil, pbErr
			}
			fmt.Println("len logs", len(slog.GetText()))
		default:
			fmt.Println("Don't know what that is")
			return nil, errors.New(fmt.Sprintf("Unknown pb type: %s", preamble.GetType().String()))
		}

	*/
	switch messageParts.Header.GetType() {
	case haneda.Preamble_SENSE_LOG:
		m := &api.SenseLog{}
		proto.Unmarshal(messageParts.Body, m)
		wrapper := &LogMessageWrapper{
			Id:      messageParts.Header.GetId(),
			Content: m,
		}

		h.Logs <- wrapper
	default:
		panic("should not happen")
	}
}

func (h *HelloServer) Run() {
	for {
		select {
		case m := <-h.Logs:
			fmt.Println("Saving log:", m.Content.GetText())
			h.Ventilator.push(&AckMessageWrapper{Id: m.Id, SenseId: m.Content.GetDeviceId()})
		case m, open := <-h.Ventilator.raw:
			if !open {
				fmt.Println("Channel was closed")
				return
			}
			h.parse(m)
		}
	}
}

func extractBasicAuth(r *http.Request, f api.AuthenticateFunc) (string, error) {
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
	var headerLen uint64
	err := binary.Read(bbuf, binary.LittleEndian, &headerLen)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	header := make([]byte, headerLen)
	n, err := bbuf.Read(header)

	if uint64(n) != headerLen {
		return nil, errors.New("header Len don't match")
	}
	var bodyLen uint64
	err = binary.Read(bbuf, binary.LittleEndian, &bodyLen)
	if err != nil {
		return nil, err
	}

	body := make([]byte, bodyLen)
	n, err = bbuf.Read(body)

	if uint64(n) != bodyLen {
		return nil, errors.New("body Len don't match")
	}

	sig := make([]byte, 20)
	n, err = bbuf.Read(sig)

	preamble := &haneda.Preamble{}
	protoErr := proto.Unmarshal(header, preamble)

	if protoErr != nil {
		return nil, protoErr
	}

	m := &api.MessageParts{
		Header: preamble,
		Body:   body,
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

func spin(s *SenseConn) {
	i := 0
	defer s.Conn.Close()
	for {

		_, content, err := s.Conn.ReadMessage()
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

		// fmt.Printf("Got message: %s %d\n", pb.GetType().String(), pb.GetId())
		msgId := mp.Header.GetId()
		if msgId%uint64(100) == 0 {
			fmt.Println("msgid:", msgId)
		}

		newId := uint64(mp.Header.GetId())
		m := serialize(newId)
		if err = s.Conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
			fmt.Println(err)
			break
		}
		i++
	}
	fmt.Println("Processed:", i)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/health")
	fmt.Fprintf(w, "%s\n", "ok")
}

func DisplayStats(s chan api.Stat) {
	running := 0
	for m := range s {
		running += m.Count
		fmt.Printf("%s", m)
		fmt.Printf("\t total so far: %d\n", running)
	}
}
