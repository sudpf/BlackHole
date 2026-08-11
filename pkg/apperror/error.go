package apperror

import (
	"errors"
	"fmt"

	"github.com/mohae/deepcopy"
)

type Code int

const Success Code = 0

type Params = map[string]any

type Error struct {
	code    Code
	cause   error
	params  Params
	details any
}

func New(code Code) *Error {
	return newError(code, nil, nil, nil)
}

func NewWithParams(code Code, params Params) *Error {
	return newError(code, nil, params, nil)
}

func NewWithDetails(code Code, details any) *Error {
	return newError(code, nil, nil, details)
}

func Wrap(code Code, cause error) *Error {
	return newError(code, cause, nil, nil)
}

func WrapWithParams(code Code, cause error, params Params) *Error {
	return newError(code, cause, params, nil)
}

func WrapWithDetails(code Code, cause error, details any) *Error {
	return newError(code, cause, nil, details)
}

func newError(code Code, cause error, params Params, details any) *Error {
	if params != nil {
		params = deepcopy.Copy(params).(Params)
	}

	return &Error{
		code:    code,
		cause:   cause,
		params:  params,
		details: details,
	}
}

func As(err error) (*Error, bool) {
	var appErr *Error
	ok := errors.As(err, &appErr)
	return appErr, ok
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.cause != nil {
		return fmt.Sprintf("application error %d: %v", e.code, e.cause)
	}
	return fmt.Sprintf("application error %d", e.code)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *Error) Code() Code {
	if e == nil {
		return 0
	}
	return e.code
}

func (e *Error) Params() Params {
	if e == nil || e.params == nil {
		return nil
	}
	return deepcopy.Copy(e.params).(Params)
}

func (e *Error) Details() any {
	if e == nil {
		return nil
	}
	return e.details
}
