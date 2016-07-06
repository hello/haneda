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
	ErrHeaderLenDontMatch = errors.New("headers len don't match")
	ErrBodyLenDontMatch   = errors.New("body len don't match")
	ErrSigDontMatch       = errors.New("signatures don't match")
	ErrEmptyKey           = errors.New("empty key")
	ErrEmptyBody          = errors.New("empty body")
	ErrNil                = errors.New("can't be nil")
	ErrInvalidProto       = errors.New("invalid proto")
)

// MessageParser converts []byte to *MessagePart
// used to read messages from Sense
type MessageParser interface {
	Parse(content []byte) (*MessageParts, error)
}

// MessageSigner converts *MessagePart to []byte
// used to send messages to Sense
type MessageSigner interface {
	Sign(mp *MessageParts) ([]byte, error)
}

// SenseAuthHmacSha1 holds the information necessary to parse and sign messages for a given Sense
type SenseAuthHmacSha1 struct {
	key     []byte
	senseId SenseId
	context string
}

func NewSenseAuthHmacSha1(key []byte, senseId SenseId) *SenseAuthHmacSha1 {
	return &SenseAuthHmacSha1{
		key:     key,
		senseId: senseId,
		context: string(senseId),
	}
}

// NewSenseAuthHmacSha1 creates a SenseAuthHmacSha1 struct
func NewSenseAuthHmacSha1WithContext(key []byte, senseId SenseId, context string) *SenseAuthHmacSha1 {
	return &SenseAuthHmacSha1{
		key:     key,
		senseId: senseId,
		context: context,
	}
}

// Parse is an implementation of MessageParser
func (s *SenseAuthHmacSha1) Parse(content []byte) (*MessageParts, error) {
	bbuf := bytes.NewReader(content)
	var headerLen uint32
	err := binary.Read(bbuf, binary.BigEndian, &headerLen)
	if err != nil {
		return nil, err
	}

	header := make([]byte, headerLen)
	n, err := bbuf.Read(header)
	if uint32(n) != headerLen {
		return nil, ErrHeaderLenDontMatch
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
		return nil, ErrBodyLenDontMatch
	}

	sig := make([]byte, 20)
	n, err = bbuf.Read(sig)

	preamble := &haneda.Preamble{}
	protoErr := proto.Unmarshal(header, preamble)

	if protoErr != nil {
		log.Printf("action=parse-message sense_id=%s error=%s\n", s.senseId, protoErr)
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
		return nil, ErrSigDontMatch
	}

	return m, nil
}

// Sign is an implementation of MessageSigner
func (s *SenseAuthHmacSha1) Sign(mp *MessageParts) ([]byte, error) {
	empty := make([]byte, 0)
	content := make([]byte, 0)

	if mp == nil || mp.Header == nil {
		return empty, ErrNil
	}

	if len(mp.Body) == 0 {
		return empty, ErrEmptyBody
	}

	if len(s.key) == 0 {
		return empty, ErrEmptyKey
	}

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

	if err != nil {
		return empty, err
	}

	bodyLen := uint32(len(mp.Body))
	err = binary.Write(bbuf, binary.BigEndian, bodyLen)
	if err != nil {
		return empty, err
	}

	n, err = bbuf.Write(mp.Body)

	if uint32(n) != bodyLen {
		return empty, ErrBodyLenDontMatch
	}

	hm := hmac.New(sha1.New, s.key)
	hm.Write(mp.Body)

	sig := hm.Sum(nil)

	bbuf.Write(sig)
	return bbuf.Bytes(), nil
}

// Match validates the hmac signature given the message, sig and key
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
