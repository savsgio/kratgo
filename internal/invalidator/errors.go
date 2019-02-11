package invalidator

import "errors"

var errInvalidAction = errors.New("Invalid action")
var errEmptyFields = errors.New("Minimum one mandatory field")
