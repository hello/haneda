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
	"log"
)

var (
	HeaderLenDontMatch = errors.New("headers len don't match")
	BodyLenDontMatch   = errors.New("body len don't match")
	SigDontMatch       = errors.New("signatures don't match")
)

type SenseAuthHmacSha1 struct {
	key     []byte
	senseId SenseId
}

type MessageParser interface {
	Parse(content []byte) (*MessageParts, error)
}

type MessageSigner interface {
	Sign(mp *MessageParts) ([]byte, error)
}

type Matcher interface {
	Match(message, messageMAC, key []byte) bool
}

func NewAuth(key []byte, senseId SenseId) *SenseAuthHmacSha1 {
	return &SenseAuth{
		key:     key,
		senseId: senseId,
	}
}

func (s *SenseAuthHmacSha1) Parse(content []byte) (*MessageParts, error) {

	bbuf := bytes.NewReader(content)
	var headerLen uint32
	err := binary.Read(bbuf, binary.BigEndian, &headerLen)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	header := make([]byte, headerLen)
	n, err := bbuf.Read(header)
	if uint32(n) != headerLen {
		return nil, HeaderLenDontMatch
	}
	var bodyLen uint32
	err = binary.Read(bbuf, binary.BigEndian, &bodyLen)
	if err != nil {
		return nil, err
	}

	body := make([]byte, bodyLen)
	n, err = bbuf.Read(body)
	if uint32(n) != bodyLen {
		msg := fmt.Sprintf("error=body-len-dont-match announced_len=%d read_len=%d", bodyLen, n)
		log.Println(msg)
		return nil, BodyLenDontMatch
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

	match := s.Match(m.Body, m.Sig, s.key)

	if !match {
		msg := fmt.Sprintf("sense_id=%s error=sig-dont-match", s.senseId)
		log.Println(msg)
		return nil, SigDontMatch
	}

	return m, nil
}

func (s *SenseAuthHmacSha1) Sign(mp *MessageParts) ([]byte, error) {
	empty := make([]byte, 0)
	content := make([]byte, 0)
	bbuf := bytes.NewBuffer(content)

	headerBytes, _ := proto.Marshal(mp.Header)

	headerLen := uint32(len(headerBytes))
	err := binary.Write(bbuf, binary.BigEndian, headerLen)
	if err != nil {
		fmt.Println(err)
		return empty, err
	}

	n, err := bbuf.Write(headerBytes)

	if uint32(n) != headerLen {
		return empty, errors.New("could not write full header")
	}

	bodyLen := uint32(len(mp.Body))
	err = binary.Write(bbuf, binary.BigEndian, bodyLen)
	if err != nil {
		return empty, err
	}

	n, err = bbuf.Write(mp.Body)

	if uint32(n) != bodyLen {
		return empty, BodyLenDontMatch
	}

	hm := hmac.New(sha1.New, s.key)
	hm.Write(mp.Body)

	sig := hm.Sum(nil)

	bbuf.Write(sig)
	return bbuf.Bytes(), nil
}

func (s *SenseAuthHmacSha1) Match(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	match := hmac.Equal(messageMAC, expectedMAC)
	if !match {
		log_message := "body_hex=%X key=%X expected_mac=%X message_mac=%X sense_id=%s"

		log.Println(fmt.Sprintf(log_message, message, key[8:], expectedMAC, messageMAC, s.senseId))
	}
	return match
}
