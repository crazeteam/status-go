// Code generated by protoc-gen-go. DO NOT EDIT.
// source: query.proto

package cbridge

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TransferHistoryStatus int32

const (
	TransferHistoryStatus_TRANSFER_UNKNOWN                      TransferHistoryStatus = 0
	TransferHistoryStatus_TRANSFER_SUBMITTING                   TransferHistoryStatus = 1
	TransferHistoryStatus_TRANSFER_FAILED                       TransferHistoryStatus = 2
	TransferHistoryStatus_TRANSFER_WAITING_FOR_SGN_CONFIRMATION TransferHistoryStatus = 3
	TransferHistoryStatus_TRANSFER_WAITING_FOR_FUND_RELEASE     TransferHistoryStatus = 4
	TransferHistoryStatus_TRANSFER_COMPLETED                    TransferHistoryStatus = 5
	TransferHistoryStatus_TRANSFER_TO_BE_REFUNDED               TransferHistoryStatus = 6
	TransferHistoryStatus_TRANSFER_REQUESTING_REFUND            TransferHistoryStatus = 7
	TransferHistoryStatus_TRANSFER_REFUND_TO_BE_CONFIRMED       TransferHistoryStatus = 8
	TransferHistoryStatus_TRANSFER_CONFIRMING_YOUR_REFUND       TransferHistoryStatus = 9
	TransferHistoryStatus_TRANSFER_REFUNDED                     TransferHistoryStatus = 10
	TransferHistoryStatus_TRANSFER_DELAYED                      TransferHistoryStatus = 11
)

var TransferHistoryStatus_name = map[int32]string{
	0:  "TRANSFER_UNKNOWN",
	1:  "TRANSFER_SUBMITTING",
	2:  "TRANSFER_FAILED",
	3:  "TRANSFER_WAITING_FOR_SGN_CONFIRMATION",
	4:  "TRANSFER_WAITING_FOR_FUND_RELEASE",
	5:  "TRANSFER_COMPLETED",
	6:  "TRANSFER_TO_BE_REFUNDED",
	7:  "TRANSFER_REQUESTING_REFUND",
	8:  "TRANSFER_REFUND_TO_BE_CONFIRMED",
	9:  "TRANSFER_CONFIRMING_YOUR_REFUND",
	10: "TRANSFER_REFUNDED",
	11: "TRANSFER_DELAYED",
}

var TransferHistoryStatus_value = map[string]int32{
	"TRANSFER_UNKNOWN":                      0,
	"TRANSFER_SUBMITTING":                   1,
	"TRANSFER_FAILED":                       2,
	"TRANSFER_WAITING_FOR_SGN_CONFIRMATION": 3,
	"TRANSFER_WAITING_FOR_FUND_RELEASE":     4,
	"TRANSFER_COMPLETED":                    5,
	"TRANSFER_TO_BE_REFUNDED":               6,
	"TRANSFER_REQUESTING_REFUND":            7,
	"TRANSFER_REFUND_TO_BE_CONFIRMED":       8,
	"TRANSFER_CONFIRMING_YOUR_REFUND":       9,
	"TRANSFER_REFUNDED":                     10,
	"TRANSFER_DELAYED":                      11,
}

func (x TransferHistoryStatus) String() string {
	return proto.EnumName(TransferHistoryStatus_name, int32(x))
}

func (TransferHistoryStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{0}
}

type Token struct {
	Symbol               string   `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Decimal              int32    `protobuf:"varint,3,opt,name=decimal,proto3" json:"decimal,omitempty"`
	XferDisabled         bool     `protobuf:"varint,4,opt,name=xfer_disabled,json=xferDisabled,proto3" json:"xfer_disabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{0}
}

func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (m *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(m, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *Token) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Token) GetDecimal() int32 {
	if m != nil {
		return m.Decimal
	}
	return 0
}

func (m *Token) GetXferDisabled() bool {
	if m != nil {
		return m.XferDisabled
	}
	return false
}

func init() {
	proto.RegisterEnum("cbridge.TransferHistoryStatus", TransferHistoryStatus_name, TransferHistoryStatus_value)
	proto.RegisterType((*Token)(nil), "cbridge.Token")
}

func init() {
	proto.RegisterFile("query.proto", fileDescriptor_5c6ac9b241082464)
}

