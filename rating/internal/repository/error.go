package repository

import "errors"

// ErrNotFound is returned whena a requested record is not found.
var ErrNotFound = errors.New("not found")
