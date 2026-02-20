package internal

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ExitCodeError struct {
	error
	exitCode int
}

func NewExitCodeError(err error, exitCode int) error {
	return &ExitCodeError{
		error:    err,
		exitCode: exitCode,
	}
}

func (e *ExitCodeError) ExitCode() int {
	return e.exitCode
}

type ApiError struct {
	error
	Code    int               `json:"code"`
	Status  int               `json:"status"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

func NewApiError(code int, message string, err error, details ...string) error {
	detailsMap := make(map[string]string)
	if len(details)%2 == 0 {
		for i := 0; i < len(details); i += 2 {
			detailsMap[details[i]] = details[i+1]
		}
	}

	status := func(code int) int {
		if 100 <= code && code < 200 {
			return 100
		} else if 200 <= code && code < 300 {
			return 200
		} else if 300 <= code && code < 400 {
			return 300
		} else if 400 <= code && code < 500 {
			return 400
		}
		return 500
	}

	return &ApiError{
		error:   err,
		Code:    code,
		Status:  status(code),
		Message: message,
		Details: detailsMap,
	}
}

func (e *ApiError) Unwrap() error {
	return e.error
}

func (e *ApiError) Error() string {
	if e.error == nil {
		return e.Message
	}

	if e.Message == "" {
		return e.error.Error()
	}
	return fmt.Sprintf("%v: %s", e.Message, e.error.Error())
}

func (e *ApiError) MarshalJSON() ([]byte, error) {
	var errStr string
	if e.error != nil {
		errStr = e.error.Error()
	}

	return json.Marshal(&struct {
		Code    int               `json:"code"`
		Status  int               `json:"status"`
		Message string            `json:"message,omitempty"`
		Error   string            `json:"error,omitempty"`
		Details map[string]string `json:"details,omitempty"`
	}{
		Code:    e.Code,
		Status:  e.Status,
		Message: e.Message,
		Error:   errStr,
		Details: e.Details,
	})
}

type ValidateError struct {
	error
}

func NewValidateError(err error) error {
	return &ValidateError{
		error: err,
	}
}

func (e *ValidateError) Unwrap() error {
	return e.error
}

type NotFoundError struct {
	error
}

func NewNotFoundError(err error) error {
	return &NotFoundError{
		error: err,
	}
}

func (e NotFoundError) Unwrap() error {
	return e.error
}

func IsNotFoundError(err error) bool {
	var notFoundError *NotFoundError
	ok := errors.As(err, &notFoundError)
	return ok
}
