package errors

import "errors"

// Unwrap wrapper function for errors.Unwrap.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is wrapper function for errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As wrapper function for errors.As.
func As(err error, target any) bool {
	return errors.As(err, target)
}
