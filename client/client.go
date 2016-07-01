// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/base64"
	// "encoding/hex"
	"flag"
	"fmt"
	"github.com/hello/haneda/sense"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func displayName(s sense.Client) {
	log.Println(s.Id())
}

// var addr = flag.String("addr", "ws-dev.hello.is", "http service address")

var (
	addr = flag.String("addr", "0.0.0.0:8082", "http service address")
	path = flag.String("path", "protobuf", "ws path")
)

func main() {
	flag.Parse()
	// log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	done := make(chan bool, 0)
	signal.Notify(interrupt, os.Interrupt)

	for i := 0; i < 1; i++ {
		name := fmt.Sprintf("Sense%d", i)
		name = "XXXXXXXXXXXXXXXX"

		privKey := []byte("1234567891234567")
		fakeSense := &sense.Sense15{
			Name:      sense.SenseId(name),
			Interrupt: interrupt,
			Done:      done,
			PrivKey:   privKey,
		}
		// fakeSense := sense.New15(name, interrupt, done, privKey)
		fakeSense = sense.NewDefaultSenseOneFive(fakeSense)
		headers := http.Header{}
		headers.Add("Authorization", "Basic "+basicAuth(name, "foo"))
		headers.Add("X-Hello-Sense-Id", name)
		// u := url.URL{Scheme: "wss", Host: *addr, Path: "/echo"}
		u := url.URL{Scheme: "ws", Host: *addr, Path: fmt.Sprintf("/%s", *path)}
		log.Printf("connecting to %s\n", u.String())

		err := fakeSense.Connect(&u, headers)
		if err != nil {
			log.Fatal(err)
		}

		displayName(fakeSense)

		go fakeSense.Receive()
		go fakeSense.Send(time.Duration(1000 * time.Millisecond))
	}
	<-done
	log.Println("Done exiting")
}
