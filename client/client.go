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

var addr = flag.String("addr", "ws-dev.hello.is:443", "http service address")

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	headers := http.Header{}
	headers.Add("Authorization", "Basic "+basicAuth("tim", "foo"))
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

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
			log.Printf("recv: %d\n", m.GetId())
		}
	}()

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

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
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}
