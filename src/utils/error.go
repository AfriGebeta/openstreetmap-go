package utils

type ApiError struct {
	StatusCode  int               `json:"code"`
	Errors      map[string]string `json:"message"`
	Message     string            `json:"message"`
	CausedError error             `json:"error"`
}

func NewApiError(statusCode int, message string, causedError error, errors map[string]string) *ApiError {
	return &ApiError{
		StatusCode:  statusCode,
		Errors:      errors,
		Message:     message,
		CausedError: causedError,
	}
}

func (err *ApiError) Error() string {
	return err.Message
}
