package errors

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

type ErrorCode int32
type Error interface {
	Error() string
	Cause() error
	Code() ErrorCode
	Format(fmt.State, rune)
	Unwrap() error
}

const (
	HandleNoProductKind ErrorCode = iota
	InvalidGrib
	ContextError
)

type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := errors.Frame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

func (s *stack) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*s)[i])
	}
	return f
}
func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

type ErrorOption func(*ErrorOptions)

func WithMessage(message string) ErrorOption {
	return func(opts *ErrorOptions) {
		opts.Message = message
	}
}

type ErrorOptions struct {
	Cause   error
	Message string
	Code    ErrorCode
}

func newErrorOptions(options ...ErrorOption) *ErrorOptions {
	errorOptions := &ErrorOptions{}

	for _, option := range options {
		option(errorOptions)
	}

	return errorOptions
}

type gribError struct {
	baseMessage string
	codeDefault ErrorCode
	stack       *stack
	options     *ErrorOptions
}

func (e gribError) Cause() error {
	if e.options.Cause != nil {
		return e.options.Cause
	}

	return nil
}

func (e gribError) Code() ErrorCode {
	if e.options != nil {
		return e.options.Code
	}

	return e.codeDefault
}

func (e gribError) Error() string {
	return e.baseMessage
}

func (e gribError) Format(st fmt.State, verb rune) {
	var message string
	if e.options != nil && e.options.Message != "" {
		message = "; " + e.options.Message
	}

	switch verb {
	case 'v':
		if st.Flag('+') {
			_, _ = fmt.Fprintf(st, "%s%s", e.Error(), message)
			if e.Cause() != nil {
				_, _ = fmt.Fprintf(st, "; %+v", e.Cause())
			} else {
				e.stack.Format(st, verb)
			}
			return
		} else {
			_, _ = fmt.Fprintf(st, "%s%s", e.Error(), message)
			if e.Cause() != nil {
				_, _ = fmt.Fprintf(st, "; %v", e.Cause())
			}
		}
	case 's', 'q':
		_, _ = fmt.Fprintf(st, "%s%s", e.Error(), message)
	}
}

func (e gribError) Unwrap() error {
	return e.Cause()
}

func newGribError(message string, codeDefault ErrorCode, options ...ErrorOption) gribError {
	return gribError{
		baseMessage: message,
		codeDefault: codeDefault,
		stack:       callers(),
		options:     newErrorOptions(options...),
	}
}

type InvalidGribError struct {
	gribError
}

func NewInavalidGribError(opts ...ErrorOption) *InvalidGribError {
	return &InvalidGribError{
		gribError: newGribError("Invalid Grib Error", InvalidGrib, opts...),
	}

}

type InvalidProductKind struct {
	gribError
}

func NewInvalidProductKind(opts ...ErrorOption) *InvalidProductKind {
	return &InvalidProductKind{gribError: newGribError("Invalid Product Kind", HandleNoProductKind, opts...)}
}

type BadContext struct {
	gribError
}

func NewContextError(opts ...ErrorOption) *BadContext {
	return &BadContext{gribError: newGribError("Context Error", ContextError, opts...)}
}
