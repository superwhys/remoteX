// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: config/config.proto

package config

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	node "github.com/superwhys/remoteX/domain/node"
	_ "github.com/superwhys/remoteX/pkg/proto/ext"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errorutils if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type TlsConfig struct {
	CertFile string `protobuf:"bytes,1,opt,name=cert_file,json=certFile,proto3" json:"certFile" yaml:"cert"`
	KeyFile  string `protobuf:"bytes,2,opt,name=key_file,json=keyFile,proto3" json:"keyFile" yaml:"key"`
}

func (m *TlsConfig) Reset()         { *m = TlsConfig{} }
func (m *TlsConfig) String() string { return proto.CompactTextString(m) }
func (*TlsConfig) ProtoMessage()    {}
func (*TlsConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc332a44e926b360, []int{0}
}
func (m *TlsConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TlsConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TlsConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TlsConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TlsConfig.Merge(m, src)
}
func (m *TlsConfig) XXX_Size() int {
	return m.ProtoSize()
}
func (m *TlsConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_TlsConfig.DiscardUnknown(m)
}

var xxx_messageInfo_TlsConfig proto.InternalMessageInfo

func (m *TlsConfig) GetCertFile() string {
	if m != nil {
		return m.CertFile
	}
	return ""
}

func (m *TlsConfig) GetKeyFile() string {
	if m != nil {
		return m.KeyFile
	}
	return ""
}

