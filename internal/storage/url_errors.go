package storage

import (
	"errors"
)

var (
	ErrorNotFound = errors.New("url not found")
	ErrorConflict = errors.New("custom short id already exists")
)
