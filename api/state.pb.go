// Code generated by protoc-gen-go.
// source: state.proto
// DO NOT EDIT!

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type AudioState struct {
	PlayingAudio     *bool   `protobuf:"varint,1,req,name=playing_audio" json:"playing_audio,omitempty"`
	DurationSeconds  *uint32 `protobuf:"varint,2,opt,name=duration_seconds" json:"duration_seconds,omitempty"`
	FilePath         *string `protobuf:"bytes,3,opt,name=file_path" json:"file_path,omitempty"`
	VolumePercent    *uint32 `protobuf:"varint,4,opt,name=volume_percent" json:"volume_percent,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AudioState) Reset()                    { *m = AudioState{} }
func (m *AudioState) String() string            { return proto.CompactTextString(m) }
func (*AudioState) ProtoMessage()               {}
func (*AudioState) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0} }

func (m *AudioState) GetPlayingAudio() bool {
	if m != nil && m.PlayingAudio != nil {
		return *m.PlayingAudio
	}
	return false
}

func (m *AudioState) GetDurationSeconds() uint32 {
	if m != nil && m.DurationSeconds != nil {
		return *m.DurationSeconds
	}
	return 0
}

func (m *AudioState) GetFilePath() string {
	if m != nil && m.FilePath != nil {
		return *m.FilePath
	}
	return ""
}

func (m *AudioState) GetVolumePercent() uint32 {
	if m != nil && m.VolumePercent != nil {
		return *m.VolumePercent
	}
	return 0
}

type SenseState struct {
	SenseId          *string     `protobuf:"bytes,1,req,name=sense_id" json:"sense_id,omitempty"`
	AudioState       *AudioState `protobuf:"bytes,2,opt,name=audio_state" json:"audio_state,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *SenseState) Reset()                    { *m = SenseState{} }
func (m *SenseState) String() string            { return proto.CompactTextString(m) }
func (*SenseState) ProtoMessage()               {}
func (*SenseState) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{1} }

func (m *SenseState) GetSenseId() string {
	if m != nil && m.SenseId != nil {
		return *m.SenseId
	}
	return ""
}

func (m *SenseState) GetAudioState() *AudioState {
	if m != nil {
		return m.AudioState
	}
	return nil
}

func init() {
	proto.RegisterType((*AudioState)(nil), "AudioState")
	proto.RegisterType((*SenseState)(nil), "SenseState")
}

var fileDescriptor7 = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x44, 0x8e, 0xbd, 0x4e, 0x03, 0x31,
	0x0c, 0x80, 0xd5, 0x96, 0xa1, 0xe7, 0xa8, 0xa8, 0x44, 0x02, 0x45, 0x4c, 0x47, 0xa7, 0x4e, 0x19,
	0x78, 0x02, 0xe0, 0x15, 0x6e, 0x63, 0x89, 0xa2, 0x8b, 0xe1, 0x2c, 0xe5, 0xe2, 0x28, 0x3f, 0x48,
	0xbc, 0x3d, 0x97, 0x30, 0x74, 0xb3, 0x3f, 0x59, 0xfe, 0x3e, 0x10, 0xb9, 0xd8, 0x82, 0x3a, 0x26,
	0x2e, 0x7c, 0x59, 0x00, 0xde, 0xab, 0x23, 0x9e, 0x1a, 0x93, 0x8f, 0x70, 0x8a, 0xde, 0xfe, 0x52,
	0xf8, 0x36, 0xb6, 0x51, 0xb5, 0x1b, 0xf7, 0xd7, 0xa3, 0x54, 0x70, 0x76, 0x35, 0xd9, 0x42, 0x1c,
	0x4c, 0xc6, 0x99, 0x83, 0xcb, 0x6a, 0x3f, 0xee, 0xae, 0x27, 0xf9, 0x00, 0xc3, 0x17, 0x79, 0x34,
	0xd1, 0x96, 0x45, 0x1d, 0x36, 0x34, 0xc8, 0x27, 0xb8, 0xff, 0x61, 0x5f, 0xd7, 0x0d, 0x62, 0x9a,
	0x31, 0x14, 0x75, 0xd7, 0x4e, 0x2f, 0x6f, 0x00, 0x13, 0x86, 0x8c, 0xff, 0xa6, 0x33, 0x1c, 0x73,
	0xdb, 0x0c, 0xb9, 0x2e, 0x19, 0xe4, 0x08, 0xa2, 0x3b, 0x4d, 0xcf, 0xeb, 0xff, 0xc5, 0xab, 0xd0,
	0xb7, 0xba, 0x8f, 0x17, 0x78, 0x9e, 0x79, 0xd5, 0x0b, 0x7a, 0xcf, 0x3a, 0xd7, 0x44, 0xb1, 0x6a,
	0x1b, 0x49, 0x53, 0x88, 0xb5, 0x7c, 0x1e, 0xb6, 0xf1, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x20, 0xba,
	0x67, 0xc3, 0xdc, 0x00, 0x00, 0x00,
}