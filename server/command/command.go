package command

import (
	"context"
	"time"
	
	"github.com/gogo/protobuf/types"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/domain/command/filesystem"
	"github.com/superwhys/remoteX/pkg/errorutils"
	fsSrv "github.com/superwhys/remoteX/server/command/filesystem"
)

type ServiceImpl struct {
	fsService filesystem.Service
}

func NewCommandService() command.Service {
	return &ServiceImpl{
		fsService: fsSrv.NewFilesystemService(),
	}
}

func (s *ServiceImpl) DoCommand(ctx context.Context, cmd *command.Command) (ret *command.Ret, err error) {
	
	switch cmd.Type {
	case command.Empty:
		ret, err = s.doEmpty(cmd)
	case command.Listdir:
		ret, err = s.doListDir(cmd)
	}
	
	return
}

func (s *ServiceImpl) doEmpty(cmd *command.Command) (ret *command.Ret, err error) {
	resp := &command.MapResp{Data: map[string]string{"now_time": time.Now().Format(time.DateTime)}}
	
	anyData, err := types.MarshalAny(resp)
	if err != nil {
		return nil, err
	}
	
	ret = &command.Ret{Command: cmd, Resp: anyData}
	return
}

func (s *ServiceImpl) doListDir(cmd *command.Command) (ret *command.Ret, err error) {
	args := cmd.GetArgs()
	if len(args) < 1 {
		err = errorutils.ErrCommandMissingArguments(int32(cmd.GetType()), cmd.GetArgs())
		return
	}
	
	path, exists := args["path"]
	if !exists {
		err = errorutils.ErrCommandMissingArguments(int32(cmd.GetType()), cmd.GetArgs())
		return
	}
	
	entries, err := s.fsService.ListDir(path)
	if err != nil {
		return nil, err
	}
	
	anyData, err := types.MarshalAny(entries)
	if err != nil {
		return nil, err
	}
	
	return &command.Ret{Command: cmd, Resp: anyData}, nil
}
