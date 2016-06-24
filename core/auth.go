package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"errors"
)

func sign(message []byte, key []byte) ([]byte, error) {
	iv := make([]byte, 16)
	for i := 0; i < len(iv); i++ {
		iv[i] = byte(i)
	}

	sha_buf := sha1.Sum(message)

	padded_sha := make([]byte, 32)

	for i, c := range sha_buf {
		padded_sha[i] = c
	}

	// key[0] = 0
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(padded_sha))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded_sha)

	c := [][]byte{message, iv, ciphertext}
	resp := bytes.Join(c, []byte(""))
	return resp, nil
}

func verify(body []byte, key []byte) error {
	if len(body) <= 48 {
		return errors.New("too short")
	}
	IV_LENGTH := 16
	SIG_LENGTH := 32

	sigStartIndex := 16
	ivStartIndex := 0

	iv := body[ivStartIndex:IV_LENGTH]
	sig := body[sigStartIndex : sigStartIndex+SIG_LENGTH]
	pb := body[48:len(body)]

	sha_buf := sha1.Sum(pb)

	padded_sha := make([]byte, 32)

	for i, c := range sig {
		padded_sha[i] = c
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(padded_sha))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded_sha)

	for i, c := range sha_buf {
		if c != ciphertext[i] {
			return errors.New("DO NOT MATCH")
		}
	}

	return nil
}
