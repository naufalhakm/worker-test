package commons

type PermanentError struct {
	Message string
}

func (e *PermanentError) Error() string {
	return e.Message
}

func NewPermanentError(msg string) *PermanentError {
	return &PermanentError{Message: msg}
}
