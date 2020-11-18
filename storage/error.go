package storage

import "errors"

var (
	ErrNotFound               = errors.New("not found")
	ErrShortNameAlreadyExists = errors.New("given short name is already used by another link")
	ErrStorageFailure         = errors.New("storage failure")
)
