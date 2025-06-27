package domain

import (
	"errors"
	"fmt"
)

// The domain.Error struct represents an error that occurred within the domain.
// It is used to provide a standardized error structure for the domain layer.
// These errors can then be used by interfaces to understand the nature of the error
// and produce a meaningful response to the user.
// These are specially helpful for the gateway interface.

// Error contains the standard error definition for this domain.
type Error struct {
	errorType string
	cause     error
	code      string
	message   string
}

var (
	// ErrInternal is the generic error for internal errors.
	ErrInternal = NewError("internal")

	// ErrFatal implies that the operation failed and cannot be retried.
	ErrFatal = NewError("fatal")

	// ErrSkip is not technically an error, but it is used to skip this task.
	ErrSkip = NewError("skip")

	// ErrAccessDenied is an error that appears when some business method
	// tries to access data that belong to other company
	ErrAccessDenied = NewError("access-denied")

	// ErrInvalid implies that something is wrong with the data, but we're not exactly sure what,
	// only that its not a network issue and retrying will probably cause the same issue.
	ErrInvalid = NewError("invalid")

	// ErrQueue is an error that appears when we need to queue the task.
	ErrQueue = NewError("queue")
)

// NewError instantiates a new error.
func NewError(errorType string) *Error {
	return &Error{errorType: errorType}
}

func (e *Error) copy() *Error {
	ne := new(Error)
	*ne = *e
	return ne
}

// WithCause attaches any error instance to the Error.
func (e *Error) WithCause(cause error) *Error {
	ne := e.copy()
	ne.cause = cause
	return ne
}

// WithMessage adds a message to the Error.
func (e *Error) WithMessage(message string) *Error {
	ne := e.copy()
	ne.message = message
	return ne
}

// WithCode adds a code to the Error.
func (e *Error) WithCode(code string) *Error {
	ne := e.copy()
	ne.code = code
	return ne
}

// WithMsgf adds the message with formatting details.
func (e *Error) WithMsgf(message string, args ...any) *Error {
	return e.WithMessage(fmt.Sprintf(message, args...))
}

// Error provides the string representation of the error.
func (e *Error) Error() string {
	if e.message == "" {
		if e.cause == nil {
			return e.errorType
		}
		return e.cause.Error()
	}
	if e.cause == nil {
		return e.message
	}
	return fmt.Sprintf("%s (%s)", e.message, e.cause.Error())
}

// Code provides the code representation of the error.
func (e *Error) Code() string {
	return e.code
}

// Is checks to see if the target error matches the current error or
// part of the chain.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return errors.Is(e.cause, target)
	}
	return e.errorType == t.errorType
}

// Cause returns the error that caused this error.
func (e *Error) Cause() error {
	return e.cause
}

// Message returns just the message component, if present
func (e *Error) Message() string {
	return e.message
}
