package domain_error

import "net/http"

type DomainError struct {
	httpStatus int
	code       string
	message    string
	err        error
}

func NewError(httpStatus int, code, message string, cause error) error {
	return &DomainError{
		httpStatus: httpStatus,
		code:       code,
		message:    message,
		err:        cause,
	}
}

func (e *DomainError) Error() string {
	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}
	return e.message
}

func (e *DomainError) Unwrap() error   { return e.err }
func (e *DomainError) Code() string    { return e.code }
func (e *DomainError) Message() string { return e.message }
func (e *DomainError) HTTPStatus() int { return e.httpStatus }

func ObjectNotExists() error {
	return &DomainError{
		httpStatus: http.StatusNotFound,
		code:       "OBJECT_NOT_EXISTS",
		message:    "Object Not Exists",
		err:        nil, //FIXME
	}
}
