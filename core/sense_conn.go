package core

import (
	"github.com/go-kit/kit/log"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hello/haneda/haneda"
	"github.com/hello/haneda/sense"
	"time"
)

// SenseConn represents a connection from Sense to our server
type SenseConn struct {
	TopFirmwareVersion    sense.TopFirmwareVersion    // Derived from HTTP headers
	MiddleFirmwareVersion sense.MiddleFirmwareVersion // Derived from HTTP headers
	SenseId               sense.SenseId               // Derived from HTTP headers
	Conn                  *websocket.Conn             // WS Conn from successful upgrade
	PrivKey               []byte                      // Key is fetched only once when connection is established
	out                   chan *sense.MessageParts    // add to this channel to send to sense. body should be signed
	internalMsgs          chan []byte                 // add to this channel to write raw messages to sense
	signer                sense.MessageSigner         // signer signs outgoing messages
	parser                sense.MessageParser         // parser parses incoming messages
	bridge                Bridge                      // bridge dispatches proto to http service
	remover               ConnectionRemover           // removes ws connection when sense disconnects
	logger                log.Logger
}

func (c *SenseConn) Send(message *sense.MessageParts) {
	c.out <- message
}

func (c *SenseConn) write() {
	c.logger.Log("action", "starting write thread")
outer:
	for {
		select {
		case m := <-c.out:
			c.logger.Log("action", "sending-message", "msg_type", m.Header.GetType().String())
			if err := c.Conn.WriteMessage(websocket.BinaryMessage, m.Body); err != nil {
				c.logger.Log("error", err)
				break outer
			}
		case m := <-c.internalMsgs: // assuming already fully assembled messages
			if err := c.Conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
				c.logger.Log("error", err)
				break outer
			}
		}
	}
	c.logger.Log("action", "write thread stopped")
}

func (c *SenseConn) Serve(stats chan *HelloStat) {
	i := 0
	defer c.Conn.Close()
	go c.write()
	for {
		// this is blocking
		_, content, err := c.Conn.ReadMessage()
		if err != nil {
			stats <- &HelloStat{ErrRead: hUint64(1)}
			break
		}
		stats <- &HelloStat{OkRead: hUint64(1)}

		mp, err := c.parser.Parse(content)
		if err != nil {
			c.logger.Log("parse_error", err)
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
			c.logger.Log("action=send-error-ack", "message_id", mp.Header.GetId())
			c.logger.Log("error", err)
			stats <- &HelloStat{ErrProxy: hUint64(1)}
		}

		body, _ := proto.Marshal(ack)

		header := &haneda.Preamble{}
		header.Type = haneda.Preamble_ACK.Enum()

		out := &sense.MessageParts{
			Header: header,
			Body:   body,
		}
		c.logger.Log("sending_ack", ack.GetMessageId())
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
			case haneda.Preamble_MORPHEUS_COMMAND:
				c.logger.Log("cmd", "command buddy")
				cmdHeader := &haneda.Preamble{}
				cmdHeader.Type = haneda.Preamble_MORPHEUS_COMMAND.Enum()
				cmdHeader.Id = proto.Uint64(uint64(time.Now().UnixNano()))

				m := &sense.MessageParts{
					Header: header,
					Body:   resp,
				}
				outbox = append(outbox, m)
			default:
				c.logger.Log("msg", "no response needed")
			}
		}

		for _, mp := range outbox {
			serialized, err := c.signer.Sign(mp)
			if err != nil {
				c.logger.Log("error", err)
				break
			}
			c.internalMsgs <- serialized
		}
		i++
	}
	c.remover.Remove(c.SenseId)
	c.logger.Log("processed", i)
}

func (c *SenseConn) Listen(in <-chan *sense.MessageParts) {
	for msg := range in {
		c.out <- msg
	}
}
