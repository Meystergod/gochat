package utils

import "github.com/pkg/errors"

var (
	ErrorGetUrlParams           = errors.New("failed to get param from query url")
	ErrorBindAndValidatePayload = errors.New("failed to validate or bind payload value")
)
