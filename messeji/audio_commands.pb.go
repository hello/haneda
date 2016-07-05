// Code generated by protoc-gen-go.
// source: audio_commands.proto
// DO NOT EDIT!

/*
Package messeji is a generated protocol buffer package.

It is generated from these files:
	audio_commands.proto
	logging.proto
	messeji.proto

It has these top-level messages:
	StopAudio
	PlayAudio
	RequestLog
	ReceiveMessageRequest
	Message
	BatchMessage
	MessageStatus
*/
package messeji

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type StopAudio struct {
	FadeOutDurationSeconds *uint32 `protobuf:"varint,1,req,name=fade_out_duration_seconds" json:"fade_out_duration_seconds,omitempty"`
	XXX_unrecognized       []byte  `json:"-"`
}

func (m *StopAudio) Reset()                    { *m = StopAudio{} }
func (m *StopAudio) String() string            { return proto.CompactTextString(m) }
func (*StopAudio) ProtoMessage()               {}
func (*StopAudio) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StopAudio) GetFadeOutDurationSeconds() uint32 {
	if m != nil && m.FadeOutDurationSeconds != nil {
		return *m.FadeOutDurationSeconds
	}
	return 0
}

type PlayAudio struct {
	FilePath                      *string `protobuf:"bytes,1,req,name=file_path" json:"file_path,omitempty"`
	VolumePercent                 *uint32 `protobuf:"varint,2,req,name=volume_percent" json:"volume_percent,omitempty"`
	DurationSeconds               *uint32 `protobuf:"varint,3,opt,name=duration_seconds" json:"duration_seconds,omitempty"`
	FadeInDurationSeconds         *uint32 `protobuf:"varint,4,req,name=fade_in_duration_seconds" json:"fade_in_duration_seconds,omitempty"`
	FadeOutDurationSeconds        *uint32 `protobuf:"varint,5,req,name=fade_out_duration_seconds" json:"fade_out_duration_seconds,omitempty"`
	TimeoutFadeOutDurationSeconds *uint32 `protobuf:"varint,6,opt,name=timeout_fade_out_duration_seconds" json:"timeout_fade_out_duration_seconds,omitempty"`
	XXX_unrecognized              []byte  `json:"-"`
}

func (m *PlayAudio) Reset()                    { *m = PlayAudio{} }
func (m *PlayAudio) String() string            { return proto.CompactTextString(m) }
func (*PlayAudio) ProtoMessage()               {}
func (*PlayAudio) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PlayAudio) GetFilePath() string {
	if m != nil && m.FilePath != nil {
		return *m.FilePath
	}
	return ""
}

func (m *PlayAudio) GetVolumePercent() uint32 {
	if m != nil && m.VolumePercent != nil {
		return *m.VolumePercent
	}
	return 0
}

func (m *PlayAudio) GetDurationSeconds() uint32 {
	if m != nil && m.DurationSeconds != nil {
		return *m.DurationSeconds
	}
	return 0
}

func (m *PlayAudio) GetFadeInDurationSeconds() uint32 {
	if m != nil && m.FadeInDurationSeconds != nil {
		return *m.FadeInDurationSeconds
	}
	return 0
}

func (m *PlayAudio) GetFadeOutDurationSeconds() uint32 {
	if m != nil && m.FadeOutDurationSeconds != nil {
		return *m.FadeOutDurationSeconds
	}
	return 0
}

func (m *PlayAudio) GetTimeoutFadeOutDurationSeconds() uint32 {
	if m != nil && m.TimeoutFadeOutDurationSeconds != nil {
		return *m.TimeoutFadeOutDurationSeconds
	}
	return 0
}

func init() {
	proto.RegisterType((*StopAudio)(nil), "StopAudio")
	proto.RegisterType((*PlayAudio)(nil), "PlayAudio")
}

var fileDescriptor0 = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x90, 0x4d, 0x4e, 0x04, 0x21,
	0x10, 0x85, 0xd3, 0xe3, 0x5f, 0xa8, 0x44, 0xa3, 0x1d, 0x35, 0xb8, 0x63, 0x66, 0xa5, 0x1b, 0xee,
	0xa0, 0x27, 0x30, 0x71, 0xe7, 0x86, 0x90, 0xa6, 0x3a, 0x8d, 0x01, 0x8a, 0x34, 0x60, 0xe2, 0xed,
	0x3c, 0x9a, 0x80, 0xee, 0x3a, 0xce, 0x12, 0xde, 0xf7, 0xea, 0xab, 0x14, 0xdc, 0xea, 0x62, 0x2c,
	0xa9, 0x89, 0xbc, 0xd7, 0xc1, 0x24, 0x19, 0x57, 0xca, 0x74, 0x90, 0xc0, 0xde, 0x32, 0xc5, 0xe7,
	0x96, 0x8d, 0x7b, 0x78, 0x98, 0xb5, 0x41, 0x45, 0x25, 0x2b, 0x53, 0x56, 0x9d, 0x2d, 0x05, 0x95,
	0x70, 0xa2, 0xca, 0xf3, 0x41, 0xec, 0x1e, 0x2f, 0x0f, 0xdf, 0x03, 0xb0, 0x57, 0xa7, 0xbf, 0x7e,
	0x0b, 0x37, 0xc0, 0x66, 0xeb, 0x50, 0x45, 0x9d, 0x97, 0x0e, 0xb0, 0xf1, 0x1e, 0xae, 0x3e, 0xc9,
	0x15, 0x5f, 0x3f, 0x71, 0x9d, 0x30, 0x64, 0xbe, 0x6b, 0xc5, 0x91, 0xc3, 0xf5, 0x66, 0xe4, 0x89,
	0x18, 0x6a, 0x22, 0x80, 0x77, 0xab, 0x0d, 0x5b, 0xe9, 0x69, 0xef, 0x1e, 0xdd, 0xeb, 0xac, 0x23,
	0x4f, 0xb0, 0xcf, 0xd6, 0x63, 0x23, 0xfe, 0x47, 0xcf, 0x9b, 0xef, 0x45, 0xc0, 0x5d, 0x3d, 0x82,
	0x5c, 0xd0, 0x39, 0x92, 0x1e, 0x53, 0xc2, 0x0f, 0x2b, 0x75, 0xb4, 0xef, 0x17, 0x7f, 0x8f, 0x9f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x7f, 0xa8, 0x76, 0xa1, 0x2b, 0x01, 0x00, 0x00,
}