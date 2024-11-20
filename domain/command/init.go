// File:		init.go
// Created by:	Hoven
// Created on:	2024-11-20
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package command

import (
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

func init() {
	protoutils.RegisterDecoderFunc("type.googleapis.com/command.ListTunnelResp", func(b []byte) (proto.Message, error) {
		resp := new(ListTunnelResp)
		err := proto.Unmarshal(b, resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})

	protoutils.RegisterDecoderFunc("type.googleapis.com/command.Tunnel", func(b []byte) (proto.Message, error) {
		resp := new(Tunnel)
		err := proto.Unmarshal(b, resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})

}
