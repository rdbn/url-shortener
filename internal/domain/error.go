package domain

import "errors"

var (
	ErrNotFound          = errors.New("url not found")
	ErrDuplicate         = errors.New("url already exists")
	ErrSaveURL           = errors.New("failed to save url")
	ErrGetURL            = errors.New("failed to get url")
	ErrURLIsRequired     = errors.New("url is required")
	ErrShortCodeNotFound = errors.New("short code not found")
)
