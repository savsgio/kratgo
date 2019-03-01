package invalidator

import "errors"

// ErrEmptyFields ...
var ErrEmptyFields = errors.New("Minimum one mandatory field")

// ErrMaxWorkersZero ...
var ErrMaxWorkersZero = errors.New("MaxWorkers must be greater than 0")
