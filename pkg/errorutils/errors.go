package errorutils

import (
	"fmt"
	"strings"
)

type BaseError struct {
	error
	Code     uint64
	Messages []string
}

type ErrOption func(connectionError *BaseError)

func WithMsg(message string) ErrOption {
	return func(baseError *BaseError) {
		if baseError.Messages == nil {
			baseError.Messages = make([]string, 0)
		}
		baseError.Messages = append(baseError.Messages, message)
	}
}

func WithCode(code uint64) ErrOption {
	return func(baseError *BaseError) {
		baseError.Code = code
	}
}

func WithError(err error) ErrOption {
	return func(baseError *BaseError) {
		baseError.SetError(err)
	}
}

func PackBaseError(opts ...ErrOption) *BaseError {
	e := new(BaseError)
	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *BaseError) SetError(err error) {
	e.error = err
}

func (e *BaseError) Error() string {
	final := make([]string, 0, 3)
	if e.Code != 0 {
		final = append(final, fmt.Sprintf("Code: %d.", e.Code))
	}
	if len(e.Messages) != 0 {
		final = append(final, strings.Join(e.Messages, ";"))
	}

	if e.error != nil {
		final = append(final, fmt.Sprintf("Error: %v.", e.error))
	}

	return strings.Join(final, " ")
}

func (e *BaseError) As(target any) bool {
	if target == nil {
		return false
	}

	switch t := target.(type) {
	case *BaseError:
		*t = *e
		return true
	default:
		return false
	}
}
