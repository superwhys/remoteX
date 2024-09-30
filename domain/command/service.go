package command

import "context"

type Service interface {
	DoCommand(ctx context.Context, cmd *Command) (*Ret, error)
}

type ServiceImpl struct{}

func NewCommandService() Service {
	return &ServiceImpl{}
}

func (s *ServiceImpl) DoCommand(ctx context.Context, cmd *Command) (*Ret, error) {
	return &Ret{Command: cmd, Resp: map[string]string{"data": "success"}}, nil
}
