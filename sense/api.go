package sense

import (
	"github.com/hello/haneda/haneda"
)

type SenseId string

type TopFirmwareVersion string

type MiddleFirmwareVersion string

type MessageParts struct {
	Header  *haneda.Preamble
	Body    []byte
	Sig     []byte
	SenseId SenseId
}
