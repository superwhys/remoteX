package svcutils

import (
	"context"
	"fmt"
	"sync"
	
	"github.com/pkg/errors"
	"github.com/thejerf/suture/v4"
)

type ExitStatus int

const (
	ExitSuccess            ExitStatus = 0
	ExitError              ExitStatus = 1
	ExitNoUpgradeAvailable ExitStatus = 2
	ExitRestart            ExitStatus = 3
	ExitUpgrade            ExitStatus = 4
)

func (s ExitStatus) AsInt() int {
	return int(s)
}

type ServiceWithError interface {
	suture.Service
	fmt.Stringer
	Error() error
}

// AsService wraps the given function to implement suture.Service. In addition,
// it keeps track of the returned error and allows querying that error.
func AsService(fn func(ctx context.Context) error, creator string) ServiceWithError {
	return &service{
		creator: creator,
		serve:   fn,
		mut:     sync.Mutex{},
	}
}

type service struct {
	creator string
	serve   func(ctx context.Context) error
	err     error
	mut     sync.Mutex
}

func (s *service) Serve(ctx context.Context) error {
	s.mut.Lock()
	s.err = nil
	s.mut.Unlock()
	
	err := waitContext(ctx, s.serve(ctx))
	
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

func waitContext(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s (non-context)", err.Error())
	}
	return err
}
