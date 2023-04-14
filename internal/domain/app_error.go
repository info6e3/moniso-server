package services

type AppError struct {
	msg string
}

func (err *AppError) Error() string {
	return err.msg
}

func NewAppError(msg string) *AppError {
	return &AppError{
		msg: msg,
	}
}
