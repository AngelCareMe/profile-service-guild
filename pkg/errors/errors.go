package errors

import "fmt"

type Error interface {
	error
	Unwrap() error
}

type appError struct {
	Message string
	Err     error
}

type httpError struct {
	Code    int
	Message string
	Err     error
}

func (ae *appError) Error() string {
	if ae.Err != nil {
		return fmt.Sprintf("%s: %v", ae.Message, ae.Err)
	}
	return ae.Message
}

func (ae *appError) Unwrap() error {
	return ae.Err
}

func NewAppError(msg string, err error) *appError {
	return &appError{
		Message: msg,
		Err:     err,
	}
}

func (he *httpError) Error() string {
	if he.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", he.Code, he.Message, he.Err)
	}

	return fmt.Sprintf("[%d] %s", he.Code, he.Message)
}

func (he *httpError) Unwrap() error {
	return he.Err
}

func NewHTTPError(code int, msg string, err error) *httpError {
	return &httpError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}
