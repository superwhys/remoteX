package command

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/server/command/filesystem"
	"github.com/superwhys/remoteX/server/command/sync"
	"github.com/superwhys/remoteX/server/command/tunnel"
)

type ServiceImpl struct {
	strategy map[command.CommandType]command.CommandProvider
}

func NewCommandService() command.Service {
	s := &ServiceImpl{
		strategy: map[command.CommandType]command.CommandProvider{},
	}

	fsSrv := filesystem.NewFilesystemService()
	syncSrv := sync.NewSyncService()
	tunnelSrv := tunnel.NewTunnelService()
	// screenshotSrv := screenshot.NewScreenshotService()

	s.registerStrategy(s, command.Empty)
	s.registerStrategy(fsSrv, command.Listdir)
	s.registerStrategy(syncSrv, command.Push, command.Pull)
	s.registerStrategy(tunnelSrv, command.Forward, command.Forwardreceive, command.Listtunnel, command.Closetunnel)
	// s.registerStrategy(command.Screenshot, screenshotSrv.Screenshot)

	return s
}

func (s *ServiceImpl) registerStrategy(handler command.CommandProvider, cmdType ...command.CommandType) {
	for _, ct := range cmdType {
		s.strategy[ct] = handler
	}
}

func (s *ServiceImpl) handleCommand(ctx context.Context, cmd *command.Command, opt *command.RemoteOpt) (proto.Message, error) {
	handler, ok := s.strategy[cmd.Type]
	if !ok {
		return nil, fmt.Errorf("unknown command type: %s", cmd.Type)
	}

	return handler.Invoke(ctx, cmd, opt)
}

func (s *ServiceImpl) packAnyData(pm proto.Message) (*types.Any, error) {
	var anyData *types.Any
	if pm == nil {
		return nil, nil
	}

	anyData, err := types.MarshalAny(pm)
	if err != nil {
		return nil, err
	}

	return anyData, nil
}

func (s *ServiceImpl) DoOriginCommand(ctx context.Context, cmd *command.Command, conn connection.StreamConnection) (*command.Ret, error) {
	pm, err := s.handleCommand(ctx, cmd, &command.RemoteOpt{Conn: conn})
	if err != nil {
		return nil, err
	}

	anyData, err := s.packAnyData(pm)
	if err != nil {
		return nil, err
	}

	return &command.Ret{Resp: anyData}, nil
}

func (s *ServiceImpl) DoAcceptCommand(ctx context.Context, cmd *command.Command, stream connection.Stream) (*command.Ret, error) {
	pm, err := s.handleCommand(ctx, cmd, &command.RemoteOpt{Stream: stream})
	if err != nil {
		return nil, err
	}

	anyData, err := s.packAnyData(pm)
	if err != nil {
		return nil, err
	}

	return &command.Ret{Resp: anyData}, nil
}

func (s *ServiceImpl) DoLocalCommand(ctx context.Context, cmd *command.Command) (ret *command.Ret, err error) {
	pm, err := s.handleCommand(ctx, cmd, nil)
	if err != nil {
		return nil, err
	}

	anyData, err := s.packAnyData(pm)
	if err != nil {
		return nil, err
	}

	return &command.Ret{Resp: anyData}, nil
}

func (s *ServiceImpl) Invoke(_ context.Context, _ *command.Command, opt *command.RemoteOpt) (proto.Message, error) {
	resp := &command.MapResp{Data: map[string]string{"now_time": time.Now().Format(time.DateTime)}}

	return resp, nil
}
