package blog

import (
	"errors"
)

var (
	ErrEntryNotFound = errors.New("No Blog Entry with that ID")
	ErrBadObjectId   = errors.New("error while validating entry ID")
)
