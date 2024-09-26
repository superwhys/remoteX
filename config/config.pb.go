// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: config/config.proto

package config

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	node "github.com/superwhys/remoteX/domain/node"
	_ "github.com/superwhys/remoteX/pkg/proto/ext"
	protocol "github.com/superwhys/remoteX/pkg/protocol"
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
	LocalNode         *node.Node                   `protobuf:"bytes,1,opt,name=local_node,json=localNode,proto3" json:"localNode" yaml:"local_node"`
	Tls               *TlsConfig                   `protobuf:"bytes,2,opt,name=tls,proto3" json:"tls" yaml:"tls"`
	TransConf         *node.NodeTransConfiguration `protobuf:"bytes,3,opt,name=trans_conf,json=transConf,proto3" json:"transConf" yaml:"trans_conf"`
	DialClients       []*protocol.Address          `protobuf:"bytes,4,rep,name=dial_clients,json=dialClients,proto3" json:"dialClients" yaml:"dial_clients"`
	HeartbeatInterval int64                        `protobuf:"varint,5,opt,name=heartbeat_interval,json=heartbeatInterval,proto3" json:"heartbeatInterval" yaml:"heartbeat_interval"`
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

func (m *Config) GetLocalNode() *node.Node {
	if m != nil {
		return m.LocalNode
	}
	return nil
}

func (m *Config) GetTls() *TlsConfig {
	if m != nil {
		return m.Tls
	}
	return nil
}

func (m *Config) GetTransConf() *node.NodeTransConfiguration {
	if m != nil {
		return m.TransConf
	}
	return nil
}

func (m *Config) GetDialClients() []*protocol.Address {
	if m != nil {
		return m.DialClients
	}
	return nil
}

func (m *Config) GetHeartbeatInterval() int64 {
	if m != nil {
		return m.HeartbeatInterval
	}
	return 0
}

func init() {
	proto.RegisterType((*TlsConfig)(nil), "config.TlsConfig")
	proto.RegisterType((*Config)(nil), "config.Config")
}

func init() { proto.RegisterFile("config/config.proto", fileDescriptor_cc332a44e926b360) }

