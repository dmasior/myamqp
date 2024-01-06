package myamqp

import (
	"errors"
	"sync"
	"time"
)

var (
	NoRetry         = NewRetry(0, 0)
	MaxRetryForever = -1
	ErrMaxRetries   = errors.New("max retries reached")
)

type Retry struct {
	count   int
	countMu sync.Mutex
	max     int
	wait    time.Duration
}

// NewRetry creates a new Retry with the given max and wait.
// If max is -1, the Retry retries forever.
// If wait is 0, the Retry never waits.
func NewRetry(max int, wait time.Duration) *Retry {
	return &Retry{
		max:  max,
		wait: wait,
	}
}

func (r *Retry) Count() int {
	return r.count
}

func (r *Retry) Max() int {
	return r.max
}

func (r *Retry) Wait() time.Duration {
	return r.wait
}

func (r *Retry) Inc() {
	r.countMu.Lock()
	defer r.countMu.Unlock()
	r.count++
}
