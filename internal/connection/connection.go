package connection

import (
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type StreamReadWriter struct {
	pr           *protoutils.ProtoReader
	pw           *protoutils.ProtoWriter
	nodeId       common.NodeID
	connectionId string
}

func (rw *StreamReadWriter) ReadMessage(message proto.Message) error {
	err := rw.pr.ReadMessage(message)
	if err != nil {
		return errorutils.ErrStreamReadMessage(err)
	}

	return nil
}

func (rw *StreamReadWriter) WriteMessage(m proto.Message) error {
	err := rw.pw.WriteMessage(m)
	if err != nil {
		return errorutils.ErrStreamWriteMessage(err)
	}

	return nil
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
