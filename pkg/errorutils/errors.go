package errorutils

import (
	"fmt"
	"strings"
)

type BaseError struct {
	error
	Code    int
	Message string
}

type ErrOption func(connectionError *BaseError)

func WithMsg(message string) ErrOption {
	return func(connectionError *BaseError) {
		connectionError.Message = message
	}
}

func WithError(err error) ErrOption {
	return func(connectionError *BaseError) {
		connectionError.SetError(err)
	}
}

func (e *BaseError) SetError(err error) {
	e.error = err
}

func (e *BaseError) Error() string {
	final := make([]string, 0, 3)
	if e.Code != 0 {
		final = append(final, fmt.Sprintf("Code: %d.", e.Code))
	}
	if e.Message != "" {
		final = append(final, e.Message)
	}
	
	if e.error != nil {
		final = append(final, fmt.Sprintf("Error: %v.", e.error))
	}
	
	return strings.Join(final, " ")
}

func PackBaseError(args ...any) *BaseError {
	e := new(BaseError)
	for _, arg := range args {
		switch t := arg.(type) {
		case error:
			e.error = t
		case int:
			e.Code = t
		case string:
			e.Message = t
		}
	}
	
	return e
}
