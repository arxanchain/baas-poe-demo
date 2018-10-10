// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common/timestampwrapper.proto

/*
Package common is a generated protocol buffer package.

It is generated from these files:
	common/timestampwrapper.proto

It has these top-level messages:
	TimestampWrapper
*/
package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Time is wrapper for Timestamp
type TimestampWrapper struct {
	Time *google_protobuf.Timestamp `protobuf:"bytes,1,opt,name=time" json:"time,omitempty"`
}

func (m *TimestampWrapper) Reset()                    { *m = TimestampWrapper{} }
func (m *TimestampWrapper) String() string            { return proto.CompactTextString(m) }
func (*TimestampWrapper) ProtoMessage()               {}
func (*TimestampWrapper) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TimestampWrapper) GetTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func init() {
	proto.RegisterType((*TimestampWrapper)(nil), "common.TimestampWrapper")
}

func init() { proto.RegisterFile("common/timestampwrapper.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 159 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4d, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2f, 0xc9, 0xcc, 0x4d, 0x2d, 0x2e, 0x49, 0xcc, 0x2d, 0x28, 0x2f, 0x4a, 0x2c,
	0x28, 0x48, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0x48, 0x4b, 0xc9, 0xa7,
	0xe7, 0xe7, 0xa7, 0xe7, 0xa4, 0xea, 0x83, 0x45, 0x93, 0x4a, 0xd3, 0x10, 0xea, 0x21, 0x0a, 0x95,
	0x9c, 0xb8, 0x04, 0x42, 0x60, 0x42, 0xe1, 0x10, 0x23, 0x84, 0xf4, 0xb8, 0x58, 0x40, 0xca, 0x24,
	0x18, 0x15, 0x18, 0x35, 0xb8, 0x8d, 0xa4, 0xf4, 0x20, 0x66, 0xe8, 0xc1, 0xcc, 0xd0, 0x83, 0x6b,
	0x08, 0x02, 0xab, 0x73, 0x32, 0x8e, 0x32, 0x4c, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce,
	0xcf, 0xd5, 0x4f, 0x2c, 0xaa, 0x48, 0xcc, 0x4b, 0xce, 0x48, 0xcc, 0xcc, 0xd3, 0x2f, 0x4e, 0xc9,
	0xd6, 0x4d, 0xcf, 0xd7, 0x85, 0x3a, 0x15, 0xac, 0xbf, 0x58, 0x1f, 0xc2, 0x4b, 0x62, 0x03, 0x73,
	0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa0, 0x4c, 0x0b, 0x78, 0xc9, 0x00, 0x00, 0x00,
}
