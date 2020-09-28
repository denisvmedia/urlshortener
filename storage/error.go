package storage

import "errors"

var ErrNotFound = errors.New("not found")
var ErrShortNameAlreadyExists = errors.New("given short name is already used by another link")
