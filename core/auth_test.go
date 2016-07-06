package core

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"
)

var (
	// signed, _ := sign([]byte("sense"), []byte("1234567891234567"))
	encodedSense = []byte{
		115, 101, 110, 115, 101, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 101,
		57, 26, 213, 59, 219, 98, 204, 227, 68, 150, 17, 211, 16, 137, 126, 27, 31, 164,
		166, 52, 143, 60, 92, 8, 70, 19, 44, 89, 66, 211, 149}

	validKey = []byte("1234567891234567")
)

func TestSign(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	var tests = []struct {
		message     []byte
		key         []byte
		signed      []byte
		shouldFail  bool
		shouldMatch bool
	}{
		{[]byte("hello"), []byte("key"), []byte{}, true, false},
		{[]byte("sense"), validKey, encodedSense, false, true},
		{[]byte("sense"), validKey, []byte("sense"), false, false},
	}

	for i, test := range tests {
		auth := suripuAuth{key: test.key}
		signed, err := auth.sign(test.message)
		if test.shouldMatch && !bytes.Equal(test.signed, signed) {
			t.Errorf("%x != %x", test.signed, signed)
		} else if test.shouldFail && err == nil {
			t.Errorf("%s should be invalid", test.key)
		} else if !test.shouldFail && err != nil {
			t.Errorf("%s - %d", "test should not fail but did fail", i)
		}
	}
}

func TestVerify(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	auth := suripuAuth{key: []byte("1234567891234567")}

	err := auth.verify([]byte("hello"))
	if err != ErrTooShort {
		t.Errorf("%v content is too short. Should have failed", err)
	}

	content := "this is some content that should be at least forty eight bytes long"
	auth = suripuAuth{key: []byte("short key")}
	err = auth.verify([]byte(content))
	if err != ErrInvalidKey {
		t.Errorf("%v key is invalid. Should have failed", err)
	}

	auth = suripuAuth{key: validKey}
	err = auth.verify([]byte(content))
	if err != ErrDontMatch {
		t.Errorf("%v content is too short. Should have failed", err)
	}
}
