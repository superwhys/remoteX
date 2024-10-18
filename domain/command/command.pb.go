// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: domain/command/command.proto

package command

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	_ "github.com/superwhys/remoteX/internal/proto/ext"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type CommandType int32

const (
	// EMPTY is the type of command sent during connection establishment
	Empty   CommandType = 0
	Listdir CommandType = 1
	Push    CommandType = 2
	Pull    CommandType = 3
)

var CommandType_name = map[int32]string{
	0: "EMPTY",
	1: "ListDir",
	2: "Push",
	3: "Pull",
}

var CommandType_value = map[string]int32{
	"EMPTY":   0,
	"ListDir": 1,
	"Push":    2,
	"Pull":    3,
}

func (x CommandType) String() string {
	return proto.EnumName(CommandType_name, int32(x))
}

func (CommandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_12622ffd59563e51, []int{0}
}

type Command struct {
	Type CommandType       `protobuf:"varint,1,opt,name=type,proto3,enum=command.CommandType" json:"type" yaml:"-"`
	Args map[string]string `protobuf:"bytes,2,rep,name=args,proto3" json:"args" yaml:"-" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_12622ffd59563e51, []int{0}
}
func (m *Command) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Command.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(m, src)
}
func (m *Command) XXX_Size() int {
	return m.ProtoSize()
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetType() CommandType {
	if m != nil {
		return m.Type
	}
	return Empty
}

func (m *Command) GetArgs() map[string]string {
	if m != nil {
		return m.Args
	}
	return nil
}

type MapResp struct {
	Data map[string]string `protobuf:"bytes,1,rep,name=data,proto3" json:"data" yaml:"-" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *MapResp) Reset()         { *m = MapResp{} }
func (m *MapResp) String() string { return proto.CompactTextString(m) }
func (*MapResp) ProtoMessage()    {}
func (*MapResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_12622ffd59563e51, []int{1}
}
func (m *MapResp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MapResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MapResp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MapResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MapResp.Merge(m, src)
}
func (m *MapResp) XXX_Size() int {
	return m.ProtoSize()
}
func (m *MapResp) XXX_DiscardUnknown() {
	xxx_messageInfo_MapResp.DiscardUnknown(m)
}

var xxx_messageInfo_MapResp proto.InternalMessageInfo

func (m *MapResp) GetData() map[string]string {
	if m != nil {
		return m.Data
	}
	return nil
}

type Ret struct {
	Command *Command   `protobuf:"bytes,1,opt,name=command,proto3" json:"command" yaml:"-"`
	Resp    *types.Any `protobuf:"bytes,2,opt,name=resp,proto3" json:"resp" yaml:"-"`
	ErrNo   uint64     `protobuf:"varint,3,opt,name=err_no,json=errNo,proto3" json:"errNo" yaml:"-"`
	ErrMsg  string     `protobuf:"bytes,4,opt,name=err_msg,json=errMsg,proto3" json:"errMsg" yaml:"-"`
}

func (m *Ret) Reset()         { *m = Ret{} }
func (m *Ret) String() string { return proto.CompactTextString(m) }
func (*Ret) ProtoMessage()    {}
func (*Ret) Descriptor() ([]byte, []int) {
	return fileDescriptor_12622ffd59563e51, []int{2}
}
func (m *Ret) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Ret) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Ret.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Ret) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ret.Merge(m, src)
}
func (m *Ret) XXX_Size() int {
	return m.ProtoSize()
}
func (m *Ret) XXX_DiscardUnknown() {
	xxx_messageInfo_Ret.DiscardUnknown(m)
}

var xxx_messageInfo_Ret proto.InternalMessageInfo

func (m *Ret) GetCommand() *Command {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *Ret) GetResp() *types.Any {
	if m != nil {
		return m.Resp
	}
	return nil
}

func (m *Ret) GetErrNo() uint64 {
	if m != nil {
		return m.ErrNo
	}
	return 0
}

func (m *Ret) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func init() {
	proto.RegisterEnum("command.CommandType", CommandType_name, CommandType_value)
	proto.RegisterType((*Command)(nil), "command.Command")
	proto.RegisterMapType((map[string]string)(nil), "command.Command.ArgsEntry")
	proto.RegisterType((*MapResp)(nil), "command.MapResp")
	proto.RegisterMapType((map[string]string)(nil), "command.MapResp.DataEntry")
	proto.RegisterType((*Ret)(nil), "command.Ret")
}

func init() { proto.RegisterFile("domain/command/command.proto", fileDescriptor_12622ffd59563e51) }

