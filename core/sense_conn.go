package core

import (
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/haneda"
	"github.com/hello/haneda/sense"
	"log"
	"time"
)

type SenseConn struct {
	TopFirmwareVersion    sense.TopFirmwareVersion
	MiddleFirmwareVersion sense.MiddleFirmwareVersion
	SenseId               sense.SenseId
	Conn                  *websocket.Conn
	PrivKey               []byte
	out                   chan *sense.MessageParts
	internalMsgs          chan []byte
	signer                *sense.MessageSigner
	parser                *sense.MessageParser
	bridge                Bridge
	remover               ConnectionRemover
}

func (c *SenseConn) Send(message *sense.MessageParts) {
	c.out <- message
}

func (c *SenseConn) write() {
	fmt.Println("starting write thread")
outer:
	for {
		select {
		case m := <-c.out:
			log.Printf("Sending %s Message to: %s\n", m.Header.GetType().String(), m.SenseId)
			if err := c.Conn.WriteMessage(websocket.BinaryMessage, m.Body); err != nil {
				fmt.Println(err)
				break outer
			}
		case m := <-c.internalMsgs: // assuming already fully assembled messages
			if err := c.Conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
				fmt.Println(err)
				break outer
			}
		}
	}
	log.Println("Writing thread stopped for sense", c.SenseId)
}

func (c *SenseConn) Serve() {
	auth := sense.NewAuth(c.PrivKey, c.SenseId)

	i := 0
	defer c.Conn.Close()

	go c.write()
	for {
		_, content, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading.", err)
			break
		}

		mp, err := auth.Parse(content)
		if err != nil {
			log.Println(c.SenseId, err)
			break
		}

		// prepare ack message
		ack := &haneda.Ack{}
		ack.MessageId = proto.Uint64(mp.Header.GetId())
		ack.Status = haneda.Ack_SUCCESS.Enum()

		outbox := make([]*sense.MessageParts, 0)

		resp, err := dispatch(c.bridge, mp, c)
		if err != nil {
			// Override status since bridge responded with error
			ack.Status = haneda.Ack_CLIENT_REQUEST_ERROR.Enum()
			fmt.Println(mp.Header.GetId(), "Ack_CLIENT_REQUEST_ERROR")
			fmt.Println(err)
		}

		body, _ := proto.Marshal(ack)

		header := &haneda.Preamble{}
		header.Type = haneda.Preamble_ACK.Enum()

		out := &sense.MessageParts{
			Header: header,
			Body:   body,
		}

		outbox = append(outbox, out)

		// response from server might be empty
		if len(resp) > 0 {
			switch mp.Header.GetType() {
			case haneda.Preamble_BATCHED_PERIODIC_DATA:
				syncHeader := &haneda.Preamble{}
				syncHeader.Type = haneda.Preamble_SYNC_RESPONSE.Enum()
				syncHeader.Id = proto.Uint64(uint64(time.Now().UnixNano()))

				out2 := &sense.MessageParts{
					Header: syncHeader,
					Body:   resp,
				}
				outbox = append(outbox, out2)
			default:
				log.Println("No response needed")
			}
		}

		for _, mp := range outbox {
			serialized, err := auth.Sign(mp)
			if err != nil {
				log.Println(err)
				break
			}
			c.internalMsgs <- serialized
		}
		i++
	}
	c.remover.Remove(c.SenseId)
	log.Println(c.SenseId, "Processed:", i)
}

func (c *SenseConn) Listen(in <-chan *sense.MessageParts) {
	for msg := range in {
		c.out <- msg
	}
}
