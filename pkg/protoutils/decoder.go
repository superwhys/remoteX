// File:		decoder.go
// Created by:	Hoven
// Created on:	2024-10-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package protoutils

import (
	"errors"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
)

type DecoderFunc func([]byte) (proto.Message, error)

var (
	typeDecoders = map[string]DecoderFunc{}
)

func RegisterDecoderFunc(name string, decoderFunc DecoderFunc) {
	typeDecoders[name] = decoderFunc
}

func DecodeAnyProto(anyProto *types.Any) (proto.Message, error) {
	decoder, exists := typeDecoders[anyProto.TypeUrl]
	if !exists {
		return nil, errors.New("unknown any.proto decoder")
	}

	return decoder(anyProto.Value)
}