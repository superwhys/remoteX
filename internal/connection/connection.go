package connection

import (
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type StreamReadWriter struct {
	pr           *protoutils.ProtoReader
	pw           *protoutils.ProtoWriter
	nodeId       common.NodeID
	connectionId string
}

func (rw *StreamReadWriter) ReadMessage(message proto.Message) error {
	return rw.pr.ReadMessage(message)
}

func (rw *StreamReadWriter) WriteMessage(m proto.Message) error {
	return rw.pw.WriteMessage(m)
}

func (rw *StreamReadWriter) GetConnectionId() string {
	return rw.connectionId
}

func (rw *StreamReadWriter) GetNodeId() common.NodeID {
	return rw.nodeId
}

func (rw *StreamReadWriter) SetNodeId(nodeId common.NodeID) {
	rw.nodeId = nodeId
}
