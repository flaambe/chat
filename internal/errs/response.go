package errs

type ResponseError struct {
	Status  int
	Message string
	Err     error
}

func (e *ResponseError) Error() string {
	return e.Message
}

func New(status int, message string, err error) *ResponseError {
	return &ResponseError{
		Status:  status,
		Message: message,
		Err:     err,
	}
}
