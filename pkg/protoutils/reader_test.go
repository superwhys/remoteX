package protoutils

import (
	"bytes"
	"encoding/binary"
	"testing"
	
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/pkg/protocol"
)

func TestProtoReader_ReadMessage_Success(t *testing.T) {
	expectedMessage := &protocol.Address{IpAddress: "127.0.0.1", Port: 80}
	msgBytes, err := proto.Marshal(expectedMessage)
	if err != nil {
		t.Fatalf("failed to marshal proto message: %v", err)
	}
	
	msgLength := uint32(len(msgBytes))
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, msgLength)
	buf.Write(msgBytes)
	
	protoReader := NewProtoReader(buf)
	actualMessage := &protocol.Address{}
	if err := protoReader.ReadMessage(actualMessage); err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}
	
	assert.Equal(t, expectedMessage, actualMessage)
}
