// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	// "github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
)

// var addr = flag.String("addr", "ws-dev.hello.is:443", "http service address")
var addr = flag.String("addr", "localhost:8082", "http service address")

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type HelloWsClient struct {
	interrupt chan os.Signal
	name      string
}

func NewHelloWsClient(name string, interrupt chan os.Signal) *HelloWsClient {
	return &HelloWsClient{
		name:      name,
		interrupt: interrupt,
	}
}

func (h *HelloWsClient) Start(addr string, sleep int64, done chan struct{}) {

	signal.Notify(h.interrupt, os.Interrupt)

	headers := http.Header{}
	headers.Add("Authorization", "Basic "+basicAuth(h.name, "foo"))
	// u := url.URL{Scheme: "wss", Host: *addr, Path: "/echo"}
	u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer c.Close()

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			m := &haneda.Preamble{}
			proto.Unmarshal(message, m)
			log.Printf("%s: %d\n", h.name, m.GetId())

		}
	}()

	ticker := time.NewTicker(time.Duration(sleep) * time.Millisecond)
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

			err := c.WriteMessage(websocket.BinaryMessage, env.Bytes())
			if err != nil {
				log.Println("write:", err)
				// return
			}
			i++
		case <-h.interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			fmt.Println("Sent", c.LocalAddr(), i)
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
				fmt.Println("Sent", c.LocalAddr(), i)
			case <-time.After(time.Second):

			}

			c.Close()
			return
		}
	}
}
func main() {
	flag.Parse()
	log.SetFlags(0)

	numWorkers := 100
	interrupt := make(chan os.Signal, 1)
	done := make(chan struct{}, 0)
	for i := 0; i < numWorkers; i++ {
		sleep := int64(20000) + int64(i)
		name := fmt.Sprintf("Sense%d", i)
		client := NewHelloWsClient(name, interrupt)
		go client.Start(*addr, sleep, done)
	}
	<-interrupt
	log.Println("Done exiting")
}
