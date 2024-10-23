// File:		init.go
// Created by:	Hoven
// Created on:	2024-10-22
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pb

import (
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

func init() {
	protoutils.RegisterDecoderFunc("type.googleapis.com/pb.SyncResp", func(b []byte) (proto.Message, error) {
		resp := new(SyncResp)
		err := proto.Unmarshal(b, resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})
}
