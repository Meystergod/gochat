package apperror

import (
	"errors"
)

var (
	ErrorDecode          = errors.New("failed to decode")
	ErrorConvert         = errors.New("failed to convert")
	ErrorConvertModel    = errors.New("failed to convert domain model to repository model")
	ErrorCreateOne       = errors.New("failed to insert object into database")
	ErrorGetOne          = errors.New("failed to get object from database")
	ErrorGetAll          = errors.New("failed to get all objects from database")
	ErrorUpdateOne       = errors.New("failed to update object in database")
	ErrorDeleteOne       = errors.New("failed to delete object from database")
	ErrorValidatePayload = errors.New("failed to validate or bind payload value")
	ErrorGetUrlParams    = errors.New("failed to get param from query url")
)

type AppError struct {
	Err          error  `json:"-"`
	ErrorMessage string `json:"error"`
	Message      string `json:"message"`
}

func NewAppError(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}
