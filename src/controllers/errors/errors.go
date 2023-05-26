package controller_errors

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	NoType = ErrorType(iota)
	BadRequest
	NotFound
	TooManyRequests
	NotImplemented
)

type ErrorType uint

type customError struct {
	errorType     ErrorType
	originalError error
}

func (error customError) Error() string {
	return error.originalError.Error()
}

func (t ErrorType) New(msg string) error {
	return t.Newf(msg)
}

func (t ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)
	return customError{errorType: t, originalError: err}
}

func (t ErrorType) Wrap(err error, msg string) error {
	return t.Wrapf(err, msg)
}

func (t ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)
	return customError{errorType: t, originalError: newErr}
}

func New(msg string) error {
	return customError{errorType: NoType, originalError: errors.New(msg)}
}

func Newf(msg string, args ...interface{}) error {
	return customError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
		}
	}

	return customError{errorType: NoType, originalError: wrappedError}
}

func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}

	return NoType
}
