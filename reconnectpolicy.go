package myamqp

import (
	"errors"
	"sync"
	"time"
)

const MaxReconnectUnlimited = -1

var (
	ErrMaxReconnectsReached = errors.New("max reconnects reached")
)

type ReconnectPolicy struct {
	count       int
	countMu     sync.Mutex
	max         int
	backoff     time.Duration
	errListener func(error)
}

// NewReconnectPolicy creates a new ReconnectPolicy with the given max and backoff.
// If max is -1, the ReconnectPolicy retries indefinitely.
// If backoff is 0, the ReconnectPolicy retries immediately.
func NewReconnectPolicy(max int, backoff time.Duration) *ReconnectPolicy {
	return &ReconnectPolicy{
		max:     max,
		backoff: backoff,
	}
}

func (r *ReconnectPolicy) WithErrorListener(handler func(error)) *ReconnectPolicy {
	r.errListener = handler
	return r
}

func (r *ReconnectPolicy) Count() int {
	return r.count
}

func (r *ReconnectPolicy) Max() int {
	return r.max
}

func (r *ReconnectPolicy) Backoff() time.Duration {
	return r.backoff
}

func (r *ReconnectPolicy) ErrListener() func(error) {
	return r.errListener
}

func (r *ReconnectPolicy) Inc() {
	r.countMu.Lock()
	defer r.countMu.Unlock()
	r.count++
}
