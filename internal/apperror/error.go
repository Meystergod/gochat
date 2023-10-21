package apperror

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
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

func HTTPAppErrorHandler(ctx context.Context) func(err error, c echo.Context) {
	logger := zerolog.Ctx(ctx)

	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var appError *AppError
		if errors.As(err, &appError) {
			switch {
			case errors.Is(appError.Err, ErrorDecode),
				errors.Is(appError.Err, ErrorCreateOne),
				errors.Is(appError.Err, ErrorConvert),
				errors.Is(appError.Err, ErrorGetOne),
				errors.Is(appError.Err, ErrorGetAll),
				errors.Is(appError.Err, ErrorUpdateOne),
				errors.Is(appError.Err, ErrorDeleteOne),
				errors.Is(appError.Err, ErrorConvertModel):

				appError.ErrorMessage = appError.Error()
				if jsonError := c.JSON(http.StatusInternalServerError, &appError); jsonError != nil {
					logger.Error().Msgf("failed to create json response: %s", jsonError.Error())
				}
				return
			case errors.Is(appError.Err, ErrorValidatePayload):
				appError.ErrorMessage = appError.Error()
				if jsonError := c.JSON(http.StatusBadRequest, &appError); jsonError != nil {
					logger.Error().Msgf("failed to validate json response: %s", jsonError.Error())
				}
				return
			}
		}
		if jsonError := c.JSON(http.StatusNotFound, map[string]string{"error": "page not found"}); jsonError != nil {
			logger.Error().Msgf("failed to validate json response: %s", jsonError.Error())
		}
	}
}
