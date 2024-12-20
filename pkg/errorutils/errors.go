package errorutils

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/pkg/common"
)

const (
	DefaultErrorErrCode = 10000400
)

type RemoteXError struct {
	Cause   error
	ErrCode uint64
	ErrMsg  string
}

func (re *RemoteXError) Code() uint64 {
	return re.ErrCode
}

func (re *RemoteXError) Unwrap() error {
	return re.Cause
}

func (re *RemoteXError) Error() string {
	var ErrMsgs []string
	var originalCause error

	current := re
	for current != nil {
		if current.ErrMsg != "" {
			ErrMsgs = append(ErrMsgs, current.ErrMsg)
		}

		if current.Cause == nil {
			break
		}

		nextError := new(RemoteXError)
		if errors.As(current.Cause, &nextError) {
			current = nextError
			originalCause = nextError.Cause
		} else {
			originalCause = current.Cause
			break
		}
	}

	slices.Reverse(ErrMsgs)

	if originalCause == nil {
		return fmt.Sprintf("RemoteXError(ErrCode=%d, ErrMsg=%s)",
			re.ErrCode,
			strings.Join(ErrMsgs, ": "))
	}

	return fmt.Sprintf("RemoteXError(ErrCode=%d, Cause=%v, ErrMsg=%s)",
		re.ErrCode,
		originalCause,
		strings.Join(ErrMsgs, ": "))
}

func WrapRemoteXError(err error, msg string) *RemoteXError {
	return &RemoteXError{
		Cause:   err,
		ErrCode: DefaultErrorErrCode,
		ErrMsg:  msg,
	}
}

func WithRemoteXError(ErrCode uint64, err error, msg string) *RemoteXError {
	return &RemoteXError{
		Cause:   err,
		ErrCode: ErrCode,
		ErrMsg:  msg,
	}
}

func WithRemoteXErrorPackerErr(ErrCode uint64, msg string) func(err error) *RemoteXError {
	return func(err error) *RemoteXError {
		return WithRemoteXError(ErrCode, errors.Wrap(err, msg), msg)
	}
}

func WithRemoteXErrorPackerMsg(ErrCode uint64, msg string) func(packMsg string) *RemoteXError {
	return func(packMsg string) *RemoteXError {
		msg = fmt.Sprintf("%s: %s", msg, packMsg)
		return WithRemoteXError(ErrCode, nil, msg)
	}
}

func WithRemoteXErrorPackerCommand(ErrCode uint64, msg string) func(cmd string, msgs ...string) *RemoteXError {
	return func(cmd string, msgs ...string) *RemoteXError {
		msg = fmt.Sprintf("%s: %s: %s", msg, cmd, strings.Join(msgs, ", "))
		return WithRemoteXError(ErrCode, nil, msg)
	}
}

func WithRemoteXErrorPackerNode(ErrCode uint64, msg string) func(nodeId common.NodeID) *RemoteXError {
	return func(nodeId common.NodeID) *RemoteXError {
		if nodeId != "" {
			msg = fmt.Sprintf("nodeId: %s. %s", nodeId, msg)
		}
		return WithRemoteXError(ErrCode, nil, msg)
	}
}

func WithRemoteXErrorPackerFilesystem(ErrCode uint64, msg string) func(path string, err error) *RemoteXError {
	return func(path string, err error) *RemoteXError {
		msg = fmt.Sprintf("Path(%s) %s", path, msg)
		return WithRemoteXError(ErrCode, err, msg)
	}
}

// inner error
var (
	// Connection errors ErrCode: 1000100 - 1000199
	ErrConnectionCert       = WithRemoteXError(1000100, nil, "Connection certificate invalidate")
	ErrConnectToMyself      = WithRemoteXError(1000101, nil, "Connected to myself")
	ErrConnectNotFound      = WithRemoteXError(1000102, nil, "Connection not found")
	ErrConnectionRemoteDead = WithRemoteXError(1000103, nil, "Remote connection is dead")
	ErrGetListenerCreator   = WithRemoteXErrorPackerErr(1000104, "Get listener creator error")
	ErrGetDialerCreator     = WithRemoteXErrorPackerErr(1000105, "Get dialer creator error")
	ErrEstablishConnection  = WithRemoteXErrorPackerErr(1000106, "Establish connection error")
	ErrStreamWriteMessage   = WithRemoteXErrorPackerErr(1000107, "Stream write message error")
	ErrStreamReadMessage    = WithRemoteXErrorPackerErr(1000108, "Stream read message error")
	ErrGetHeartbeatStream   = WithRemoteXErrorPackerErr(10000109, "Get heartbeat stream error")

	// Command errors ErrCode: 10000200 - 10000299
	ErrCommandTypeNotSupport   = WithRemoteXErrorPackerCommand(10000200, "Command type not support")
	ErrCommandArgsError        = WithRemoteXErrorPackerCommand(10000201, "Command args error")
	ErrCommandMissingArguments = WithRemoteXErrorPackerCommand(10000202, "Missing arguments")
	ErrInvokeProvider          = WithRemoteXErrorPackerCommand(10000203, "Invoke provider error")
	ErrHandleCommand           = WithRemoteXErrorPackerCommand(10000204, "Command handle error")
	ErrReadCommand             = WithRemoteXErrorPackerErr(10000205, "Read command error")

	// Node errors ErrCode: 10000300 - 10000399
	ErrSameOnlineNode      = WithRemoteXErrorPackerNode(10000300, "Same online node error")
	ErrNodeMissingHostPort = WithRemoteXErrorPackerNode(10000300, "Node missing host port error")
	ErrNodeNotFound        = WithRemoteXErrorPackerNode(10000301, "Node not found error")

	// Filesystem errors ErrCode: 20000100-20000199
	ErrFilePathNotFound = WithRemoteXErrorPackerFilesystem(20000100, "File path not found")
	ErrDirPathNotFound  = WithRemoteXErrorPackerFilesystem(20000101, "Directory path not found")
	ErrDirPathNotAbs    = WithRemoteXErrorPackerFilesystem(20000102, "Directory path must be absolute")
	ErrReadDir          = WithRemoteXErrorPackerFilesystem(20000103, "Read directory error")
	ErrCreateDir        = WithRemoteXErrorPackerFilesystem(20000104, "Create directory error")
	ErrRemoveDir        = WithRemoteXErrorPackerFilesystem(20000105, "Remove directory error")
	ErrListDir          = WithRemoteXErrorPackerFilesystem(20000106, "List directory error")
	ErrOpenFile         = WithRemoteXErrorPackerFilesystem(20000107, "Open file error")
	ErrCreateFile       = WithRemoteXErrorPackerFilesystem(20000108, "Create file error")
	ErrGetInfo          = WithRemoteXErrorPackerFilesystem(20000109, "Get entry info error")
	ErrStat             = WithRemoteXErrorPackerFilesystem(20000110, "Stat error")
	ErrLstat            = WithRemoteXErrorPackerFilesystem(20000111, "Lstat error")
	ErrReadLink         = WithRemoteXErrorPackerFilesystem(20000112, "Read entry link error")
)

func IsRemoteDead(err error) bool {
	return errors.Is(err, ErrConnectionRemoteDead)
}
