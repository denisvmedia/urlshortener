package storage

import "errors"

var (
	// ErrNotFound is returned when a link is not found in the storage
	ErrNotFound = errors.New("not found")
	// ErrShortNameAlreadyExists is returned when a link already exists in the storage
	ErrShortNameAlreadyExists = errors.New("given short name is already used by another link")
	// ErrStorageFailure is returned in case of a storage problem
	ErrStorageFailure = errors.New("storage failure")
)
