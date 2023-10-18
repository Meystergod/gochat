package utils

import "github.com/pkg/errors"

var (
	ErrorDecode                    = errors.New("failed to decode")
	ErrorMarshal                   = errors.New("failed to marshal")
	ErrorConvert                   = errors.New("failed to convert")
	ErrorUnmarshal                 = errors.New("failed to unmarshal")
	ErrorExecuteQuery              = errors.New("failed to execute query")
	ErrorConvertDomainToRepository = errors.New("failed to convert domain model to repository model")
	ErrorGetUrlParams              = errors.New("failed to get param from query url")
	ErrorBindAndValidatePayload    = errors.New("failed to validate or bind payload value")
)
