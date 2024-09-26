package common

import (
	"crypto/sha256"
	"encoding/hex"
)

type NodeID string

func (n NodeID) String() string {
	return string(n)
}

func NewNodeID(rawCert []byte) NodeID {
	hash := sha256.New()
	hash.Write(rawCert)

	return NodeID(hex.EncodeToString(hash.Sum(nil)))
}
