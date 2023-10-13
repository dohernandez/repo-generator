package errors

import (
	"fmt"
)

type errorString struct {
	message string
}

// Error implements the standard library error interface.
func (s *errorString) Error() string {
	return s.message
}

// New returns an error with the supplied message without cause.
func New(message string) error {
	return &errorString{
		message: message,
	}
}

// Newf returns an error without cause with the formats according to a format specifier.
func Newf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)

	return &errorString{
		message: message,
	}
}

// Is implements future error.Is functionality.
// An Error is equivalent if err message identical.
func (s *errorString) Is(err error) bool {
	return s.message == err.Error()
}

type withMessage struct {
	message string
	err     error
}

// Error implements the standard library error interface.
func (wm *withMessage) Error() string {
	return wm.message
}

// Unwrap implements errors.Unwrap for Error.
func (wm *withMessage) Unwrap() error {
	return wm.err
}

// Wrap returns an error annotating
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf("%s: %s", message, err)

	return &withMessage{
		// message is the full concatenate error message (top to bottom)
		message: msg,
		// err is the original error
		err: err,
	}
}

// Wrapf returns an error annotating
// at the point Wrapf is called, and the supplied message.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, args...)

	return Wrap(err, message)
}

type withError struct {
	// message is the full concatenate error message (top to bottom)
	message string
	// err is the supplied error most of the time the sentinel error.
	err error
	// cause is the original error.
	cause error
}

// Error implements the standard library error interface.
func (we *withError) Error() string {
	return we.message
}

// Unwrap implements errors.Unwrap for Error.
func (we *withError) Unwrap() error {
	return we.err
}

// Cause returns the underlying cause of error.
func (we *withError) Cause() error {
	return we.cause
}

// WrapWithError returns an error annotating err with cause
// at the point WrapWithError is called, and the supplied err.
//
// If err is nil, WrapWithError returns supplied err.
// If supplied err is nil, WrapWithError returns err.
func WrapWithError(err error, supplied error) error {
	if err == nil {
		return supplied
	}

	if supplied == nil {
		return err
	}

	msg := fmt.Sprintf("%s: %s", supplied, err)

	return &withError{
		message: msg,
		err:     supplied,
		cause:   err,
	}
}

// Is implements future error.Is functionality.
// An Error is equivalent if err message or any of the underlying cause message are identical.
func (we *withError) Is(target error) bool {
	if Is(we.err, target) {
		return true
	}

	cause := Cause(we)
	if cause == nil {
		return false
	}

	return Is(cause, target)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// If the error does not implement Cause, the error nil will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	//nolint:errorlint
	cause, ok := err.(causer)
	if !ok {
		return nil
	}

	return cause.Cause()
}
