// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: domain/node/node.proto

package node

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/superwhys/remoteX/internal/proto/ext"
	github_com_superwhys_remoteX_pkg_common "github.com/superwhys/remoteX/pkg/common"
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

type NodeStatus int32

const (
	NodeStatusUnknown    NodeStatus = 0
	NodeStatusOnline     NodeStatus = 1
	NodeStatusOffline    NodeStatus = 2
	NodeStatusConnecting NodeStatus = 3
)

var NodeStatus_name = map[int32]string{
	0: "NODE_STATUS_UNKNOWN",
	1: "NODE_STATUS_ONLINE",
	2: "NODE_STATUS_OFFLINE",
	3: "NODE_STATUS_CONNECTING",
}

var NodeStatus_value = map[string]int32{
	"NODE_STATUS_UNKNOWN":    0,
	"NODE_STATUS_ONLINE":     1,
	"NODE_STATUS_OFFLINE":    2,
	"NODE_STATUS_CONNECTING": 3,
}

func (x NodeStatus) String() string {
	return proto.EnumName(NodeStatus_name, int32(x))
}

func (NodeStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_80ab9acc8f1311ae, []int{0}
}

type NodeOS int32

const (
	NodeOsUnknown NodeOS = 0
	NodeOsWin     NodeOS = 1
	NodeOsLinux   NodeOS = 2
	NodeOsDarwin  NodeOS = 3
)

var NodeOS_name = map[int32]string{
	0: "NODE_OS_UNKNOWN",
	1: "NODE_OS_WIN",
	2: "NODE_OS_LINUX",
	3: "NODE_OS_DARWIN",
}

var NodeOS_value = map[string]int32{
	"NODE_OS_UNKNOWN": 0,
	"NODE_OS_WIN":     1,
	"NODE_OS_LINUX":   2,
	"NODE_OS_DARWIN":  3,
}

func (x NodeOS) String() string {
	return proto.EnumName(NodeOS_name, int32(x))
}

func (NodeOS) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_80ab9acc8f1311ae, []int{1}
}

type NodeArch int32

const (
	NodeArchUnknown NodeArch = 0
	NodeArchX86     NodeArch = 1
	NodeArchArm     NodeArch = 2
)

var NodeArch_name = map[int32]string{
	0: "NODE_ARCH_UNKNOWN",
	1: "NODE_ARCH_X86",
	2: "NODE_ARCH_ARM",
}

var NodeArch_value = map[string]int32{
	"NODE_ARCH_UNKNOWN": 0,
	"NODE_ARCH_X86":     1,
	"NODE_ARCH_ARM":     2,
}

func (x NodeArch) String() string {
	return proto.EnumName(NodeArch_name, int32(x))
}

func (NodeArch) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_80ab9acc8f1311ae, []int{2}
}

type NodeConnectRole int32

const (
	NodeConnectRoleUnknown NodeConnectRole = 0
	NodeConnectRoleServer  NodeConnectRole = 1
	NodeConnectRoleClient  NodeConnectRole = 2
)

var NodeConnectRole_name = map[int32]string{
	0: "NODE_CONNECT_ROLE_UNKNOWN",
	1: "NODE_CONNECT_ROLE_SERVER",
	2: "NODE_CONNECT_ROLE_CLIENT",
}

var NodeConnectRole_value = map[string]int32{
	"NODE_CONNECT_ROLE_UNKNOWN": 0,
	"NODE_CONNECT_ROLE_SERVER":  1,
	"NODE_CONNECT_ROLE_CLIENT":  2,
}

func (x NodeConnectRole) String() string {
	return proto.EnumName(NodeConnectRole_name, int32(x))
}

func (NodeConnectRole) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_80ab9acc8f1311ae, []int{3}
}

