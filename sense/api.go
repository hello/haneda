package sense

import (
	"github.com/hello/haneda/haneda"
)

type MessageParts struct {
	Header  *haneda.Preamble
	Body    []byte
	Sig     []byte
	SenseId string
}
