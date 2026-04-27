package clamp

import "fmt"

// NewFromConfig constructs a Clamper from the provided Config,
// returning an error if validation fails.
func NewFromConfig(cfg Config) (*Clamper, error) {
	return New(cfg)
}

// NewDefault constructs a Clamper using DefaultConfig.
func NewDefault() *Clamper {
	c, _ := New(DefaultConfig())
	return c
}

// MustNewFromConfig constructs a Clamper from the provided Config,
// panicking if validation fails. Useful in init() or test setup where
// an invalid config is a programming error.
func MustNewFromConfig(cfg Config) *Clamper {
	c, err := New(cfg)
	if err != nil {
		panic(fmt.Sprintf("clamp: invalid config: %v", err))
	}
	return c
}