var fileDescriptor_5c6ac9b241082464 = []byte{
	// 373 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x51, 0x6f, 0xd3, 0x30,
	0x14, 0x85, 0x49, 0xb7, 0xb6, 0xdb, 0xdd, 0x10, 0xc6, 0x63, 0x5b, 0x04, 0x12, 0x14, 0xa6, 0x49,
	0x85, 0x87, 0xf2, 0xc0, 0x2f, 0x48, 0xe7, 0x9b, 0x11, 0x91, 0x3a, 0xe0, 0x38, 0x9a, 0xca, 0x8b,
	0x95, 0x34, 0x2e, 0x8a, 0x68, 0x1b, 0x70, 0x52, 0xa9, 0xfd, 0xe9, 0xbc, 0xa1, 0xa4, 0x69, 0x44,
	0xd1, 0x1e, 0xcf, 0x77, 0xce, 0x3d, 0xb6, 0xae, 0x2e, 0x9c, 0xfd, 0x5e, 0x6b, 0xb3, 0x1d, 0xfd,
	0x32, 0x79, 0x99, 0xd3, 0xfe, 0x2c, 0x31, 0x59, 0xfa, 0x43, 0xbf, 0xdb, 0x40, 0x57, 0xe6, 0x3f,
	0xf5, 0x8a, 0x5e, 0x41, 0xaf, 0xd8, 0x2e, 0x93, 0x7c, 0x61, 0x5b, 0x03, 0x6b, 0x78, 0x2a, 0x1a,
	0x45, 0x6d, 0xe8, 0xc7, 0x69, 0x6a, 0x74, 0x51, 0xd8, 0x9d, 0xda, 0xd8, 0xcb, 0xca, 0x49, 0xf5,
	0x2c, 0x5b, 0xc6, 0x0b, 0xfb, 0x68, 0x60, 0x0d, 0xbb, 0x62, 0x2f, 0xe9, 0x0d, 0x3c, 0xdd, 0xcc,
	0xb5, 0x51, 0x69, 0x56, 0xc4, 0xc9, 0x42, 0xa7, 0xf6, 0xf1, 0xc0, 0x1a, 0x9e, 0x88, 0xf3, 0x0a,
	0xb2, 0x86, 0x7d, 0xf8, 0xd3, 0x81, 0x4b, 0x69, 0xe2, 0x55, 0x31, 0xd7, 0xe6, 0x73, 0x56, 0x94,
	0xb9, 0xd9, 0x86, 0x65, 0x5c, 0xae, 0x0b, 0xfa, 0x02, 0x88, 0x14, 0x0e, 0x0f, 0x5d, 0x14, 0x2a,
	0xe2, 0x5f, 0x78, 0xf0, 0xc0, 0xc9, 0x13, 0x7a, 0x0d, 0x17, 0x2d, 0x0d, 0xa3, 0xf1, 0xc4, 0x93,
	0xd2, 0xe3, 0xf7, 0xc4, 0xa2, 0x17, 0xf0, 0xac, 0x35, 0x5c, 0xc7, 0xf3, 0x91, 0x91, 0x0e, 0x7d,
	0x0f, 0xb7, 0x2d, 0x7c, 0x70, 0xbc, 0x2a, 0xaa, 0xdc, 0x40, 0xa8, 0xf0, 0x9e, 0xab, 0xbb, 0x80,
	0xbb, 0x9e, 0x98, 0x38, 0xd2, 0x0b, 0x38, 0x39, 0xa2, 0xb7, 0xf0, 0xf6, 0xd1, 0xa8, 0x1b, 0x71,
	0xa6, 0x04, 0xfa, 0xe8, 0x84, 0x48, 0x8e, 0xe9, 0x15, 0xd0, 0x36, 0x76, 0x17, 0x4c, 0xbe, 0xfa,
	0x28, 0x91, 0x91, 0x2e, 0x7d, 0x05, 0xd7, 0x2d, 0x97, 0x81, 0x1a, 0xa3, 0x12, 0x58, 0x8d, 0x22,
	0x23, 0x3d, 0xfa, 0x1a, 0x5e, 0xb6, 0xa6, 0xc0, 0x6f, 0x11, 0x86, 0x75, 0xfd, 0x2e, 0x41, 0xfa,
	0xf4, 0x06, 0xde, 0xfc, 0xe3, 0xd7, 0x2f, 0xee, 0x3a, 0x9a, 0x3f, 0x22, 0x23, 0x27, 0x07, 0xa1,
	0x86, 0x57, 0x25, 0xd3, 0x20, 0xda, 0x0f, 0x91, 0x53, 0x7a, 0x09, 0xcf, 0xff, 0x6b, 0x42, 0x46,
	0xe0, 0x60, 0x97, 0x0c, 0x7d, 0x67, 0x8a, 0x8c, 0x9c, 0x8d, 0xcf, 0xbf, 0xc3, 0x68, 0xf4, 0xb1,
	0xb9, 0x81, 0xa4, 0x57, 0xdf, 0xc4, 0xa7, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x66, 0xef, 0xfe,
	0x5f, 0x22, 0x02, 0x00, 0x00,
}