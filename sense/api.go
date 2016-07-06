package sense

import (
	"github.com/hello/haneda/haneda"
)

type SenseId string

type TopFirmwareVersion string

type MiddleFirmwareVersion string

// MessageParts holds all the parts required to act on a message from or to sense
type MessageParts struct {
	Header  *haneda.Preamble
	Body    []byte  // serialized protobuf. Type of protobuf is defined by Header
	Sig     []byte  // holds the signature
	SenseId SenseId // optional SenseId for which the message is generated
}
