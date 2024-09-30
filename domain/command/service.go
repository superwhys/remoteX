package command

import "context"

type Service interface {
	DoCommand(ctx context.Context, cmd *Command) (any, error)
}

type ServiceImpl struct{}

func NewCommandService() Service {
	return &ServiceImpl{}
}

func (s *ServiceImpl) DoCommand(ctx context.Context, cmd *Command) (any, error) {
	return nil, nil
}
