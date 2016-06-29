package sense

import (
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Client interface {
	Id() string
	Connect(u *url.URL, h http.Header) error
	Disconnect() error
	Send(t time.Duration)
	Receive()
	Write(message []byte)
}

type Sense15 struct {
	conn      *websocket.Conn
	interrupt chan os.Signal
	done      chan bool
	name      string
	logs      chan string
	privKey   []byte
	auth      *SenseAuth
	store     *Store
}

func (s *Sense15) Id() string {
	return s.name
}

func (s *Sense15) Connect(u *url.URL, headers http.Header) error {
	headers.Add("X-Hello-Sense-Id", s.name)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return err
	}
	s.conn = c
	return nil
}

func (s *Sense15) Write(message []byte) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int31n(1000)
	time.Sleep(time.Duration(n) * time.Millisecond)
	err := s.conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		log.Println("write:", err)
		s.done <- true
	}
}

func (s *Sense15) periodic(messageId uint64) *MessageParts {
	header := &haneda.Preamble{}
	header.Type = haneda.Preamble_BATCHED_PERIODIC_DATA.Enum()
	header.Id = proto.Uint64(messageId)

	batched := &api.BatchedPeriodicData{}
	periodic := &api.PeriodicData{}
	periodic.Temperature = proto.Int32(27)

	batched.DeviceId = &s.name
	batched.FirmwareVersion = proto.Int32(888)
	batched.Data = append(batched.Data, periodic)

	body, pbErr := proto.Marshal(batched)
	if pbErr != nil {
		log.Println("pbErr", pbErr)
		s.done <- true
	}

	mp := &MessageParts{
		Header: header,
		Body:   body,
	}
	return mp
}

func (s *Sense15) genLogs(messageId uint64, logs []string) *MessageParts {
	pb := &haneda.Preamble{}
	pb.Type = haneda.Preamble_SENSE_LOG.Enum()
	pb.Id = proto.Uint64(messageId)

	sLog := &api.SenseLog{}
	combined := strings.Join(logs, "\n")
	sLog.Text = &combined
	sLog.DeviceId = proto.String(s.name)
	body, _ := proto.Marshal(sLog)

	mp := &MessageParts{
		Header: pb,
		Body:   body,
	}
	return mp
}

func (s *Sense15) Send(sleep time.Duration) {
	ticker := time.NewTicker(sleep)
	defer ticker.Stop()
	i := 0
	msgId := uint64(1)

	logs := make([]string, 0)
	for {
		select {
		case <-ticker.C:

			mp := s.periodic(msgId)
			env, err := s.auth.Sign(mp)
			if err != nil {
				log.Println(err)
				s.done <- true
			}

			duplicate := s.store.Save(msgId)
			if duplicate != nil {
				log.Println("duplicate:", msgId)
				s.done <- true
			}
			s.Write(env)
			log.Println("<--", mp.Header.GetType(), msgId)
			i++
			msgId++
		case logMessage := <-s.logs:
			if len(logs) == 10 {
				mp := s.genLogs(msgId, logs)

				env, err := s.auth.Sign(mp)
				if err != nil {
					log.Println(err)
					s.done <- true
				}

				s.Write(env)
				duplicate := s.store.Save(msgId)
				if duplicate != nil {
					log.Println("duplicate:", msgId)
					s.done <- true
				}
				log.Println("<--", mp.Header.GetType(), msgId)
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
			log.Println("Sent", s.conn.LocalAddr(), i)
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
	signal.Notify(s.interrupt, syscall.SIGTERM)
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

		mp, parseErr := s.auth.Parse(message)
		if parseErr != nil {
			log.Println("parseErr", parseErr)
			continue
		}

		switch mp.Header.GetType() {
		case haneda.Preamble_ACK:
			ackMessage := &haneda.Ack{}
			err := proto.Unmarshal(mp.Body, ackMessage)
			if err != nil || ackMessage.Status.String() != haneda.Ack_SUCCESS.String() {
				log.Println(err, ackMessage.GetMessageId(), ackMessage.GetStatus())
				continue
			}
			s.store.Expire(ackMessage.GetMessageId())
			log.Println("-->", mp.Header.GetType(), ackMessage.GetMessageId())
		case haneda.Preamble_SYNC_RESPONSE:
			syncResp := &api.SyncResponse{}
			err := proto.Unmarshal(mp.Body, syncResp)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("-->", mp.Header.GetType(), syncResp.GetRingTimeAck())
		default:
			log.Println("-->", mp.Header.GetType())

		}

		lm := fmt.Sprintf("%s: %d\n", s.name, mp.Header.GetId())
		s.logs <- lm
	}

	log.Println("Done receiving")
}

func New15(name string, interrupt chan os.Signal, done chan bool, privKey []byte) *Sense15 {
	store := NewStore()
	return &Sense15{
		interrupt: interrupt,
		name:      name,
		done:      done,
		logs:      make(chan string, 16),
		privKey:   privKey,
		auth:      &SenseAuth{key: privKey},
		store:     store,
	}
}
