package command

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/server/command/filesystem"
)

type strategyHandler func(ctx context.Context, args command.Args) (proto.Message, error)

type ServiceImpl struct {
	strategy map[command.CommandType]strategyHandler
}

func NewCommandService() command.Service {
	s := &ServiceImpl{
		strategy: map[command.CommandType]strategyHandler{},
	}

	fsSrv := filesystem.NewFilesystemService()

	s.registerStrategy(command.Empty, s.doEmpty)
	s.registerStrategy(command.Listdir, fsSrv.ListDir)

	return s
}

func (s *ServiceImpl) registerStrategy(cmdType command.CommandType, handler strategyHandler) {
	s.strategy[cmdType] = handler
}

func (s *ServiceImpl) handleCommand(ctx context.Context, cmdType command.CommandType, args command.Args) (proto.Message, error) {
	handler, ok := s.strategy[cmdType]
	if !ok {
		return nil, fmt.Errorf("unknown command type: %s", cmdType)
	}

	return handler(ctx, args)
}

func (s *ServiceImpl) DoCommand(ctx context.Context, cmd *command.Command) (ret *command.Ret, err error) {
	pm, err := s.handleCommand(ctx, cmd.GetType(), cmd.GetArgs())
	if err != nil {
		return nil, err
	}

	anyData, err := types.MarshalAny(pm)
	if err != nil {
		return nil, err
	}

	return &command.Ret{Command: cmd, Resp: anyData}, nil
}

func (s *ServiceImpl) doEmpty(_ context.Context, _ command.Args) (proto.Message, error) {
	resp := &command.MapResp{Data: map[string]string{"now_time": time.Now().Format(time.DateTime)}}

	return resp, nil
}
