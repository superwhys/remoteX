package filesystem

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/domain/command"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/pkg/protoutils"

	fsDomain "github.com/superwhys/remoteX/domain/command/filesystem"
)

type ServiceImpl struct {
	fs filesystem.FileSystem
}

func NewFilesystemService() fsDomain.Service {
	return &ServiceImpl{
		fs: filesystem.NewBasicFileSystem(),
	}
}

func (s *ServiceImpl) Name() string {
	return "filesystem"
}

func (s *ServiceImpl) SupportedCommands() []command.CommandType {
	return []command.CommandType{
		command.Listdir,
	}
}

func (s *ServiceImpl) Invoke(ctx context.Context, cmd *command.Command, cmdCtx *command.CommandContext) (proto.Message, error) {
	switch cmd.Type {
	case command.Listdir:
		if !cmdCtx.IsRemote {
			return s.ListDir(ctx, cmd.GetArgs())
		}

		err := cmdCtx.Remote.WriteMessage(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "writeCommand")
		}

		resp := new(command.Ret)
		if err = cmdCtx.Remote.ReadMessage(resp); err != nil {
			return nil, errors.Wrap(err, "readRespMessage")
		}

		listResp, err := protoutils.DecodeAnyProto(resp.GetResp())
		if err != nil {
			return nil, errors.Wrap(err, "decodeAnyProto")
		}

		return listResp, nil
	default:
		return nil, errors.New("unsupport command")
	}
}

func (s *ServiceImpl) ListDir(ctx context.Context, args command.Args) (proto.Message, error) {
	path, exists := args.GetArg("path")
	if !exists {
		return nil, errorutils.ErrCommandMissingArguments(int32(command.Listdir), args)
	}

	entries, err := s.fs.List(path.GetStrValue())
	if err != nil {
		return nil, err
	}

	return &filesystem.ListResp{
		Entries: entries,
	}, nil
}
