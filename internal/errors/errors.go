package errors

import "fmt"

type ApplicationError interface {
	Status() int
	Error() string
}

type conflictError struct {
	code    int
	message string
}

func ConflictError(message string) *conflictError {
	return &conflictError{
		code:    409,
		message: message,
	}
}

func (e *conflictError) Error() string {
	return fmt.Sprintf("RESOURCE_ALREADY_EXISTS: %s", e.message)
}

func (e *conflictError) Status() int {
	return e.code
}

type unprocessableError struct {
	code    int
	message string
}

func UnprocessableError(message string) *unprocessableError {
	return &unprocessableError{
		code:    422,
		message: message,
	}
}

func (e *unprocessableError) Error() string {
	return fmt.Sprintf("UNPROCESSABLE_ENTITY: %s", e.message)
}

func (e *unprocessableError) Status() int {
	return e.code
}

type unauthorizedError struct {
	code    int
	message string
}

func UnauthorizedError(message string) *unauthorizedError {
	return &unauthorizedError{
		code:    401,
		message: message,
	}
}

func (e *unauthorizedError) Error() string {
	return fmt.Sprintf("UNAUTHORIZED: %s", e.message)
}

func (e *unauthorizedError) Status() int {
	return e.code
}

type badRequestError struct {
	code    int
	message string
}

func BadRequestError(message string) *badRequestError {
	return &badRequestError{
		code:    400,
		message: message,
	}
}

func (e *badRequestError) Error() string {
	return fmt.Sprintf("BAD_REQUEST: %s", e.message)
}

func (e *badRequestError) Status() int {
	return e.code
}
