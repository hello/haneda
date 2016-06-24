package sense

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/haneda"
)

func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

type SenseAuth struct {
	key []byte
}

func NewAuth(key []byte) *SenseAuth {
	return &SenseAuth{
		key: key,
	}
}

func (s *SenseAuth) Parse(content []byte) (*MessageParts, error) {

	bbuf := bytes.NewReader(content)
	var headerLen uint64
	err := binary.Read(bbuf, binary.LittleEndian, &headerLen)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	header := make([]byte, headerLen)
	n, err := bbuf.Read(header)

	if uint64(n) != headerLen {
		return nil, errors.New("header Len don't match")
	}
	var bodyLen uint64
	err = binary.Read(bbuf, binary.LittleEndian, &bodyLen)
	if err != nil {
		return nil, err
	}

	body := make([]byte, bodyLen)
	n, err = bbuf.Read(body)

	if uint64(n) != bodyLen {
		return nil, errors.New("body Len don't match")
	}

	sig := make([]byte, 20)
	n, err = bbuf.Read(sig)

	preamble := &haneda.Preamble{}
	protoErr := proto.Unmarshal(header, preamble)

	if protoErr != nil {
		return nil, protoErr
	}

	m := &MessageParts{
		Header: preamble,
		Body:   body,
		Sig:    sig,
	}

	match := CheckMAC(m.Body, m.Sig, s.key)
	if !match {
		fmt.Println("don't match!!!")
		fmt.Printf("%v\n", m.Body)
		fmt.Printf("%x\n", m.Sig)
		return nil, errors.New("sig don't match")
	}

	return m, nil
}

func (s *SenseAuth) Sign(mp *MessageParts) ([]byte, error) {
	empty := make([]byte, 0)
	content := make([]byte, 0)
	bbuf := bytes.NewBuffer(content)

	headerBytes, _ := proto.Marshal(mp.Header)

	headerLen := uint64(len(headerBytes))
	err := binary.Write(bbuf, binary.LittleEndian, headerLen)
	if err != nil {
		fmt.Println(err)
		return empty, err
	}

	n, err := bbuf.Write(headerBytes)

	if uint64(n) != headerLen {
		return empty, errors.New("could not write full header")
	}

	bodyLen := uint64(len(mp.Body))
	err = binary.Write(bbuf, binary.LittleEndian, bodyLen)
	if err != nil {
		return empty, err
	}

	n, err = bbuf.Write(mp.Body)

	if uint64(n) != bodyLen {
		return empty, errors.New("body Len don't match")
	}

	hm := hmac.New(sha1.New, s.key)
	hm.Write(mp.Body)

	sig := hm.Sum(nil)

	bbuf.Write(sig)
	return bbuf.Bytes(), nil
}
