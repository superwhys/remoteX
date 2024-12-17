package command

import (
	"context"

	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type Args map[string]Command_Arg

func (a Args) GetArg(name string) (Command_Arg, bool) {
	arg, exists := a[name]
	if !exists {
		return Command_Arg{}, false
	}

	return arg, true
}

type RemoteOpt struct {
	Conn   connection.StreamConnection
	Stream connection.Stream
}

func (ro *RemoteOpt) IsOrigin() bool {
	return ro.Conn != nil
}

func (ro *RemoteOpt) GetRemoteChannel() (stream protoutils.ProtoMessageReadWriter, err error) {
	stream = ro.Stream
	if ro.Conn != nil {
		stream, err = ro.Conn.OpenStream()
	}

	return
}

type Service interface {
	RegisterProvider(provider CommandProvider)
	Execute(ctx context.Context, cmd *Command, cmdCtx *CommandContext) (*Ret, error)
}
