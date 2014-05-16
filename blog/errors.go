package blog

import (
	"errors"
)

var ErrBadObjectId = errors.New("error while validating entry ID")
