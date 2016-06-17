package sense

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Sense interface {
	Id() string
	Connect(u url.URL, h http.Header) error
	Disconnect() error
	Send()
	Receive()
}

type Sense15 struct {
	conn      *websocket.Conn
	sleep     time.Duration
	interrupt chan os.Signal
	done      chan bool
	name      string
	logs      chan string
}

func (s *Sense15) Id() string {
	return s.name
}

func (s *Sense15) Connect(u url.URL, headers http.Header) error {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return err
	}
	s.conn = c
	return nil
}

func assemble(header, body, key []byte) []byte {
	hm := hmac.New(sha1.New, key)
	hm.Write(body)

	sig := hm.Sum(nil)
	signed := make([]byte, 0)
	envBuff := make([]byte, 0)

	buffer := bytes.NewBuffer(signed)
	buffer.Write(header)

	env := bytes.NewBuffer(envBuff)
	payload := buffer.Bytes()
	headerLen := uint64(len(payload))
	binary.Write(env, binary.LittleEndian, headerLen)

	env.Write(payload)

	bodyLen := uint64(len(body))
	binary.Write(env, binary.LittleEndian, bodyLen)
	env.Write(body)
	env.Write(sig)
	return env.Bytes()
}

func (s *Sense15) Send() {
	ticker := time.NewTicker(s.sleep)
	defer ticker.Stop()
	i := 0
	msgId := uint64(1)

	logs := make([]string, 0)
	for {
		select {
		case <-ticker.C:
			pb := &haneda.Preamble{}
			pb.Type = haneda.Preamble_BATCHED_PERIODIC_DATA.Enum()
			pb.Id = proto.Uint64(msgId)
			header, _ := proto.Marshal(pb)

			batched := &api.BatchedPeriodicData{}
			periodic := &api.PeriodicData{}
			periodic.Temperature = proto.Int32(27)

			batched.DeviceId = &s.name
			batched.FirmwareVersion = proto.Int32(888)
			batched.Data = append(batched.Data, periodic)

			body, pbErr := proto.Marshal(batched)
			if pbErr != nil {
				log.Println(pbErr)
				s.done <- true
			}

			env := assemble(header, body, []byte("abc"))
			err := s.conn.WriteMessage(websocket.BinaryMessage, env)
			if err != nil {
				log.Println("write:", err)
				s.done <- true
			}
			i++
			msgId++
		case logMessage := <-s.logs:
			if len(logs) == 200 {
				fmt.Println("Sending logMessage", logMessage)
				pb := &haneda.Preamble{}
				pb.Type = haneda.Preamble_SENSE_LOG.Enum()
				pb.Id = proto.Uint64(msgId)
				header, _ := proto.Marshal(pb)

				sLog := &api.SenseLog{}
				combined := strings.Join(logs, "\n")
				sLog.Text = &combined
				body, _ := proto.Marshal(sLog)

				env := assemble(header, body, []byte("abc"))
				err := s.conn.WriteMessage(websocket.BinaryMessage, env)
				if err != nil {
					log.Println("write:", err)
					s.done <- true
				}
				i++
				msgId++
				logs = make([]string, 0)
			} else {
				logs = append(logs, logMessage)
			}

		case <-s.interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			fmt.Println("Sent", s.conn.LocalAddr(), i)
			err := s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)

			}

			s.Disconnect()
			s.done <- true
		}
	}
}

func (s *Sense15) Disconnect() error {
	log.Println("Disconnecting")
	return s.conn.Close()
}

func (s *Sense15) Receive() {
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		m := &haneda.Preamble{}
		proto.Unmarshal(message, m)
		lm := fmt.Sprintf("%s: %d\n", s.name, m.GetId())
		s.logs <- lm
	}

	log.Println("Done receiving")
}

func New15(name string, sleep time.Duration, interrupt chan os.Signal, done chan bool) *Sense15 {
	return &Sense15{
		sleep:     sleep,
		interrupt: interrupt,
		name:      name,
		done:      done,
		logs:      make(chan string, 16),
	}
}
