package repositories

import (
	"errors"
)

// Contains common code designed to be used in gui web-client interface.

var EntityNotExistsErr error = errors.New("Entity doesn't exists")
var DuplicatedEntityErr error = errors.New("Duplicated Entity")
var InvalidEntityStateErr error = errors.New("Entity state is invalid")
