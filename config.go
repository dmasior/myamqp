package myamqp

import "errors"

var (
	ErrDialFuncCannotBeNil = errors.New("dial func cannot be nil")
)

// Config represents a configuration.
// It is used to configure a MyAMQP instance.
type Config struct {
	dialFunc        DialFunc
	reconnectPolicy *ReconnectPolicy
	onConnect       func(*MyAMQP)
	qos             *Qos
}

// NewConfig creates a new Config with the given URL.
func NewConfig(dialFunc DialFunc) (*Config, error) {
	if dialFunc == nil {
		return nil, ErrDialFuncCannotBeNil
	}

	return &Config{
		dialFunc: dialFunc,
	}, nil
}

// WithReconnectPolicy sets the ReconnectPolicy on the Config.
func (c *Config) WithReconnectPolicy(rp *ReconnectPolicy) *Config {
	c.reconnectPolicy = rp
	return c
}

// WithOnConnect sets the OnConnect callback on the Config.
func (c *Config) WithOnConnect(callback func(*MyAMQP)) *Config {
	c.onConnect = callback
	return c
}

// WithQos sets the Qos on the Config.
func (c *Config) WithQos(qos *Qos) *Config {
	c.qos = qos
	return c
}

func (c *Config) OnConnect() func(*MyAMQP) {
	return c.onConnect
}

func (c *Config) Qos() *Qos {
	return c.qos
}

func (c *Config) ReconnectPolicy() *ReconnectPolicy {
	return c.reconnectPolicy
}

func (c *Config) DialFunc() DialFunc {
	return c.dialFunc
}
