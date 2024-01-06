package myamqp

type Config struct {
	dialFn DialFunc
	retry  *Retry
}

// NewConfig creates a new Config with the given DialFunc.
func NewConfig(dialFn DialFunc) *Config {
	return &Config{
		dialFn: dialFn,
	}
}

// WithRetry sets the Retry on the Config.
func (c *Config) WithRetry(retry *Retry) *Config {
	c.retry = retry
	return c
}

func (c *Config) DialFn() DialFunc {
	return c.dialFn
}

func (c *Config) Retry() *Retry {
	return c.retry
}
