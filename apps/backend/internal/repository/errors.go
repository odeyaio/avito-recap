package repository

import "errors"

var (
	ErrProfileNotFound    = errors.New("profile not found")
	ErrRecapNotFound      = errors.New("recap not found")
	ErrCatalogUnavailable = errors.New("catalog unavailable")
)