var fileDescriptor_12622ffd59563e51 = []byte{
	// 530 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0x4f, 0x8b, 0xd3, 0x40,
	0x18, 0xc6, 0x33, 0x6d, 0xba, 0xdd, 0x4e, 0x41, 0xca, 0x58, 0x25, 0x86, 0x25, 0x09, 0x45, 0xa1,
	0xb8, 0x9a, 0x40, 0xbd, 0xc8, 0x1e, 0x84, 0x8d, 0x5b, 0x51, 0xb4, 0xb2, 0x94, 0x3d, 0xa8, 0x17,
	0x99, 0x6e, 0xc7, 0xb4, 0x98, 0x7f, 0xcc, 0x4c, 0xd5, 0x7c, 0x85, 0x9e, 0xfc, 0x02, 0x85, 0xfd,
	0x18, 0x7e, 0x84, 0x1e, 0x7b, 0xf4, 0x14, 0x70, 0x7b, 0xcb, 0xd1, 0x93, 0xde, 0x64, 0x26, 0xc9,
	0xb2, 0xa9, 0xe7, 0x3d, 0xb5, 0xef, 0xef, 0x7d, 0xe6, 0xc9, 0x33, 0x6f, 0xde, 0xc0, 0x83, 0x69,
	0x14, 0xe0, 0x79, 0xe8, 0x9c, 0x47, 0x41, 0x80, 0xc3, 0x69, 0xf9, 0x6b, 0xc7, 0x34, 0xe2, 0x11,
	0x6a, 0x16, 0xa5, 0xde, 0x22, 0xdf, 0x78, 0xce, 0xf4, 0x7b, 0x5e, 0x14, 0x79, 0x3e, 0x71, 0x64,
	0x35, 0x59, 0x7c, 0x72, 0x70, 0x98, 0xe4, 0xad, 0xde, 0x1f, 0x00, 0x9b, 0xcf, 0xf3, 0x13, 0xe8,
	0x19, 0x54, 0x79, 0x12, 0x13, 0x0d, 0x58, 0xa0, 0x7f, 0x6b, 0xd0, 0xb5, 0x4b, 0xe3, 0xa2, 0x7f,
	0x96, 0xc4, 0xc4, 0xed, 0x66, 0xa9, 0x29, 0x55, 0xbf, 0x53, 0x73, 0x3f, 0xc1, 0x81, 0x7f, 0xd4,
	0x7b, 0xdc, 0x1b, 0x4b, 0x82, 0x5e, 0x42, 0x15, 0x53, 0x8f, 0x69, 0x35, 0xab, 0xde, 0x6f, 0x0f,
	0xf4, 0xdd, 0xf3, 0xf6, 0x31, 0xf5, 0xd8, 0x30, 0xe4, 0x34, 0x71, 0xb5, 0x75, 0x6a, 0x2a, 0xc2,
	0x49, 0xe8, 0xab, 0x4e, 0x82, 0xe8, 0xe7, 0xb0, 0x75, 0x25, 0x46, 0x0f, 0x60, 0xfd, 0x33, 0x49,
	0x64, 0xaa, 0x96, 0x7b, 0x3b, 0x4b, 0x4d, 0x51, 0x56, 0x0e, 0x09, 0x80, 0x0e, 0x61, 0xe3, 0x0b,
	0xf6, 0x17, 0x44, 0xab, 0x49, 0xe1, 0x9d, 0x2c, 0x35, 0x73, 0x50, 0x91, 0xe6, 0xe8, 0xa8, 0xf6,
	0x14, 0xf4, 0x7e, 0x00, 0xd8, 0x1c, 0xe1, 0x78, 0x4c, 0x58, 0x2c, 0xa2, 0x4f, 0x31, 0xc7, 0x1a,
	0xd8, 0x89, 0x5e, 0xf4, 0xed, 0x13, 0xcc, 0xf1, 0x4e, 0x74, 0xa1, 0xaf, 0x46, 0x17, 0x44, 0x44,
	0xbf, 0x12, 0xdf, 0x58, 0xf4, 0xbf, 0x00, 0xd6, 0xc7, 0x84, 0xa3, 0xd7, 0xb0, 0x7c, 0xdd, 0xf2,
	0x19, 0xed, 0x41, 0x67, 0x77, 0xe8, 0xae, 0xb9, 0x4e, 0x4d, 0x90, 0xa5, 0x66, 0x29, 0xbc, 0x6e,
	0x79, 0xb1, 0xb9, 0x0f, 0xc6, 0x65, 0x03, 0xbd, 0x80, 0x2a, 0x25, 0x2c, 0x96, 0x21, 0xda, 0x83,
	0xae, 0x9d, 0x2f, 0x8d, 0x5d, 0x2e, 0x8d, 0x7d, 0x1c, 0x26, 0xee, 0x41, 0xe1, 0x26, 0x95, 0xff,
	0x59, 0x49, 0x8a, 0x1e, 0xc1, 0x3d, 0x42, 0xe9, 0xc7, 0x30, 0xd2, 0xea, 0x16, 0xe8, 0xab, 0xf9,
	0x75, 0x08, 0xa5, 0x6f, 0xa3, 0xea, 0x75, 0x24, 0x42, 0x0e, 0x6c, 0x0a, 0x75, 0xc0, 0x3c, 0x4d,
	0x95, 0xb7, 0xbf, 0x9b, 0xa5, 0xa6, 0x30, 0x18, 0x31, 0xaf, 0xa2, 0x2f, 0xd8, 0x43, 0x02, 0xdb,
	0xd7, 0x16, 0x12, 0x75, 0x61, 0x63, 0x38, 0x3a, 0x3d, 0x7b, 0xdf, 0x51, 0xf4, 0xd6, 0x72, 0x65,
	0x35, 0x86, 0x41, 0xcc, 0x13, 0xa4, 0xc1, 0xe6, 0x9b, 0x39, 0xe3, 0x27, 0x73, 0xda, 0x01, 0x7a,
	0x7b, 0xb9, 0xb2, 0x64, 0x39, 0x9d, 0x53, 0x84, 0xa0, 0x7a, 0xba, 0x60, 0xb3, 0x4e, 0x4d, 0xdf,
	0x5f, 0xae, 0x2c, 0xf9, 0x3f, 0x67, 0xbe, 0xdf, 0xa9, 0x97, 0xcc, 0xf7, 0xdd, 0x57, 0x9b, 0x5f,
	0x86, 0xb2, 0xbe, 0x34, 0xc0, 0xe6, 0xd2, 0x00, 0xdf, 0xb7, 0x86, 0x72, 0xb1, 0x35, 0xc0, 0x66,
	0x6b, 0x28, 0x3f, 0xb7, 0x86, 0xf2, 0xe1, 0xd0, 0x9b, 0xf3, 0xd9, 0x62, 0x22, 0xa6, 0xed, 0xb0,
	0x45, 0x4c, 0xe8, 0xd7, 0x59, 0xc2, 0x1c, 0x4a, 0x82, 0x88, 0x93, 0x77, 0x4e, 0xf5, 0x03, 0x9d,
	0xec, 0xc9, 0x11, 0x3e, 0xf9, 0x17, 0x00, 0x00, 0xff, 0xff, 0x35, 0x9c, 0xb4, 0xe6, 0xb9, 0x03,
	0x00, 0x00,
}