var fileDescriptor_cc332a44e926b360 = []byte{
	// 497 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0x31, 0x6f, 0xd3, 0x4e,
	0x18, 0xc6, 0xe3, 0xbf, 0xf3, 0x0f, 0xf5, 0x85, 0x25, 0xae, 0x84, 0x42, 0x41, 0x76, 0x94, 0x66,
	0xe8, 0x14, 0x4b, 0x65, 0x41, 0x20, 0x21, 0xea, 0x4a, 0x88, 0x2e, 0x0c, 0x51, 0x07, 0xc4, 0x62,
	0x5d, 0xec, 0xb7, 0xc9, 0xc9, 0x17, 0x5f, 0x74, 0x77, 0xa1, 0xf8, 0x2b, 0x30, 0xf1, 0x09, 0x50,
	0x3f, 0x4e, 0x37, 0x32, 0x32, 0x9d, 0x44, 0xb3, 0x65, 0xcc, 0xc8, 0x84, 0xee, 0xce, 0xb5, 0x83,
	0xca, 0x92, 0xbc, 0x77, 0xef, 0xfb, 0xfc, 0x9e, 0xf7, 0xb9, 0x28, 0xe8, 0x30, 0x65, 0xc5, 0x15,
	0x99, 0x45, 0xf6, 0x6b, 0xbc, 0xe4, 0x4c, 0x32, 0xbf, 0x63, 0x4f, 0x47, 0x4f, 0x32, 0xb6, 0xc0,
	0xa4, 0x88, 0x0a, 0x96, 0x81, 0xf9, 0xb0, 0xfd, 0x23, 0x0f, 0xbe, 0xc8, 0xaa, 0x7c, 0xb6, 0xcc,
	0x67, 0x91, 0x29, 0x53, 0x46, 0xeb, 0xc2, 0x36, 0x87, 0x5f, 0x1d, 0xe4, 0x5d, 0x52, 0x71, 0x6e,
	0x68, 0xfe, 0x5b, 0xe4, 0xa5, 0xc0, 0x65, 0x72, 0x45, 0x28, 0xf4, 0x9d, 0x81, 0x73, 0xe2, 0xc5,
	0xc7, 0x5b, 0x15, 0x1e, 0xe8, 0xcb, 0x77, 0x84, 0xc2, 0x4e, 0x85, 0xdd, 0x12, 0x2f, 0xe8, 0xab,
	0xa1, 0xbe, 0x19, 0xfe, 0xfe, 0x31, 0x6a, 0xeb, 0x62, 0x52, 0x0f, 0xf8, 0xaf, 0xd1, 0x41, 0x0e,
	0xa5, 0x05, 0xfc, 0x67, 0x00, 0x83, 0xad, 0x0a, 0x1f, 0xe5, 0x50, 0x56, 0x7a, 0x64, 0xf5, 0x39,
	0x94, 0x5a, 0xee, 0xe6, 0x50, 0x4e, 0xee, 0xbb, 0xc3, 0xef, 0x6d, 0xd4, 0xa9, 0x36, 0x49, 0x11,
	0xa2, 0x2c, 0xc5, 0x34, 0xd1, 0x99, 0xcc, 0x2a, 0xdd, 0x53, 0x34, 0x36, 0x01, 0x3f, 0xb0, 0x0c,
	0xe2, 0x97, 0xb7, 0x2a, 0x74, 0xb6, 0x2a, 0xf4, 0xcc, 0x94, 0xbe, 0xda, 0xa9, 0xb0, 0x67, 0xd9,
	0x8d, 0x50, 0x5b, 0xec, 0x71, 0x6e, 0xd6, 0x23, 0x67, 0xd2, 0x28, 0xfc, 0x0b, 0xe4, 0x4a, 0x2a,
	0xcc, 0x9e, 0xdd, 0xd3, 0xde, 0xb8, 0x7a, 0xe0, 0xfa, 0x39, 0xe2, 0x51, 0x65, 0xa2, 0xa7, 0x9a,
	0xd5, 0x25, 0x15, 0x66, 0x75, 0x49, 0x85, 0x01, 0xea, 0xc2, 0x5f, 0x21, 0x24, 0x39, 0x2e, 0x44,
	0xa2, 0x21, 0x7d, 0xd7, 0x10, 0x9f, 0x37, 0xfb, 0x5e, 0xea, 0x9e, 0xa5, 0xae, 0x38, 0x96, 0x84,
	0x15, 0x4d, 0x02, 0x79, 0xdf, 0x6b, 0x12, 0x34, 0x28, 0x93, 0xa0, 0x39, 0xda, 0x04, 0xb5, 0xc2,
	0x97, 0xe8, 0x71, 0x46, 0x30, 0x4d, 0x52, 0x4a, 0xa0, 0x90, 0xa2, 0xdf, 0x1e, 0xb8, 0x26, 0x4a,
	0xfd, 0x2b, 0x9f, 0x65, 0x19, 0x07, 0x21, 0xe2, 0x37, 0x95, 0x5b, 0x57, 0x8f, 0x9f, 0xdb, 0xe9,
	0x9d, 0x0a, 0x0f, 0xad, 0xdf, 0x3e, 0x43, 0x3b, 0xfe, 0x05, 0x35, 0x9e, 0xfb, 0x3a, 0xff, 0x1a,
	0xf9, 0x73, 0xc0, 0x5c, 0x4e, 0x01, 0xcb, 0x84, 0x14, 0x12, 0xf8, 0x67, 0x4c, 0xfb, 0xff, 0x0f,
	0x9c, 0x13, 0x37, 0x7e, 0xbf, 0x55, 0x61, 0xaf, 0xee, 0x5e, 0x54, 0xcd, 0x9d, 0x0a, 0x9f, 0x5a,
	0xab, 0x87, 0x42, 0x6d, 0xf8, 0x0f, 0xde, 0xe4, 0x21, 0x25, 0x3e, 0x5b, 0xff, 0x0a, 0x5a, 0xb7,
	0x77, 0x81, 0xb3, 0xbe, 0x0b, 0x9c, 0x6f, 0x9b, 0xa0, 0x75, 0xb3, 0x09, 0x9c, 0xf5, 0x26, 0x68,
	0xfd, 0xdc, 0x04, 0xad, 0x4f, 0xc7, 0x33, 0x22, 0xe7, 0xab, 0xe9, 0x38, 0x65, 0x8b, 0x48, 0xac,
	0x96, 0xc0, 0xaf, 0xe7, 0xa5, 0x88, 0x38, 0x2c, 0x98, 0x84, 0x8f, 0xd5, 0xdf, 0x67, 0xda, 0x31,
	0x4f, 0xf3, 0xe2, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe4, 0xfd, 0xe5, 0x7c, 0x56, 0x03, 0x00,
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
	if m.HeartbeatInterval != 0 {
		i = encodeVarintConfig(dAtA, i, uint64(m.HeartbeatInterval))
		i--
		dAtA[i] = 0x28
	}
	if len(m.DialClients) > 0 {
		for iNdEx := len(m.DialClients) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DialClients[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintConfig(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.TransConf != nil {
		{
			size, err := m.TransConf.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintConfig(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Tls != nil {
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
	}
	if m.LocalNode != nil {
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
	}
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
	if m.LocalNode != nil {
		l = m.LocalNode.ProtoSize()
		n += 1 + l + sovConfig(uint64(l))
	}
	if m.Tls != nil {
		l = m.Tls.ProtoSize()
		n += 1 + l + sovConfig(uint64(l))
	}
	if m.TransConf != nil {
		l = m.TransConf.ProtoSize()
		n += 1 + l + sovConfig(uint64(l))
	}
	if len(m.DialClients) > 0 {
		for _, e := range m.DialClients {
			l = e.ProtoSize()
			n += 1 + l + sovConfig(uint64(l))
		}
	}
	if m.HeartbeatInterval != 0 {
		n += 1 + sovConfig(uint64(m.HeartbeatInterval))
	}
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
			if m.LocalNode == nil {
				m.LocalNode = &node.Node{}
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
			if m.Tls == nil {
				m.Tls = &TlsConfig{}
			}
			if err := m.Tls.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransConf", wireType)
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
			if m.TransConf == nil {
				m.TransConf = &node.NodeTransConfiguration{}
			}
			if err := m.TransConf.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DialClients", wireType)
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
			m.DialClients = append(m.DialClients, &protocol.Address{})
			if err := m.DialClients[len(m.DialClients)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HeartbeatInterval", wireType)
			}
			m.HeartbeatInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HeartbeatInterval |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