type NodeConfiguration struct {
	Os             NodeOS                  `protobuf:"varint,1,opt,name=os,proto3,enum=node.NodeOS" json:"os" yaml:"-"`
	Arch           NodeArch                `protobuf:"varint,2,opt,name=arch,proto3,enum=node.NodeArch" json:"arch" yaml:"-"`
	Transmission   *NodeTransConfiguration `protobuf:"bytes,3,opt,name=transmission,proto3" json:"transmission" yaml:"-"`
	AdditionalInfo map[string]string       `protobuf:"bytes,4,rep,name=additional_info,json=additionalInfo,proto3" json:"additionalInfo" yaml:"-" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *NodeConfiguration) Reset()         { *m = NodeConfiguration{} }
func (m *NodeConfiguration) String() string { return proto.CompactTextString(m) }
func (*NodeConfiguration) ProtoMessage()    {}
func (*NodeConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ab9acc8f1311ae, []int{0}
}
func (m *NodeConfiguration) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NodeConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NodeConfiguration.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NodeConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeConfiguration.Merge(m, src)
}
func (m *NodeConfiguration) XXX_Size() int {
	return m.ProtoSize()
}
func (m *NodeConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_NodeConfiguration proto.InternalMessageInfo

func (m *NodeConfiguration) GetOs() NodeOS {
	if m != nil {
		return m.Os
	}
	return NodeOsUnknown
}

func (m *NodeConfiguration) GetArch() NodeArch {
	if m != nil {
		return m.Arch
	}
	return NodeArchUnknown
}

func (m *NodeConfiguration) GetTransmission() *NodeTransConfiguration {
	if m != nil {
		return m.Transmission
	}
	return nil
}

func (m *NodeConfiguration) GetAdditionalInfo() map[string]string {
	if m != nil {
		return m.AdditionalInfo
	}
	return nil
}

type NodeTransConfiguration struct {
	MaxSendKbps int `protobuf:"varint,1,opt,name=max_send_kbps,json=maxSendKbps,proto3,casttype=int" json:"maxSendKbps" yaml:"max_send_kbps"`
	MaxRecvKbps int `protobuf:"varint,2,opt,name=max_recv_kbps,json=maxRecvKbps,proto3,casttype=int" json:"maxRecvKbps" yaml:"max_recv_kbps"`
}

func (m *NodeTransConfiguration) Reset()         { *m = NodeTransConfiguration{} }
func (m *NodeTransConfiguration) String() string { return proto.CompactTextString(m) }
func (*NodeTransConfiguration) ProtoMessage()    {}
func (*NodeTransConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ab9acc8f1311ae, []int{1}
}
func (m *NodeTransConfiguration) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NodeTransConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NodeTransConfiguration.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NodeTransConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeTransConfiguration.Merge(m, src)
}
func (m *NodeTransConfiguration) XXX_Size() int {
	return m.ProtoSize()
}
func (m *NodeTransConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeTransConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_NodeTransConfiguration proto.InternalMessageInfo

func (m *NodeTransConfiguration) GetMaxSendKbps() int {
	if m != nil {
		return m.MaxSendKbps
	}
	return 0
}

func (m *NodeTransConfiguration) GetMaxRecvKbps() int {
	if m != nil {
		return m.MaxRecvKbps
	}
	return 0
}

type Node struct {
	NodeId        github_com_superwhys_remoteX_pkg_common.NodeID `protobuf:"bytes,1,opt,name=id,proto3,customtype=github.com/superwhys/remoteX/pkg/common.NodeID" json:"node_id" yaml:"-"`
	ConnectionId  string                                         `protobuf:"bytes,2,opt,name=connection_id,json=connectionId,proto3" json:"connectionId" yaml:"-"`
	Name          string                                         `protobuf:"bytes,3,opt,name=name,proto3" json:"name" yaml:"name"`
	Address       protocol.Address                               `protobuf:"bytes,4,opt,name=address,proto3" json:"address" yaml:"address"`
	Status        NodeStatus                                     `protobuf:"varint,5,opt,name=status,proto3,enum=node.NodeStatus" json:"status" yaml:"-"`
	IsLocal       bool                                           `protobuf:"varint,6,opt,name=is_local,json=isLocal,proto3" json:"isLocal" yaml:"-"`
	Role          NodeConnectRole                                `protobuf:"varint,7,opt,name=role,proto3,enum=node.NodeConnectRole" json:"role" yaml:"-"`
	Configuration *NodeConfiguration                             `protobuf:"bytes,8,opt,name=configuration,proto3" json:"configuration" yaml:"-"`
	LastHeartbeat int64                                          `protobuf:"varint,9,opt,name=last_heartbeat,json=lastHeartbeat,proto3" json:"lastHeartbeat" yaml:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_80ab9acc8f1311ae, []int{2}
}
func (m *Node) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Node.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return m.ProtoSize()
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetConnectionId() string {
	if m != nil {
		return m.ConnectionId
	}
	return ""
}

