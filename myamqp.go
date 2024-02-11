package myamqp

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var (
	ErrNotConnected = errors.New("not connected")
)

// MyAMQP represents a connection to an AMQP server.
type MyAMQP struct {
	config  *Config
	conn    *amqp091.Connection
	connMu  sync.Mutex
	rCtx    context.Context
	rCancel context.CancelFunc
}

// DialFunc is a function that returns a new AMQP connection.
// It is used by the Config to establish a connection to the AMQP server.
type DialFunc func() (*amqp091.Connection, error)

// New creates a new MyAMQP with the given Config.
func New(config *Config) (*MyAMQP, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	return &MyAMQP{
		config: config,
	}, nil
}

// Run runs the MyAMQP. It connects to the AMQP server and handles reconnects.
func (s *MyAMQP) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	errListener := s.config.ReconnectPolicy().ErrListener()
	if errListener == nil {
		errListener = func(err error) {}
	}

	s.rCtx, s.rCancel = context.WithCancel(ctx)

	reconnErrCh := make(chan error, 1)
	_, errCh, err := s.connect()
	if err != nil {
		errListener(err)
		reconnErrCh <- err
	}

	for {
		select {
		case <-s.rCtx.Done():
			if s.conn != nil && !s.conn.IsClosed() {
				cErr := s.conn.Close()
				if cErr != nil {
					errListener(cErr)
				}
			}
			errListener(s.rCtx.Err())
			return s.rCtx.Err()
		case <-reconnErrCh:
			reconnectPolicy := s.config.ReconnectPolicy()
			if reconnectPolicy.Max() != MaxReconnectUnlimited && reconnectPolicy.Count() >= reconnectPolicy.Max() {
				errListener(ErrMaxReconnectsReached)
				return ErrMaxReconnectsReached
			}

			// Setup new error channel and reconnect.
			reconnectPolicy.Inc()
			<-time.After(reconnectPolicy.Backoff())

			var err error
			_, errCh, err = s.connect()
			if err != nil {
				errListener(err)
				reconnErrCh <- err
			}
		case rcErr := <-errCh:
			reconnErrCh <- rcErr
		default:
			<-time.After(time.Millisecond * 500)
		}
	}
}

func (s *MyAMQP) connect() (*amqp091.Connection, chan *amqp091.Error, error) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	errCh := make(chan *amqp091.Error)

	// Run to the AMQP server.
	conn, err := s.config.DialFunc()()
	if err != nil {
		return nil, errCh, err
	}
	s.conn = conn

	if s.config.OnConnect() != nil {
		s.config.OnConnect()(s)
	}

	s.conn.NotifyClose(errCh)

	return s.conn, errCh, nil
}

// Close closes the connection to the AMQP server and cancels the reconnect goroutine.
func (s *MyAMQP) Close() error {
	// Cancel the context, stopping the setupReconnect goroutine
	if s.rCancel != nil {
		s.rCancel()
	}

	// Close the connection and nil it
	var err error
	if s.conn != nil && !s.conn.IsClosed() {
		err = s.conn.Close()
		s.conn = nil
	}

	return err
}
