package clamp

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