func (m *Node) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Node) GetAddress() protocol.Address {
	if m != nil {
		return m.Address
	}
	return protocol.Address{}
}

func (m *Node) GetStatus() NodeStatus {
	if m != nil {
		return m.Status
	}
	return NodeStatusUnknown
}

func (m *Node) GetIsLocal() bool {
	if m != nil {
		return m.IsLocal
	}
	return false
}

func (m *Node) GetRole() NodeConnectRole {
	if m != nil {
		return m.Role
	}
	return NodeConnectRoleUnknown
}

func (m *Node) GetConfiguration() *NodeConfiguration {
	if m != nil {
		return m.Configuration
	}
	return nil
}

func (m *Node) GetLastHeartbeat() int64 {
	if m != nil {
		return m.LastHeartbeat
	}
	return 0
}

func init() {
	proto.RegisterEnum("node.NodeStatus", NodeStatus_name, NodeStatus_value)
	proto.RegisterEnum("node.NodeOS", NodeOS_name, NodeOS_value)
	proto.RegisterEnum("node.NodeArch", NodeArch_name, NodeArch_value)
	proto.RegisterEnum("node.NodeConnectRole", NodeConnectRole_name, NodeConnectRole_value)
	proto.RegisterType((*NodeConfiguration)(nil), "node.NodeConfiguration")
	proto.RegisterMapType((map[string]string)(nil), "node.NodeConfiguration.AdditionalInfoEntry")
	proto.RegisterType((*NodeTransConfiguration)(nil), "node.NodeTransConfiguration")
	proto.RegisterType((*Node)(nil), "node.Node")
}

func init() { proto.RegisterFile("domain/node/node.proto", fileDescriptor_80ab9acc8f1311ae) }

