package protoutils

import (
	"bytes"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/pkg/protocol"
)

func TestProtoWriter_WriteMessage_Success(t *testing.T) {
	message := &protocol.Address{IpAddress: "127.0.0.1", Port: 80}
	
	buf := new(bytes.Buffer)
	
	protoWriter := NewProtoWriter(buf)
	protoReader := NewProtoReader(buf)
	
	err := protoWriter.WriteMessage(message)
	if err != nil {
		t.Fatalf("WriteMessage failed: %v", err)
	}
	
	expectedMessage := new(protocol.Address)
	if err := protoReader.ReadMessage(expectedMessage); err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}
	
	assert.Equal(t, message, expectedMessage)
}
