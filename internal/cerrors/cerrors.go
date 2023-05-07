package cerrors

type ApplicationError struct {
	code    int
	message string
}

func (e ApplicationError) Error() string {
	return e.message
}

func (e ApplicationError) Code() int {
	return e.code
}

type conflictError struct {
	ApplicationError
}

func ConflictError(message string) *conflictError {
	return &conflictError{
		ApplicationError: ApplicationError{
			code:    409,
			message: message,
		},
	}
}

type unprocessableError struct {
	ApplicationError
}

func UnprocessableError(message string) *unprocessableError {
	return &unprocessableError{
		ApplicationError: ApplicationError{
			code:    422,
			message: message,
		},
	}
}

type unauthorizedError struct {
	ApplicationError
}

func UnauthorizedError(message string) *unauthorizedError {
	return &unauthorizedError{
		ApplicationError: ApplicationError{
			code:    401,
			message: message,
		},
	}
}

type badRequestError struct {
	ApplicationError
}

func BadRequestError(message string) *badRequestError {
	return &badRequestError{
		ApplicationError: ApplicationError{
			code:    400,
			message: message,
		},
	}
}
