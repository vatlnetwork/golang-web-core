package srverr

type ServerError struct {
	Message string
	Code    int
}

func New(message string, code ...int) ServerError {
	actualCode := 500
	if len(code) > 0 {
		actualCode = code[0]
	}

	return ServerError{
		Message: message,
		Code:    actualCode,
	}
}

func (e ServerError) Error() string {
	return e.Message
}
