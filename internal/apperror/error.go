package apperror

import (
	"encoding/json"
	"errors"
)

var (
	ErrorDecode       = errors.New("failed to decode")
	ErrorConvertId    = errors.New("failed to convert id")
	ErrorConvertModel = errors.New("failed to convert domain model to repository model")
	ErrorCreateOne    = errors.New("failed to insert object into database")
	ErrorGetOne       = errors.New("failed to get object from database")
	ErrorGetAll       = errors.New("failed to get all objects from database")
	ErrorUpdateOne    = errors.New("failed to update object in database")
	ErrorDeleteOne    = errors.New("failed to delete object from database")
	ErrorNotFound     = errors.New("failed to get items")
)

type AppError struct {
	err     error
	message string
}

func NewAppError(err error, message string) *AppError {
	return &AppError{
		err:     err,
		message: message,
	}
}

func (e *AppError) Error() string {
	return e.err.Error()
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}
