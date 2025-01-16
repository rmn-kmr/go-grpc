package errors

import "github.com/rmnkmr/go-common/errors"

var (
	ErrProviderErr    = errors.From(400, "provider_error")
	ErrInvalidRequest = errors.From(400, "invalid_request")
)
