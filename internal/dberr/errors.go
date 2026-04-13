package dberr

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrDBNil       = errors.New("database is not initialized")
)
