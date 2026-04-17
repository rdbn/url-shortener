package domain

import "errors"

var (
	ErrNotFound             = errors.New("url not found")
	ErrDuplicate            = errors.New("url already exists")
	ErrSaveURL              = errors.New("failed to save url")
	ErrGetURL               = errors.New("failed to get url")
	ErrBatchIncrementClicks = errors.New("failed batch to increment clicks")
	ErrURLIsRequired        = errors.New("url is required")
	ErrShortCodeNotFound    = errors.New("short code not found")
	ErrUnexpectedType       = errors.New("unexpected type from singleflight")
)
