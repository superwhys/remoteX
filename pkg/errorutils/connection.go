package errorutils

import (
	"fmt"
	
	"github.com/superwhys/remoteX/pkg/common"
)

type ConnectionError struct {
	*BaseError
	connectionId string
}

func ErrConnection(connId string, opts ...ErrOption) *ConnectionError {
	e := &ConnectionError{
		BaseError:    &BaseError{},
		connectionId: connId,
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

func ErrConnectToMyself(remoteId common.NodeID, connId string) *ConnectToMyselfError {
	return &ConnectToMyselfError{
		ConnectionError: ErrConnection(connId),
		remoteId:        remoteId,
	}
}

func (err *ConnectToMyselfError) Error() string {
	return fmt.Sprintf("connected to myself (%s) at %s", err.remoteId, err.connectionId)
}

type ConnectNotFoundError struct {
	*ConnectionError
}

func ErrConnectNotFound(connectionId string) *ConnectNotFoundError {
	return &ConnectNotFoundError{
		ConnectionError: &ConnectionError{connectionId: connectionId},
	}
}
