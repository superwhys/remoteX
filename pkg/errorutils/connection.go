// File:		connection.go
// Created by:	Hoven
// Created on:	2024-09-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package errorutils

import (
	"fmt"
	
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/common"
)

type ConnectionError struct {
	*BaseError
	connectionId string
}

func ErrConnection(conn connection.TlsConn, opts ...ErrOption) *ConnectionError {
	e := &ConnectionError{
		BaseError:    &BaseError{},
		connectionId: conn.String(),
	}
	
	for _, opt := range opts {
		opt(e.BaseError)
	}
	
	return e
}

func (ce *ConnectionError) Error() string {
	return ce.String()
}

func (ce *ConnectionError) String() string {
	return fmt.Sprintf("ConnectionError{connectionId:%s}. Error{%v}", ce.connectionId, ce.BaseError)
}

type ConnectToMyselfError struct {
	*ConnectionError
	remoteId common.NodeID
}

func ErrConnectToMyself(remoteId common.NodeID, conn connection.TlsConn) *ConnectToMyselfError {
	return &ConnectToMyselfError{
		ConnectionError: ErrConnection(conn),
		remoteId:        remoteId,
	}
}

func (err *ConnectToMyselfError) Error() string {
	return fmt.Sprintf("connected to myself (%s) at %s", err.remoteId, err.connectionId)
}
