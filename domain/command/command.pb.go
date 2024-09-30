// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: domain/command/command.proto

package command

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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
	Empty CommandType = 0
)

var CommandType_name = map[int32]string{
	0: "EMPTY",
}

var CommandType_value = map[string]int32{
	"EMPTY": 0,
}

func (x CommandType) String() string {
	return proto.EnumName(CommandType_name, int32(x))
}

func (CommandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_12622ffd59563e51, []int{0}
}

type CommandParam struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key" yaml:"-"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value" yaml:"-"`
}

func (m *CommandParam) Reset()         { *m = CommandParam{} }
func (m *CommandParam) String() string { return proto.CompactTextString(m) }
func (*CommandParam) ProtoMessage()    {}
func (*CommandParam) Descriptor() ([]byte, []int) {
	return fileDescriptor_12622ffd59563e51, []int{0}
}
func (m *CommandParam) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CommandParam) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CommandParam.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CommandParam) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandParam.Merge(m, src)
}
func (m *CommandParam) XXX_Size() int {
	return m.ProtoSize()
}
func (m *CommandParam) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandParam.DiscardUnknown(m)
}

var xxx_messageInfo_CommandParam proto.InternalMessageInfo

func (m *CommandParam) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *CommandParam) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Command struct {
	Type   CommandType    `protobuf:"varint,1,opt,name=type,proto3,enum=command.CommandType" json:"type" yaml:"-"`
	Args   []string       `protobuf:"bytes,2,rep,name=args,proto3" json:"args" yaml:"-"`
	Params []CommandParam `protobuf:"bytes,3,rep,name=params,proto3" json:"params" yaml:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_12622ffd59563e51, []int{1}
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

func (m *Command) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *Command) GetParams() []CommandParam {
	if m != nil {
		return m.Params
	}
	return nil
}

func init() {
	proto.RegisterEnum("command.CommandType", CommandType_name, CommandType_value)
	proto.RegisterType((*CommandParam)(nil), "command.CommandParam")
	proto.RegisterType((*Command)(nil), "command.Command")
}

func init() { proto.RegisterFile("domain/command/command.proto", fileDescriptor_12622ffd59563e51) }

var fileDescriptor_12622ffd59563e51 = []byte{
	// 318 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x49, 0xc9, 0xcf, 0x4d,
	0xcc, 0xcc, 0xd3, 0x4f, 0xce, 0xcf, 0xcd, 0x4d, 0xcc, 0x4b, 0x81, 0xd1, 0x7a, 0x05, 0x45, 0xf9,
	0x25, 0xf9, 0x42, 0xec, 0x50, 0xae, 0x52, 0x12, 0x17, 0x8f, 0x33, 0x84, 0x19, 0x90, 0x58, 0x94,
	0x98, 0x2b, 0xa4, 0xca, 0xc5, 0x9c, 0x9d, 0x5a, 0x29, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0xe9, 0x24,
	0xfc, 0xea, 0x9e, 0x3c, 0x88, 0xfb, 0xe9, 0x9e, 0x3c, 0x47, 0x65, 0x62, 0x6e, 0x8e, 0x95, 0x92,
	0xae, 0x52, 0x10, 0x48, 0x40, 0x48, 0x9b, 0x8b, 0xb5, 0x2c, 0x31, 0xa7, 0x34, 0x55, 0x82, 0x09,
	0xac, 0x50, 0xf4, 0xd5, 0x3d, 0x79, 0x88, 0x00, 0x8a, 0x52, 0x88, 0x90, 0xd2, 0x3e, 0x46, 0x2e,
	0x76, 0xa8, 0x25, 0x42, 0x76, 0x5c, 0x2c, 0x25, 0x95, 0x05, 0xa9, 0x60, 0x0b, 0xf8, 0x8c, 0x44,
	0xf4, 0x60, 0xce, 0x82, 0xca, 0x87, 0x54, 0x16, 0xa4, 0x3a, 0x89, 0xbc, 0xba, 0x27, 0x0f, 0x56,
	0x85, 0x62, 0x18, 0x58, 0x44, 0x48, 0x83, 0x8b, 0x25, 0xb1, 0x28, 0xbd, 0x58, 0x82, 0x49, 0x81,
	0x59, 0x83, 0x13, 0xa2, 0x12, 0xc4, 0x47, 0x55, 0x09, 0x12, 0x11, 0xf2, 0xe4, 0x62, 0x2b, 0x00,
	0x79, 0xa9, 0x58, 0x82, 0x59, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x14, 0xdd, 0x2e, 0xb0, 0x87, 0x9d,
	0xa4, 0x4e, 0xdc, 0x93, 0x67, 0x78, 0x75, 0x4f, 0x1e, 0xaa, 0x18, 0xc5, 0x20, 0xa8, 0x98, 0x96,
	0x32, 0x17, 0x37, 0x92, 0xfb, 0x84, 0x44, 0xb8, 0x58, 0x5d, 0x7d, 0x03, 0x42, 0x22, 0x05, 0x18,
	0xa4, 0x38, 0xbb, 0xe6, 0x2a, 0xb0, 0xba, 0xe6, 0x16, 0x94, 0x54, 0x3a, 0x79, 0x5e, 0x78, 0x28,
	0xc7, 0x70, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x2c, 0x78,
	0x2c, 0xc7, 0x78, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0xda, 0xe9, 0x99, 0x25,
	0x19, 0xa5, 0x49, 0x20, 0x37, 0xe8, 0x17, 0x97, 0x16, 0xa4, 0x16, 0x95, 0x67, 0x54, 0x16, 0xeb,
	0x17, 0xa5, 0xe6, 0xe6, 0x97, 0xa4, 0x46, 0xe8, 0xa3, 0xc6, 0x55, 0x12, 0x1b, 0x38, 0x92, 0x8c,
	0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8c, 0xc0, 0xe1, 0x54, 0xc4, 0x01, 0x00, 0x00,
}

func (m *CommandParam) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CommandParam) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CommandParam) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintCommand(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintCommand(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
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
	if len(m.Params) > 0 {
		for iNdEx := len(m.Params) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Params[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCommand(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Args) > 0 {
		for iNdEx := len(m.Args) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Args[iNdEx])
			copy(dAtA[i:], m.Args[iNdEx])
			i = encodeVarintCommand(dAtA, i, uint64(len(m.Args[iNdEx])))
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
func (m *CommandParam) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovCommand(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovCommand(uint64(l))
	}
	return n
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
		for _, s := range m.Args {
			l = len(s)
			n += 1 + l + sovCommand(uint64(l))
		}
	}
	if len(m.Params) > 0 {
		for _, e := range m.Params {
			l = e.ProtoSize()
			n += 1 + l + sovCommand(uint64(l))
		}
	}
	return n
}

func sovCommand(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommand(x uint64) (n int) {
	return sovCommand(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CommandParam) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: CommandParam: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CommandParam: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
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
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
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
			m.Value = string(dAtA[iNdEx:postIndex])
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
			m.Args = append(m.Args, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
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
			m.Params = append(m.Params, CommandParam{})
			if err := m.Params[len(m.Params)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
