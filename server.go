package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	// "github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"net/http"
	"strings"
)

type WSHandler struct {
}

func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func (h *WSHandler) Get(username, password string) bool {
	if username == "tim" && password == "foo" {
		return true
	}
	return false
}

func (h WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("origin: %s\n", r.Header.Get("Origin"))

	fmt.Printf("%v\n", r.Header)
	var sense string
	if len(r.Header["Authorization"]) > 0 {

		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !h.Get(pair[0], pair[1]) {
			fmt.Println("Bad auth")
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		sense = pair[0]
	} else {
		fmt.Println(r.Header)
		fmt.Println("yo")
		http.Error(w, http.StatusText(401), 401)
		return
	}

	fmt.Printf("Method: %s\n", r.Method)
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go echo(conn, sense)
}

func basicAuth(auth string) (string, string) {
	payload, _ := base64.StdEncoding.DecodeString(auth)
	pair := strings.SplitN(string(payload), ":", 2)
	return pair[0], pair[1]
}

func echo(conn *websocket.Conn, senseId string) {
	i := 0
	for {

		_, content, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading json.", err)
			break
		}

		bbuf := bytes.NewReader(content)
		var llen uint64
		err = binary.Read(bbuf, binary.LittleEndian, &llen)
		if err != nil {
			fmt.Println(err)
			break
		}

		payload := make([]byte, llen)

		_, err = bbuf.Read(payload)

		sig := make([]byte, 20)
		_, err = bbuf.Read(sig)

		match := CheckMAC(payload, sig, []byte("abc"))
		if !match {
			fmt.Println("don't match!!!")
			fmt.Printf("%v\n", payload)
			fmt.Println("len(content)", len(content))

			fmt.Println("llen", llen)
			fmt.Printf("%x\n", sig)
		}

		pb := &haneda.Preamble{}
		protoErr := proto.Unmarshal(payload, pb)
		if protoErr != nil {
			fmt.Println(protoErr)
			break
		}

		fmt.Printf("Got message: %s %d\n", pb.GetType().String(), pb.GetId())
		newId := uint64(88888)
		pb.Id = &newId
		m, _ := proto.Marshal(pb)
		if err = conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
			fmt.Println(err)
			break
		}
		i++
	}
	fmt.Println(i)
	conn.Close()
}

// func wsHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("%v\n", r)
// 	// if r.Header.Get("Origin") != "http://"+r.Host {
// 	// 	fmt.Println("Origin not allowed")
// 	// 	http.Error(w, "Origin not allowed", 403)
// 	// 	return
// 	// }
// 	fmt.Printf("Method: %s\n", r.Method)
// 	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
// 	if err != nil {
// 		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
// 	}

// 	go echo(conn)
// }

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/health")
	fmt.Fprintf(w, "%s\n", "ok")
}

func main() {

	wsHandler := WSHandler{}
	http.HandleFunc("/health", healthHandler)
	http.Handle("/echo", wsHandler)
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
