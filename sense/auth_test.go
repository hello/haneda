package sense

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"io"
	"testing"
)

var (
	testKey = []byte("hello")
)

func TestSign(t *testing.T) {
	defaultHeader := &haneda.Preamble{}
	defaultHeader.Id = proto.Uint64(1)

	tests := []struct {
		name        string
		key         []byte
		body        []byte
		header      *haneda.Preamble
		shouldErr   bool
		expectedErr error
	}{
		{
			name:        "ideal_case",
			key:         testKey,
			body:        []byte("body"),
			header:      defaultHeader,
			shouldErr:   false,
			expectedErr: nil,
		},
		{
			name:        "nil_header",
			key:         testKey,
			body:        []byte("body"),
			header:      nil,
			shouldErr:   true,
			expectedErr: ErrNil,
		},
		{
			name:        "empty_body",
			key:         testKey,
			body:        []byte(""),
			header:      defaultHeader,
			shouldErr:   true,
			expectedErr: ErrEmptyBody,
		},
		{
			name:        "empty_key",
			key:         []byte(""),
			body:        []byte("body"),
			header:      defaultHeader,
			shouldErr:   true,
			expectedErr: ErrEmptyKey,
		},
	}

	for _, test := range tests {
		auth := NewSenseAuthHmacSha1(test.key, "test")
		mp := &MessageParts{Header: test.header, Body: test.body}
		_, err := auth.Sign(mp)

		if !test.shouldErr && err != nil {
			t.Fatalf("name=%s, err=%v", test.name, err)
		}

		if err != test.expectedErr {
			t.Fatalf("name=%s got=%v want=%v", test.name, err, test.expectedErr)
		}
	}
}

func makeMessage(annoucedHeaderLen, annoucedBodyLen int32, header, body, sig []byte) []byte {
	internal := make([]byte, 0)
	buffer := bytes.NewBuffer(internal)
	binary.Write(buffer, binary.BigEndian, &annoucedHeaderLen)
	buffer.Write(header)
	binary.Write(buffer, binary.BigEndian, &annoucedBodyLen)
	buffer.Write(body)
	buffer.Write(sig)

	return buffer.Bytes()
}

func sig(body, key []byte) []byte {
	hm := hmac.New(sha1.New, key)
	hm.Write(body)
	sig := hm.Sum(nil)
	return sig
}

func header(header *haneda.Preamble) []byte {
	header.Type = haneda.Preamble_BATCHED_PERIODIC_DATA.Enum()
	body, err := proto.Marshal(header)
	if err != nil {
		panic(err)
	}
	return body
}

func TestParseMalformedMessages(t *testing.T) {

	tests := []struct {
		name        string
		key         []byte
		body        []byte
		expectedErr error
	}{
		{
			name:        "empty_message",
			key:         []byte(""),
			body:        []byte(""),
			expectedErr: io.EOF,
		},
		{
			name:        "header_len_too_short",
			key:         []byte(""),
			body:        []byte{10, 10, 10},
			expectedErr: io.ErrUnexpectedEOF,
		},
		{
			name:        "invalid_header_len",
			key:         []byte(""),
			body:        []byte("garbage"),
			expectedErr: ErrHeaderLenDontMatch,
		},
	}

	for _, test := range tests {
		auth := NewSenseAuthHmacSha1(test.key, "test")
		_, err := auth.Parse(test.body)
		if err == nil {
			t.Fatalf("test name=%s should have failed.", test.name)
		}

		if test.expectedErr != err {
			t.Fatalf("name=%s got=%v want=%v", test.name, err, test.expectedErr)
		}

	}
}

func TestParse(t *testing.T) {
	defaultHeader := &haneda.Preamble{}
	defaultHeader.Id = proto.Uint64(1)

	tests := []struct {
		name      string
		headerLen *int32
		header    *haneda.Preamble
		bodyLen   *int32
		body      []byte
		key       []byte
		sig       []byte
		shouldErr bool

		expectedErr error
	}{
		{
			name:        "ideal_case",
			header:      &haneda.Preamble{},
			body:        []byte("body"),
			key:         []byte("abc"),
			shouldErr:   false,
			expectedErr: nil,
		},
		// {
		// 	name:        "wrong_header_len",
		// 	header:      &haneda.Preamble{},
		// 	headerLen:   proto.Int32(1),
		// 	body:        []byte("body"),
		// 	key:         []byte("abc"),
		// 	shouldErr:   true,
		// 	expectedErr: ErrInvalidProto,
		// },
	}

	for _, test := range tests {
		auth := NewSenseAuthHmacSha1(test.key, "test")

		s := sig(test.body, test.key)
		if len(test.sig) > 0 {
			s = test.sig
		}

		headerLen := int32(len(header(test.header)))
		if test.headerLen != nil {
			headerLen = *test.headerLen
		}

		bodyLen := int32(len(test.body))
		if test.bodyLen != nil {
			bodyLen = *test.bodyLen
		}

		messageBytes := makeMessage(headerLen, bodyLen, header(test.header), test.body, s)

		_, err := auth.Parse(messageBytes)

		if !test.shouldErr && err != nil {
			t.Fatalf("name=%s, err=%v", test.name, err)
		}

		if err != test.expectedErr {
			t.Fatalf("name=%s got=%v want=%v", test.name, err, test.expectedErr)
		}
	}
}

func TestSignParse(t *testing.T) {
	auth := NewSenseAuthHmacSha1([]byte("key"), "test")
	ack := "foo"
	id := uint64(8888)

	header := &haneda.Preamble{}
	header.Id = proto.Uint64(id)
	header.Type = haneda.Preamble_ACK.Enum()

	syncResp := &api.SyncResponse{}
	syncResp.RingTimeAck = proto.String(ack)
	body, _ := proto.Marshal(syncResp)

	mp := &MessageParts{Header: header, Body: body}

	signed, err := auth.Sign(mp)
	if err != nil {
		t.Fatalf("%v", err)
	}

	rmp, err := auth.Parse(signed)
	if err != nil {
		t.Fatalf("%v", err)
	}

	rSyncResp := &api.SyncResponse{}
	protoErr := proto.Unmarshal(rmp.Body, rSyncResp)
	if protoErr != nil {
		t.Fatalf("%v", protoErr)
	}

	if mp.Header.GetId() != rmp.Header.GetId() {
		t.Fatalf("got=%d want=%d", rmp.Header.GetId(), mp.Header.GetId())
	}

	if syncResp.GetRingTimeAck() != rSyncResp.GetRingTimeAck() {
		t.Fatalf("got=%s want=%s", rSyncResp.GetRingTimeAck(), syncResp.GetRingTimeAck())
	}
}