var fileDescriptor_80ab9acc8f1311ae = []byte{
	// 1110 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x56, 0xcd, 0x6f, 0xe2, 0x46,
	0x14, 0xc7, 0x40, 0xb2, 0x61, 0xd8, 0x10, 0x32, 0xf9, 0x58, 0x2f, 0xdd, 0x62, 0x0b, 0xb1, 0x15,
	0xca, 0x6e, 0x49, 0x95, 0xad, 0xb6, 0xe9, 0x5e, 0xb6, 0x98, 0xb0, 0x5d, 0xb4, 0xd4, 0x48, 0x26,
	0x69, 0xd2, 0x5e, 0x90, 0x83, 0x87, 0xc4, 0x0a, 0x9e, 0x41, 0xb6, 0xc9, 0xc7, 0xa5, 0xea, 0xb1,
	0xe2, 0x54, 0xf5, 0x58, 0x09, 0x69, 0xff, 0x92, 0xf6, 0x9a, 0x5b, 0x38, 0x56, 0x3d, 0x58, 0x6a,
	0x72, 0xe3, 0x98, 0x63, 0x4e, 0xd5, 0x8c, 0x31, 0x78, 0x48, 0xda, 0x4b, 0x34, 0x7e, 0xbf, 0x8f,
	0xf7, 0xe6, 0xcd, 0xbc, 0x21, 0x60, 0xdd, 0x20, 0x96, 0x6e, 0xe2, 0x4d, 0x4c, 0x0c, 0xc4, 0xfe,
	0x14, 0xbb, 0x36, 0x71, 0x09, 0x8c, 0xd3, 0x75, 0x26, 0x81, 0xce, 0x5d, 0x3f, 0x90, 0xf9, 0xa4,
	0x7b, 0x72, 0xb4, 0xc9, 0x96, 0x2d, 0xd2, 0x99, 0x2c, 0x7c, 0x30, 0x37, 0x8c, 0x81, 0x65, 0x95,
	0x18, 0xa8, 0x4c, 0x70, 0xdb, 0x3c, 0xea, 0xd9, 0xba, 0x6b, 0x12, 0x0c, 0xbf, 0x00, 0x51, 0xe2,
	0x88, 0x82, 0x2c, 0x14, 0x52, 0x5b, 0x8f, 0x8b, 0xcc, 0x9c, 0x92, 0xea, 0x0d, 0x05, 0x8e, 0x3c,
	0x29, 0x4a, 0x9c, 0x5b, 0x4f, 0x5a, 0xb8, 0xd0, 0xad, 0xce, 0x9b, 0xdc, 0xe7, 0x39, 0x2d, 0x4a,
	0x1c, 0xb8, 0x0d, 0xe2, 0xba, 0xdd, 0x3a, 0x16, 0xa3, 0x4c, 0x93, 0x9a, 0x6a, 0x4a, 0x76, 0xeb,
	0x58, 0x59, 0x1d, 0x79, 0x12, 0xc3, 0x39, 0x1d, 0x8b, 0xc0, 0x36, 0x78, 0xec, 0xda, 0x3a, 0x76,
	0x2c, 0xd3, 0x71, 0x4c, 0x82, 0xc5, 0x98, 0x2c, 0x14, 0x92, 0x5b, 0xcf, 0xa6, 0x0e, 0xbb, 0x14,
	0xe5, 0xea, 0x53, 0xf2, 0x97, 0x9e, 0x24, 0x8c, 0x3c, 0x89, 0x53, 0x86, 0xbd, 0x3f, 0x0e, 0xf3,
	0x82, 0xc6, 0xa1, 0xf0, 0x1c, 0x2c, 0xe9, 0x86, 0x61, 0x52, 0xbd, 0xde, 0x69, 0x9a, 0xb8, 0x4d,
	0xc4, 0xb8, 0x1c, 0x2b, 0x24, 0xb7, 0x5e, 0x4c, 0x53, 0x71, 0x59, 0x8a, 0xa5, 0x09, 0xbd, 0x8a,
	0xdb, 0xa4, 0x82, 0x5d, 0xfb, 0x42, 0xc9, 0x5d, 0x7a, 0x52, 0x64, 0xe4, 0x49, 0x29, 0x9d, 0x03,
	0xb9, 0x7d, 0xcd, 0x60, 0x19, 0x0b, 0xac, 0x3c, 0x60, 0x05, 0x9f, 0x83, 0xd8, 0x09, 0xba, 0x60,
	0x5d, 0x4e, 0x28, 0x2b, 0x23, 0x4f, 0xa2, 0x9f, 0x9c, 0x11, 0x0d, 0xc0, 0x17, 0x60, 0xee, 0x54,
	0xef, 0xf4, 0x10, 0x6b, 0x6d, 0x42, 0x59, 0x1b, 0x79, 0x92, 0x1f, 0xe0, 0xa8, 0x7e, 0xe8, 0x4d,
	0x74, 0x5b, 0xc8, 0x5d, 0x0b, 0x60, 0xfd, 0xe1, 0xbe, 0xc1, 0x16, 0x58, 0xb4, 0xf4, 0xf3, 0xa6,
	0x83, 0xb0, 0xd1, 0x3c, 0x39, 0xec, 0xfa, 0x47, 0x3c, 0xa7, 0xbc, 0x1d, 0x79, 0x52, 0xd2, 0xd2,
	0xcf, 0x1b, 0x08, 0x1b, 0x1f, 0x0e, 0xbb, 0xf4, 0x74, 0x57, 0x7d, 0x67, 0x8e, 0x9d, 0xbb, 0xf3,
	0xa4, 0x98, 0x89, 0xdd, 0xbb, 0xab, 0x3c, 0x6f, 0xa3, 0x85, 0xc5, 0x41, 0x12, 0x1b, 0xb5, 0x4e,
	0xfd, 0x24, 0x51, 0x2e, 0x89, 0x86, 0x5a, 0xa7, 0xf7, 0x93, 0x4c, 0xd8, 0xb3, 0x49, 0x26, 0x80,
	0x16, 0x16, 0xe7, 0xfe, 0x9c, 0x03, 0x71, 0xba, 0x49, 0xf8, 0xb3, 0x00, 0xa2, 0xa6, 0x31, 0xee,
	0x62, 0x97, 0x9e, 0xce, 0xdf, 0x9e, 0x54, 0x3c, 0x32, 0xdd, 0xe3, 0xde, 0x61, 0xb1, 0x45, 0xac,
	0x4d, 0xa7, 0xd7, 0x45, 0xf6, 0xd9, 0xf1, 0x85, 0xb3, 0x69, 0x23, 0x8b, 0xb8, 0xe8, 0x60, 0x93,
	0xce, 0x43, 0x8b, 0x58, 0x16, 0xc1, 0xec, 0xd0, 0xab, 0x3b, 0xd7, 0x9e, 0x34, 0xcf, 0x56, 0xc6,
	0xc8, 0x93, 0x1e, 0xd1, 0xdb, 0xd0, 0x34, 0x8d, 0x70, 0x7b, 0x6f, 0xaf, 0xf2, 0x41, 0xf8, 0x97,
	0x61, 0x5e, 0xf8, 0x6d, 0x98, 0x1f, 0xf3, 0xb5, 0xa8, 0x69, 0xc0, 0x12, 0x58, 0x6c, 0x11, 0x8c,
	0x51, 0x8b, 0xf6, 0xb8, 0x69, 0x1a, 0xe3, 0x93, 0x7a, 0x46, 0x2f, 0xe8, 0x14, 0xa8, 0x72, 0x8e,
	0x1a, 0x87, 0xc0, 0x57, 0x20, 0x8e, 0x75, 0x0b, 0xb1, 0xcb, 0x9f, 0x50, 0x24, 0x3a, 0x2e, 0xf4,
	0xfb, 0xd6, 0x93, 0x92, 0xbe, 0x82, 0x7e, 0xe5, 0xee, 0xae, 0xf2, 0x2c, 0xac, 0xb1, 0xbf, 0xf0,
	0x07, 0xf0, 0x48, 0x37, 0x0c, 0x1b, 0x39, 0x8e, 0x18, 0x67, 0x43, 0xb3, 0x5c, 0x9c, 0x4c, 0x77,
	0xc9, 0x07, 0x94, 0x97, 0xe3, 0xfb, 0x1a, 0x30, 0x6f, 0x3d, 0x29, 0xe5, 0x3b, 0x8e, 0x03, 0xd4,
	0x34, 0x00, 0xb5, 0x60, 0x01, 0xbf, 0x01, 0xf3, 0x8e, 0xab, 0xbb, 0x3d, 0x47, 0x9c, 0x63, 0x03,
	0x9d, 0x9e, 0xce, 0x48, 0x83, 0xc5, 0x95, 0xf5, 0x91, 0x27, 0x8d, 0x39, 0xdc, 0xbe, 0xc6, 0x31,
	0xb8, 0x05, 0x16, 0x4c, 0xa7, 0xd9, 0x21, 0x2d, 0xbd, 0x23, 0xce, 0xcb, 0x42, 0x61, 0x41, 0x79,
	0x42, 0xcb, 0x30, 0x9d, 0x1a, 0x0d, 0x71, 0x92, 0x20, 0x08, 0xdf, 0x82, 0xb8, 0x4d, 0x3a, 0x48,
	0x7c, 0xc4, 0x72, 0xae, 0x71, 0x73, 0x49, 0x7b, 0xa5, 0x91, 0x0e, 0xf2, 0xdf, 0x12, 0x4a, 0xe3,
	0xdf, 0x12, 0x1a, 0x81, 0x06, 0x3b, 0x89, 0xe9, 0x85, 0x17, 0x17, 0x58, 0x5f, 0x9e, 0xfc, 0xc7,
	0x84, 0x2b, 0xcf, 0xc7, 0xef, 0x08, 0xaf, 0xba, 0xf7, 0x90, 0xf0, 0x30, 0xdc, 0x01, 0xa9, 0x8e,
	0xee, 0xb8, 0xcd, 0x63, 0xa4, 0xdb, 0xee, 0x21, 0xd2, 0x5d, 0x31, 0x21, 0x0b, 0x85, 0x98, 0xf2,
	0x29, 0x75, 0xa2, 0xc8, 0xfb, 0x00, 0xe0, 0x4a, 0xe4, 0xa1, 0x8d, 0xa1, 0x00, 0xc0, 0xb4, 0x9f,
	0xb0, 0x08, 0x56, 0xd4, 0xfa, 0x4e, 0xa5, 0xd9, 0xd8, 0x2d, 0xed, 0xee, 0x35, 0x9a, 0x7b, 0xea,
	0x07, 0xb5, 0xbe, 0xaf, 0xa6, 0x23, 0x99, 0xb5, 0xfe, 0x40, 0x5e, 0x9e, 0x12, 0xf7, 0xf0, 0x09,
	0x26, 0x67, 0x18, 0xbe, 0x04, 0x30, 0xcc, 0xaf, 0xab, 0xb5, 0xaa, 0x5a, 0x49, 0x0b, 0x99, 0xd5,
	0xfe, 0x40, 0x4e, 0x4f, 0xe9, 0x75, 0xdc, 0x31, 0x31, 0x9a, 0x75, 0xaf, 0xbf, 0x7b, 0xc7, 0xe8,
	0xd1, 0x59, 0xf7, 0x7a, 0xbb, 0xcd, 0xf8, 0x5f, 0x82, 0xf5, 0x30, 0xbf, 0x5c, 0x57, 0xd5, 0x4a,
	0x79, 0xb7, 0xaa, 0x7e, 0x9b, 0x8e, 0x65, 0xc4, 0xfe, 0x40, 0x5e, 0x9d, 0x4a, 0xc6, 0x67, 0x63,
	0xe2, 0xa3, 0x8d, 0xdf, 0x05, 0x30, 0xef, 0xff, 0x4e, 0xc0, 0xcf, 0xc0, 0x12, 0x33, 0xa8, 0x87,
	0xb7, 0xb2, 0xdc, 0x1f, 0xc8, 0x8b, 0x8c, 0x30, 0xd9, 0x46, 0x16, 0x24, 0x03, 0xde, 0x7e, 0x55,
	0x4d, 0x0b, 0x99, 0xc5, 0xfe, 0x40, 0x4e, 0xf8, 0x9c, 0x7d, 0x13, 0xc3, 0x1c, 0x58, 0x0c, 0xf0,
	0x5a, 0x55, 0xdd, 0x3b, 0x48, 0x47, 0x33, 0x4b, 0xfd, 0x81, 0x9c, 0xf4, 0x19, 0x35, 0x13, 0xf7,
	0xce, 0x61, 0x1e, 0xa4, 0x02, 0xce, 0x4e, 0x49, 0xa3, 0x36, 0xb1, 0x4c, 0xba, 0x3f, 0x90, 0x1f,
	0xfb, 0xa4, 0x1d, 0xdd, 0x3e, 0x33, 0xf1, 0xc6, 0x4f, 0x60, 0x21, 0xf8, 0x3d, 0x82, 0x1b, 0x60,
	0x99, 0x29, 0x4a, 0x5a, 0xf9, 0x7d, 0xa8, 0xbe, 0x95, 0xfe, 0x40, 0x5e, 0x0a, 0x48, 0x41, 0x85,
	0x41, 0x05, 0x8c, 0x7b, 0xb0, 0xfd, 0x3a, 0x2d, 0x4c, 0x2b, 0xa0, 0xbc, 0x83, 0xed, 0xd7, 0x3c,
	0xa7, 0xa4, 0x7d, 0x17, 0xae, 0x92, 0x72, 0x4a, 0xb6, 0xb5, 0xf1, 0x87, 0x00, 0x96, 0x66, 0xee,
	0x32, 0xfc, 0x1a, 0x3c, 0x65, 0xba, 0x71, 0x7f, 0x9b, 0x5a, 0xbd, 0x56, 0x09, 0xd5, 0x93, 0xe9,
	0x0f, 0xe4, 0xf5, 0x19, 0x4d, 0x50, 0xd6, 0x57, 0x40, 0xbc, 0x2f, 0x6d, 0x54, 0xb4, 0xef, 0x2b,
	0x5a, 0x5a, 0xc8, 0x3c, 0xed, 0x0f, 0xe4, 0xb5, 0x19, 0x65, 0x03, 0xd9, 0xa7, 0xc8, 0x7e, 0x58,
	0x58, 0xae, 0x55, 0x2b, 0xea, 0x6e, 0x3a, 0xfa, 0xa0, 0xb0, 0xdc, 0x31, 0x11, 0x76, 0x95, 0x77,
	0xc3, 0x7f, 0xb2, 0x91, 0xcb, 0xeb, 0xac, 0x30, 0xbc, 0xce, 0x0a, 0xbf, 0xde, 0x64, 0x23, 0x1f,
	0x6f, 0xb2, 0xc2, 0xf0, 0x26, 0x1b, 0xf9, 0xeb, 0x26, 0x1b, 0xf9, 0xb1, 0xf0, 0xbf, 0x4f, 0x6e,
	0xe8, 0x7f, 0x95, 0xc3, 0x79, 0xf6, 0x48, 0xbd, 0xfa, 0x37, 0x00, 0x00, 0xff, 0xff, 0x48, 0xa7,
	0xe2, 0x0c, 0xc1, 0x08, 0x00, 0x00,
}

