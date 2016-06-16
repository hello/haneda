package sense

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/haneda"
	"log"
	"net/http"
	"net/url"
	"os"
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

func (s *Sense15) Send() {
	ticker := time.NewTicker(s.sleep)
	defer ticker.Stop()
	i := 0
	for {
		select {
		case <-ticker.C:

			pb := &haneda.Preamble{}
			pb.Type = haneda.Preamble_BATCHED_PERIODIC_DATA.Enum()
			id := uint64(99999999)
			pb.Id = &id
			buff, _ := proto.Marshal(pb)

			hm := hmac.New(sha1.New, []byte("abc"))
			hm.Write(buff)

			sig := hm.Sum(nil)
			signed := make([]byte, 0)
			envBuff := make([]byte, 0)

			buffer := bytes.NewBuffer(signed)
			buffer.Write(buff)

			env := bytes.NewBuffer(envBuff)
			payload := buffer.Bytes()
			llen := uint64(len(payload))
			binary.Write(env, binary.LittleEndian, llen)

			env.Write(payload)
			env.Write(sig)

			err := s.conn.WriteMessage(websocket.BinaryMessage, env.Bytes())
			if err != nil {
				log.Println("write:", err)
				s.done <- true
			}
			i++
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
		log.Printf("%s: %d\n", s.name, m.GetId())
	}

	log.Println("Done receiving")
}

func New15(name string, sleep time.Duration, interrupt chan os.Signal, done chan bool) *Sense15 {
	return &Sense15{
		sleep:     sleep,
		interrupt: interrupt,
		name:      name,
		done:      done,
	}
}
