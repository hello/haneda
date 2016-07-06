package sense

import (
	"errors"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"log"
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
	Interrupt chan os.Signal
	Done      chan bool
	Name      SenseId
	Logs      chan string
	PrivKey   []byte
	Auth      *SenseAuthHmacSha1
	Store     *Store
}

func (s *Sense15) Id() string {
	return string(s.Name)
}

func (s *Sense15) Connect(u *url.URL, headers http.Header) error {
	headers.Add("X-Hello-Sense-Id", string(s.Name))
	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return err
	}
	s.conn = c
	return nil
}

func (s *Sense15) Write(message []byte) error {
	return s.conn.WriteMessage(websocket.BinaryMessage, message)
}

func (s *Sense15) periodic(messageId uint64) *MessageParts {
	header := &haneda.Preamble{}
	header.Type = haneda.Preamble_BATCHED_PERIODIC_DATA.Enum()
	header.Id = proto.Uint64(messageId)

	batched := &api.BatchedPeriodicData{}
	periodic := &api.PeriodicData{}
	periodic.Temperature = proto.Int32(3500)

	n := string(s.Name)
	batched.DeviceId = &n
	batched.FirmwareVersion = proto.Int32(888)
	batched.Data = append(batched.Data, periodic)

	body, pbErr := proto.Marshal(batched)
	if pbErr != nil {
		log.Println("pbErr", pbErr)
		s.Done <- true
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
	sLog.DeviceId = proto.String(string(s.Name))
	body, _ := proto.Marshal(sLog)

	mp := &MessageParts{
		Header: pb,
		Body:   body,
	}
	return mp
}

type Result struct {
	body []byte
	err  error
}

func (s *Sense15) Read(timeout time.Duration) ([]byte, error) {

	results := make(chan *Result, 1)
	go func() {
		_, body, err := s.conn.ReadMessage()
		results <- &Result{body: body, err: err}
	}()

outer:
	for {
		select {
		case res := <-results:
			return res.body, res.err
		case <-time.After(timeout):
			fmt.Println("timeout 1")
			break outer
		}
	}

	return []byte{}, errors.New("timeout")
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
			env, err := s.Auth.Sign(mp)
			if err != nil {
				log.Println(err)
				s.Done <- true
			}

			duplicate := s.Store.Save(msgId)
			if duplicate != nil {
				log.Println("duplicate:", msgId)
				s.Done <- true
			}
			s.Write(env)
			log.Println("<--", mp.Header.GetType(), msgId)
			i++
			msgId++
		case logMessage := <-s.Logs:
			if len(logs) == 10 {
				mp := s.genLogs(msgId, logs)

				env, err := s.Auth.Sign(mp)
				if err != nil {
					log.Println(err)
					s.Done <- true
				}

				s.Write(env)
				duplicate := s.Store.Save(msgId)
				if duplicate != nil {
					log.Println("duplicate:", msgId)
					s.Done <- true
				}
				log.Println("<--", mp.Header.GetType(), msgId)
				i++
				msgId++
				logs = make([]string, 0)

			} else {
				logs = append(logs, logMessage)
			}

		case <-s.Interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			log.Println("Sent", s.conn.LocalAddr(), i)
			err := s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)

			}
			s.Disconnect()
			s.Done <- true
		}
	}
}

func (s *Sense15) Disconnect() error {
	signal.Notify(s.Interrupt, syscall.SIGTERM)
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

		log.Println("len:", len(message))
		mp, parseErr := s.Auth.Parse(message)
		if parseErr != nil {
			log.Println("parseErr", parseErr)
			continue
		}

		switch mp.Header.GetType() {
		case haneda.Preamble_ACK:
			ackMessage := &haneda.Ack{}
			err := proto.Unmarshal(mp.Body, ackMessage)
			if err != nil || ackMessage.Status.String() != haneda.Ack_SUCCESS.String() {
				log.Println("-->", err, ackMessage.GetMessageId(), ackMessage.GetStatus())
				continue
			}
			s.Store.Expire(ackMessage.GetMessageId())
			log.Println("-->", mp.Header.GetType(), ackMessage.GetMessageId())
		case haneda.Preamble_SYNC_RESPONSE:
			syncResp := &api.SyncResponse{}
			err := proto.Unmarshal(mp.Body, syncResp)
			if err != nil {
				log.Println(err)
				continue
			}
			// log.Println("-->", mp.Header.GetType(), syncResp.GetRingTimeAck())
			log.Println("-->", mp.Header.GetType(), syncResp.GetRoomConditions().Enum())
		default:
			log.Println("-->", mp.Header.GetType())

		}

		lm := fmt.Sprintf("%s: %d\n", s.Name, mp.Header.GetId())
		s.Logs <- lm
	}

	log.Println("Done receiving")
}

func NewDefaultSenseOneFive(sense *Sense15) *Sense15 {
	if len(sense.PrivKey) == 0 {
		sense.PrivKey = []byte{'d', 'e', 'f', 'a', 'u', 'l', 't', ' ', 'k', 'e', 'y'}
	}
	if sense.Logs == nil {
		sense.Logs = make(chan string, 16)
	}
	if sense.Store == nil {
		sense.Store = NewStore()
	}

	if sense.Auth == nil {
		sense.Auth = &SenseAuthHmacSha1{key: sense.PrivKey}
	}
	return sense
}

func New15(name SenseId, interrupt chan os.Signal, done chan bool, privKey []byte) *Sense15 {
	store := NewStore()
	return &Sense15{
		Interrupt: interrupt,
		Name:      name,
		Done:      done,
		Logs:      make(chan string, 16),
		PrivKey:   privKey,
		Auth:      &SenseAuthHmacSha1{key: privKey},
		Store:     store,
	}
}