func (m *NodeConfiguration) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NodeConfiguration) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NodeConfiguration) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AdditionalInfo) > 0 {
		for k := range m.AdditionalInfo {
			v := m.AdditionalInfo[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintNode(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintNode(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintNode(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Transmission != nil {
		{
			size, err := m.Transmission.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintNode(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Arch != 0 {
		i = encodeVarintNode(dAtA, i, uint64(m.Arch))
		i--
		dAtA[i] = 0x10
	}
	if m.Os != 0 {
		i = encodeVarintNode(dAtA, i, uint64(m.Os))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *NodeTransConfiguration) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NodeTransConfiguration) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NodeTransConfiguration) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MaxRecvKbps != 0 {
		i = encodeVarintNode(dAtA, i, uint64(m.MaxRecvKbps))
		i--
		dAtA[i] = 0x10
	}
	if m.MaxSendKbps != 0 {
		i = encodeVarintNode(dAtA, i, uint64(m.MaxSendKbps))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Node) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Node) MarshalTo(dAtA []byte) (int, error) {
	size := m.ProtoSize()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Node) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastHeartbeat != 0 {
		i = encodeVarintNode(dAtA, i, uint64(m.LastHeartbeat))
		i--
		dAtA[i] = 0x48
	}
	if m.Configuration != nil {
		{
			size, err := m.Configuration.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintNode(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if m.Role != 0 {
		i = encodeVarintNode(dAtA, i, uint64(m.Role))
		i--
		dAtA[i] = 0x38
	}
	if m.IsLocal {
		i--
		if m.IsLocal {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if m.Status != 0 {
		i = encodeVarintNode(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x28
	}
	{
		size, err := m.Address.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintNode(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintNode(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ConnectionId) > 0 {
		i -= len(m.ConnectionId)
		copy(dAtA[i:], m.ConnectionId)
		i = encodeVarintNode(dAtA, i, uint64(len(m.ConnectionId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.NodeId) > 0 {
		i -= len(m.NodeId)
		copy(dAtA[i:], m.NodeId)
		i = encodeVarintNode(dAtA, i, uint64(len(m.NodeId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintNode(dAtA []byte, offset int, v uint64) int {
	offset -= sovNode(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *NodeConfiguration) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Os != 0 {
		n += 1 + sovNode(uint64(m.Os))
	}
	if m.Arch != 0 {
		n += 1 + sovNode(uint64(m.Arch))
	}
	if m.Transmission != nil {
		l = m.Transmission.ProtoSize()
		n += 1 + l + sovNode(uint64(l))
	}
	if len(m.AdditionalInfo) > 0 {
		for k, v := range m.AdditionalInfo {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovNode(uint64(len(k))) + 1 + len(v) + sovNode(uint64(len(v)))
			n += mapEntrySize + 1 + sovNode(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *NodeTransConfiguration) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxSendKbps != 0 {
		n += 1 + sovNode(uint64(m.MaxSendKbps))
	}
	if m.MaxRecvKbps != 0 {
		n += 1 + sovNode(uint64(m.MaxRecvKbps))
	}
	return n
}

func (m *Node) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.NodeId)
	if l > 0 {
		n += 1 + l + sovNode(uint64(l))
	}
	l = len(m.ConnectionId)
	if l > 0 {
		n += 1 + l + sovNode(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovNode(uint64(l))
	}
	l = m.Address.ProtoSize()
	n += 1 + l + sovNode(uint64(l))
	if m.Status != 0 {
		n += 1 + sovNode(uint64(m.Status))
	}
	if m.IsLocal {
		n += 2
	}
	if m.Role != 0 {
		n += 1 + sovNode(uint64(m.Role))
	}
	if m.Configuration != nil {
		l = m.Configuration.ProtoSize()
		n += 1 + l + sovNode(uint64(l))
	}
	if m.LastHeartbeat != 0 {
		n += 1 + sovNode(uint64(m.LastHeartbeat))
	}
	return n
}

func sovNode(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozNode(x uint64) (n int) {
	return sovNode(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *NodeConfiguration) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNode
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
			return fmt.Errorf("proto: NodeConfiguration: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NodeConfiguration: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Os", wireType)
			}
			m.Os = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Os |= NodeOS(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Arch", wireType)
			}
			m.Arch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Arch |= NodeArch(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Transmission", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
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
				return ErrInvalidLengthNode
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthNode
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Transmission == nil {
				m.Transmission = &NodeTransConfiguration{}
			}
			if err := m.Transmission.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AdditionalInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
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
				return ErrInvalidLengthNode
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthNode
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.AdditionalInfo == nil {
				m.AdditionalInfo = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowNode
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
							return ErrIntOverflowNode
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
						return ErrInvalidLengthNode
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthNode
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
							return ErrIntOverflowNode
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
						return ErrInvalidLengthNode
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthNode
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipNode(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthNode
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.AdditionalInfo[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipNode(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthNode
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
func (m *NodeTransConfiguration) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNode
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
			return fmt.Errorf("proto: NodeTransConfiguration: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NodeTransConfiguration: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxSendKbps", wireType)
			}
			m.MaxSendKbps = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxSendKbps |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxRecvKbps", wireType)
			}
			m.MaxRecvKbps = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxRecvKbps |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipNode(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthNode
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
func (m *Node) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNode
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
			return fmt.Errorf("proto: Node: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Node: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
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
				return ErrInvalidLengthNode
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNode
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeId = github_com_superwhys_remoteX_pkg_common.NodeID(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
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
				return ErrInvalidLengthNode
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNode
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
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
				return ErrInvalidLengthNode
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNode
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
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
				return ErrInvalidLengthNode
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthNode
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Address.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= NodeStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsLocal", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsLocal = bool(v != 0)
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Role", wireType)
			}
			m.Role = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Role |= NodeConnectRole(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Configuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
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
				return ErrInvalidLengthNode
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthNode
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Configuration == nil {
				m.Configuration = &NodeConfiguration{}
			}
			if err := m.Configuration.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastHeartbeat", wireType)
			}
			m.LastHeartbeat = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastHeartbeat |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipNode(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthNode
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
func skipNode(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowNode
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
					return 0, ErrIntOverflowNode
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
					return 0, ErrIntOverflowNode
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
				return 0, ErrInvalidLengthNode
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupNode
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthNode
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthNode        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowNode          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupNode = fmt.Errorf("proto: unexpected end of group")
)