type Config struct {
	LocalNode node.Node `protobuf:"bytes,1,opt,name=local_node,json=localNode,proto3" json:"localNode" yaml:"local_node"`
	Tls       TlsConfig `protobuf:"bytes,2,opt,name=tls,proto3" json:"tls" yaml:"tls"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc332a44e926b360, []int{1}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return m.ProtoSize()
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetLocalNode() node.Node {
	if m != nil {
		return m.LocalNode
	}
	return node.Node{}
}

func (m *Config) GetTls() TlsConfig {
	if m != nil {
		return m.Tls
	}
	return TlsConfig{}
}

func init() {
	proto.RegisterType((*TlsConfig)(nil), "config.TlsConfig")
	proto.RegisterType((*Config)(nil), "config.Config")
}

func init() { proto.RegisterFile("config/config.proto", fileDescriptor_cc332a44e926b360) }

var fileDescriptor_cc332a44e926b360 = []byte{
	// 337 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x51, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x8d, 0x29, 0x2a, 0x8d, 0x3b, 0x35, 0x48, 0x08, 0x75, 0x70, 0xaa, 0x94, 0x81, 0x29, 0x91,
	0x0a, 0x13, 0x2c, 0x10, 0x04, 0x23, 0x43, 0xc4, 0x80, 0x58, 0xaa, 0x36, 0xbd, 0xb6, 0x51, 0x9c,
	0xba, 0x4a, 0x5c, 0x41, 0x7e, 0x81, 0x89, 0x4f, 0xe0, 0x13, 0xf8, 0x8c, 0x6e, 0x74, 0x64, 0xb2,
	0x44, 0xb3, 0x75, 0xec, 0xc8, 0x84, 0x6c, 0x87, 0x86, 0xe5, 0x72, 0xb9, 0x7b, 0xef, 0xdd, 0xf3,
	0x1d, 0x3e, 0x0c, 0xd9, 0x6c, 0x1c, 0x4d, 0x3c, 0xfd, 0x71, 0xe7, 0x29, 0xe3, 0xcc, 0xaa, 0xeb,
	0xbf, 0xf6, 0xd1, 0x88, 0x25, 0x83, 0x68, 0xe6, 0xcd, 0xd8, 0x08, 0x54, 0xd0, 0xfd, 0xb6, 0x09,
	0x2f, 0x5c, 0xa7, 0xce, 0x2b, 0xc2, 0xe6, 0x03, 0xcd, 0x6e, 0x14, 0xc1, 0xba, 0xc2, 0x66, 0x08,
	0x29, 0xef, 0x8f, 0x23, 0x0a, 0xc7, 0xa8, 0x83, 0x4e, 0x4d, 0xbf, 0xbb, 0x11, 0x76, 0x43, 0x16,
	0xef, 0x22, 0x0a, 0x5b, 0x61, 0x37, 0xf3, 0x41, 0x42, 0x2f, 0x1c, 0x59, 0x71, 0x7e, 0x3e, 0x4f,
	0xf6, 0x65, 0x12, 0xec, 0x00, 0xd6, 0x25, 0x6e, 0xc4, 0x90, 0x6b, 0x81, 0x3d, 0x25, 0xd0, 0xd9,
	0x08, 0xfb, 0x20, 0x86, 0xbc, 0xe4, 0x63, 0xcd, 0x8f, 0x21, 0x97, 0xf4, 0x5a, 0x0c, 0x79, 0xf0,
	0xd7, 0x75, 0x3e, 0x10, 0xae, 0x97, 0x4e, 0xfa, 0x18, 0x53, 0x16, 0x0e, 0x68, 0x5f, 0xda, 0x56,
	0x56, 0x9a, 0x3d, 0xec, 0xaa, 0x37, 0xdc, 0xb3, 0x11, 0xf8, 0xe7, 0x4b, 0x61, 0x1b, 0x1b, 0x61,
	0x9b, 0x0a, 0x25, 0x4b, 0x5b, 0x61, 0xb7, 0xb4, 0x76, 0x45, 0x94, 0x23, 0xfe, 0xe9, 0x04, 0x15,
	0xda, 0xba, 0xc5, 0x35, 0x4e, 0x33, 0xe5, 0xb1, 0xd9, 0x6b, 0xb9, 0xe5, 0xfe, 0x76, 0xab, 0xf0,
	0x3b, 0xe5, 0x00, 0x89, 0xaa, 0x6c, 0x73, 0x9a, 0x29, 0xdb, 0x9c, 0x66, 0x81, 0x0c, 0xfe, 0xf5,
	0xea, 0x9b, 0x18, 0xcb, 0x35, 0x41, 0xab, 0x35, 0x41, 0x6f, 0x05, 0x31, 0xde, 0x0b, 0x82, 0x56,
	0x05, 0x31, 0xbe, 0x0a, 0x62, 0x3c, 0x75, 0x27, 0x11, 0x9f, 0x2e, 0x86, 0x6e, 0xc8, 0x12, 0x2f,
	0x5b, 0xcc, 0x21, 0x7d, 0x9e, 0xe6, 0x99, 0x97, 0x42, 0xc2, 0x38, 0x3c, 0x96, 0x37, 0x1b, 0xd6,
	0xd5, 0x25, 0xce, 0x7e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xef, 0x54, 0x6b, 0xcd, 0xcb, 0x01, 0x00,
	0x00,
}

func (m *TlsConfig) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TlsConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TlsConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.KeyFile) > 0 {
		i -= len(m.KeyFile)
		copy(dAtA[i:], m.KeyFile)
		i = encodeVarintConfig(dAtA, i, uint64(len(m.KeyFile)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.CertFile) > 0 {
		i -= len(m.CertFile)
		copy(dAtA[i:], m.CertFile)
		i = encodeVarintConfig(dAtA, i, uint64(len(m.CertFile)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Config) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Tls.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintConfig(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.LocalNode.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintConfig(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TlsConfig) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CertFile)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	l = len(m.KeyFile)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	return n
}

func (m *Config) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LocalNode.ProtoSize()
	n += 1 + l + sovConfig(uint64(l))
	l = m.Tls.ProtoSize()
	n += 1 + l + sovConfig(uint64(l))
	return n
}

func sovConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozConfig(x uint64) (n int) {
	return sovConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TlsConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowConfig
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
			return fmt.Errorf("proto: TlsConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TlsConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CertFile", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CertFile = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeyFile", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KeyFile = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthConfig
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
func (m *Config) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowConfig
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
			return fmt.Errorf("proto: Config: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Config: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LocalNode", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LocalNode.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tls", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Tls.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthConfig
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
func skipConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowConfig
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
					return 0, ErrIntOverflowConfig
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
					return 0, ErrIntOverflowConfig
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
				return 0, ErrInvalidLengthConfig
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupConfig
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthConfig
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthConfig        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowConfig          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupConfig = fmt.Errorf("proto: unexpected end of group")
)
