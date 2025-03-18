package openroutergo

import "errors"

var (
	// ErrBaseURLRequired is returned when the base URL is needed but not provided.
	ErrBaseURLRequired = errors.New("the base URL is required")

	// ErrAPIKeyRequired is returned when the API key is needed but not provided.
	ErrAPIKeyRequired = errors.New("the API key is required")

	// ErrMessagesRequired is returned when no messages are found and they are needed.
	ErrMessagesRequired = errors.New("at least one message is required")
)
