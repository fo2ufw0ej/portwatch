package decay

import "errors"

// ErrInvalidHalfLife is returned when HalfLife is zero or negative.
var ErrInvalidHalfLife = errors.New("decay: HalfLife must be positive")
