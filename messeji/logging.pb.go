// Code generated by protoc-gen-go.
// source: logging.proto
// DO NOT EDIT!

package messeji

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type RequestLog_Type int32

const (
	RequestLog_MESSAGE                 RequestLog_Type = 0
	RequestLog_RECEIVE_MESSAGE_REQUEST RequestLog_Type = 1
)

var RequestLog_Type_name = map[int32]string{
	0: "MESSAGE",
	1: "RECEIVE_MESSAGE_REQUEST",
}
var RequestLog_Type_value = map[string]int32{
	"MESSAGE":                 0,
	"RECEIVE_MESSAGE_REQUEST": 1,
}

func (x RequestLog_Type) Enum() *RequestLog_Type {
	p := new(RequestLog_Type)
	*p = x
	return p
}
func (x RequestLog_Type) String() string {
	return proto.EnumName(RequestLog_Type_name, int32(x))
}
func (x *RequestLog_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(RequestLog_Type_value, data, "RequestLog_Type")
	if err != nil {
		return err
	}
	*x = RequestLog_Type(value)
	return nil
}
func (RequestLog_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

type RequestLog struct {
	Type                  *RequestLog_Type       `protobuf:"varint,1,req,name=type,enum=RequestLog_Type" json:"type,omitempty"`
	Timestamp             *uint64                `protobuf:"varint,2,opt,name=timestamp" json:"timestamp,omitempty"`
	MessageRequest        *Message               `protobuf:"bytes,3,opt,name=message_request" json:"message_request,omitempty"`
	ReceiveMessageRequest *ReceiveMessageRequest `protobuf:"bytes,4,opt,name=receive_message_request" json:"receive_message_request,omitempty"`
	SenseId               *string                `protobuf:"bytes,5,opt,name=sense_id" json:"sense_id,omitempty"`
	XXX_unrecognized      []byte                 `json:"-"`
}

func (m *RequestLog) Reset()                    { *m = RequestLog{} }
func (m *RequestLog) String() string            { return proto.CompactTextString(m) }
func (*RequestLog) ProtoMessage()               {}
func (*RequestLog) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *RequestLog) GetType() RequestLog_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return RequestLog_MESSAGE
}

func (m *RequestLog) GetTimestamp() uint64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *RequestLog) GetMessageRequest() *Message {
	if m != nil {
		return m.MessageRequest
	}
	return nil
}

func (m *RequestLog) GetReceiveMessageRequest() *ReceiveMessageRequest {
	if m != nil {
		return m.ReceiveMessageRequest
	}
	return nil
}

func (m *RequestLog) GetSenseId() string {
	if m != nil && m.SenseId != nil {
		return *m.SenseId
	}
	return ""
}

func init() {
	proto.RegisterType((*RequestLog)(nil), "RequestLog")
	proto.RegisterEnum("RequestLog_Type", RequestLog_Type_name, RequestLog_Type_value)
}

var fileDescriptor1 = []byte{
	// 238 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0xc9, 0x4f, 0x4f,
	0xcf, 0xcc, 0x4b, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e, 0x4e,
	0xcd, 0xca, 0x84, 0x70, 0x95, 0xee, 0x33, 0x72, 0x71, 0x05, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97,
	0xf8, 0xe4, 0xa7, 0x0b, 0xc9, 0x71, 0xb1, 0x94, 0x54, 0x16, 0xa4, 0x4a, 0x30, 0x2a, 0x30, 0x69,
	0xf0, 0x19, 0x09, 0xe8, 0x21, 0xa4, 0xf4, 0x42, 0x80, 0xe2, 0x42, 0x82, 0x5c, 0x9c, 0x25, 0x99,
	0x40, 0x13, 0x4a, 0x12, 0x73, 0x0b, 0x24, 0x98, 0x14, 0x18, 0x35, 0x58, 0x84, 0x14, 0xb9, 0xf8,
	0x41, 0x46, 0x26, 0xa6, 0xa7, 0xc6, 0x17, 0x41, 0x54, 0x4b, 0x30, 0x03, 0x25, 0xb8, 0x8d, 0x38,
	0xf4, 0x7c, 0x21, 0xe2, 0x42, 0xe6, 0x5c, 0xe2, 0x45, 0xa9, 0xc9, 0xa9, 0x99, 0x65, 0xa9, 0xf1,
	0xe8, 0x4a, 0x59, 0xc0, 0x4a, 0xc5, 0x80, 0x16, 0x81, 0xe5, 0xa1, 0x3a, 0xa0, 0xd6, 0x0a, 0x09,
	0x70, 0x71, 0x14, 0xa7, 0xe6, 0x15, 0xa7, 0xc6, 0x67, 0xa6, 0x48, 0xb0, 0x02, 0x55, 0x72, 0x2a,
	0x19, 0x70, 0xb1, 0x80, 0x1d, 0xc2, 0xcd, 0xc5, 0xee, 0xeb, 0x1a, 0x1c, 0xec, 0xe8, 0xee, 0x2a,
	0xc0, 0x20, 0x24, 0xcd, 0x25, 0x1e, 0xe4, 0xea, 0xec, 0xea, 0x19, 0xe6, 0x1a, 0x0f, 0x15, 0x8c,
	0x0f, 0x72, 0x0d, 0x0c, 0x75, 0x0d, 0x0e, 0x11, 0x60, 0x74, 0x52, 0xe0, 0x12, 0x4d, 0xce, 0xcf,
	0xd5, 0xcb, 0x48, 0xcd, 0xc9, 0xc9, 0xd7, 0x83, 0x79, 0x3e, 0xb1, 0x20, 0x33, 0x8a, 0x1d, 0xca,
	0x01, 0x04, 0x00, 0x00, 0xff, 0xff, 0x02, 0x0d, 0x79, 0xb5, 0x22, 0x01, 0x00, 0x00,
}
