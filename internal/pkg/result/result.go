package result

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorTrace struct {
	caller     string
	message    string
	cause      ErrorCause
	IsExpected bool
}

type Result[T any] struct {
	errors []*ErrorTrace
	value  T
}

func NewError[T any](message string, Expected ...bool) Result[T] {
	return Err(Result[T]{}, message, Expected...)
}

func Err[T any](err Result[T], message string, Expected ...bool) Result[T] {
	_, file, line, _ := runtime.Caller(2)
	caller := fmt.Sprintf("%s:%d", file, line)

	var knErr bool
	if len(Expected) > 0 {
		knErr = Expected[0]
	}

	_error := &ErrorTrace{caller, message, INTERNAL_SERVICE_ERROR, knErr}
	if err.errors == nil {
		err.errors = make([]*ErrorTrace, 1, 4)
		err.errors[0] = _error
	} else {
		err.errors = append(err.errors, _error)
	}

	return Result[T]{errors: err.errors}
}

func Ok[T any](v T) Result[T] {
	return Result[T]{value: v, errors: nil}
}

func (r Result[T]) Value() T {
	return r.value
}

func (r Result[T]) IsError() bool {
	return len(r.errors) > 0
}

func (r Result[T]) RootError() *ErrorTrace {
	if len(r.errors) > 0 {
		return r.errors[0]
	}

	return nil
}

func (r Result[T]) LastError() *ErrorTrace {
	if len(r.errors) > 0 {
		return r.errors[len(r.errors)-1]
	}

	return nil
}

func (r Result[T]) Error() string {
	if len(r.errors) <= 0 {
		return ""
	}

	errorString := strings.Builder{}
	for _, v := range r.errors {
		errorString.WriteString(fmt.Sprintf("%s:\n\t%s\n", v.caller, v.message))
	}

	return errorString.String()
}

func (r Result[T]) ExpectedError() *ErrorTrace {
	for _, err := range r.errors {
		if err.IsExpected {
			return err
		}
	}

	return nil
}

func (r Result[T]) WithCause(cause ErrorCause) Result[T] {
	r.errors[0].cause = cause
	return r
}

func (e ErrorTrace) Error() string {
	return e.message
}

func (e ErrorTrace) Cause() ErrorCause {
	return e.cause
}
