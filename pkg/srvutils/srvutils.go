package srvutils

import (
	"context"
	"errors"
	"fmt"
	"sync"
	
	"github.com/thejerf/suture/v4"
)

type Service interface {
	suture.Service
	fmt.Stringer
	
	Error() error
}

type service struct {
	creator string
	serve   func(context.Context) error
	err     error
	mut     sync.Mutex
}

func AsService(fn func(ctx context.Context) error, creator string) Service {
	return &service{
		creator: creator,
		serve:   fn,
		mut:     sync.Mutex{},
	}
}

func waitContext(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%s (non-context)", err.Error())
		}
		
		return err
	}
}

func (s *service) Serve(ctx context.Context) error {
	s.mut.Lock()
	s.err = nil
	s.mut.Unlock()
	
	// err := waitContext(ctx, s.serve(ctx))
	err := s.serve(ctx)
	
	s.mut.Lock()
	s.err = err
	s.mut.Unlock()
	
	return err
}

func (s *service) Error() error {
	s.mut.Lock()
	defer s.mut.Unlock()
	return s.err
}

func (s *service) String() string {
	return fmt.Sprintf("Service@%p created by %v", s, s.creator)
}