func (m *Command) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Command) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Command) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Args) > 0 {
		for k := range m.Args {
			v := m.Args[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintCommand(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintCommand(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintCommand(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Type != 0 {
		i = encodeVarintCommand(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MapResp) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MapResp) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MapResp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		for k := range m.Data {
			v := m.Data[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintCommand(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintCommand(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintCommand(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Ret) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Ret) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Ret) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ErrMsg) > 0 {
		i -= len(m.ErrMsg)
		copy(dAtA[i:], m.ErrMsg)
		i = encodeVarintCommand(dAtA, i, uint64(len(m.ErrMsg)))
		i--
		dAtA[i] = 0x22
	}
	if m.ErrNo != 0 {
		i = encodeVarintCommand(dAtA, i, uint64(m.ErrNo))
		i--
		dAtA[i] = 0x18
	}
	if m.Resp != nil {
		{
			size, err := m.Resp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommand(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Command != nil {
		{
			size, err := m.Command.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommand(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommand(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommand(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Command) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovCommand(uint64(m.Type))
	}
	if len(m.Args) > 0 {
		for k, v := range m.Args {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovCommand(uint64(len(k))) + 1 + len(v) + sovCommand(uint64(len(v)))
			n += mapEntrySize + 1 + sovCommand(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *MapResp) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Data) > 0 {
		for k, v := range m.Data {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovCommand(uint64(len(k))) + 1 + len(v) + sovCommand(uint64(len(v)))
			n += mapEntrySize + 1 + sovCommand(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *Ret) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Command != nil {
		l = m.Command.ProtoSize()
		n += 1 + l + sovCommand(uint64(l))
	}
	if m.Resp != nil {
		l = m.Resp.ProtoSize()
		n += 1 + l + sovCommand(uint64(l))
	}
	if m.ErrNo != 0 {
		n += 1 + sovCommand(uint64(m.ErrNo))
	}
	l = len(m.ErrMsg)
	if l > 0 {
		n += 1 + l + sovCommand(uint64(l))
	}
	return n
}

func sovCommand(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommand(x uint64) (n int) {
	return sovCommand(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Command) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Command: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Command: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= CommandType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Args", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommand
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Args == nil {
				m.Args = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCommand
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCommand
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthCommand
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthCommand
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCommand
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthCommand
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthCommand
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCommand(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthCommand
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Args[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MapResp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MapResp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MapResp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommand
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowCommand
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCommand
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthCommand
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthCommand
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowCommand
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthCommand
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthCommand
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipCommand(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthCommand
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Data[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Ret) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Ret: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Ret: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Command", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommand
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Command == nil {
				m.Command = &Command{}
			}
			if err := m.Command.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Resp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommand
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Resp == nil {
				m.Resp = &types.Any{}
			}
			if err := m.Resp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrNo", wireType)
			}
			m.ErrNo = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ErrNo |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrMsg", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommand
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ErrMsg = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCommand(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommand
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthCommand
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommand
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommand
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommand        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommand          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommand = fmt.Errorf("proto: unexpected end of group")
)
