package drain

import "errors"

// ErrInvalidCapacity is returned when Config.Capacity is not positive.
var ErrInvalidCapacity = errors.New("drain: capacity must be greater than zero")

// ErrInvalidFlushTimeout is returned when Config.FlushTimeout is not positive.
var ErrInvalidFlushTimeout = errors.New("drain: flush timeout must be greater than zero")
