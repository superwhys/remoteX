package command

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/superwhys/remoteX/domain/command/filesystem"
	"github.com/superwhys/remoteX/pkg/errorutils"
)

type Service interface {
	DoCommand(ctx context.Context, cmd *Command) (*Ret, error)
}

type ServiceImpl struct {
	fsService filesystem.Service
}

func NewCommandService() Service {
	return &ServiceImpl{
		fsService: filesystem.NewFilesystemService(),
	}
}

func (s *ServiceImpl) DoCommand(ctx context.Context, cmd *Command) (ret *Ret, err error) {

	switch cmd.Type {
	case Empty:
		ret, err = s.doEmpty(cmd)
	case Listdir:
		ret, err = s.doListDir(cmd)
	}

	return
}

func (s *ServiceImpl) doEmpty(cmd *Command) (ret *Ret, err error) {
	resp := &MapResp{Data: map[string]string{"now_time": time.Now().Format(time.DateTime)}}

	anyData, err := types.MarshalAny(resp)
	if err != nil {
		return nil, err
	}

	ret = &Ret{Command: cmd, Resp: anyData}
	return
}

func (s *ServiceImpl) doListDir(cmd *Command) (ret *Ret, err error) {
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

	return &Ret{Command: cmd, Resp: anyData}, nil
}
