package command

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/server/command/filesystem"
	"github.com/superwhys/remoteX/server/command/sync"
	"github.com/superwhys/remoteX/server/command/tunnel"
)

type ServiceImpl struct {
	providers map[command.CommandType]command.CommandProvider
}

func NewCommandService() command.Service {
	s := &ServiceImpl{
		providers: make(map[command.CommandType]command.CommandProvider),
	}

	fsSrv := filesystem.NewFilesystemService()
	syncSrv := sync.NewSyncService()
	tunnelSrv := tunnel.NewTunnelService()

	s.RegisterProvider(s)
	s.RegisterProvider(fsSrv)
	s.RegisterProvider(syncSrv)
	s.RegisterProvider(tunnelSrv)

	return s
}

func (s *ServiceImpl) Name() string {
	return "defaultCommand"
}

func (s *ServiceImpl) SupportedCommands() []command.CommandType {
	return []command.CommandType{command.Empty}
}

func (s *ServiceImpl) Invoke(ctx context.Context, cmd *command.Command, cmdCtx *command.CommandContext) (proto.Message, error) {
	resp := &command.MapResp{
		Data: map[string]string{
			"now_time": time.Now().Format(time.DateTime),
		},
	}
	return resp, nil
}

func (s *ServiceImpl) RegisterProvider(provider command.CommandProvider) {
	if provider == nil {
		return
	}

	for _, t := range provider.SupportedCommands() {
		s.providers[t] = provider
	}
}

func (s *ServiceImpl) Execute(ctx context.Context, cmd *command.Command, cmdCtx *command.CommandContext) (*command.Ret, error) {
	provider, ok := s.providers[cmd.Type]
	if !ok {
		return nil, fmt.Errorf("unknown command type: %s", cmd.Type)
	}

	resp, err := provider.Invoke(ctx, cmd, cmdCtx)
	if err != nil {
		return nil, err
	}

	anyData, err := types.MarshalAny(resp)
	if err != nil {
		return nil, err
	}

	return &command.Ret{
		Resp: anyData,
	}, nil
}
