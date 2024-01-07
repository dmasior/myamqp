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
	channel *amqp091.Channel
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

// Connect establishes a connection to the AMQP server and returns the connection object.
// If the connection is already established, it returns the existing connection object.
// If the connection is closed, it attempts to reconnect.
// If the context is done, it returns an error.
func (s *MyAMQP) Connect(ctx context.Context) (*amqp091.Connection, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.conn != nil && !s.conn.IsClosed() {
		return s.conn, nil
	}

	// Connect to the AMQP server.
	conn, err := s.config.DialFunc()()
	if err != nil {
		return nil, err
	}
	s.conn = conn

	// Create a new channel.
	s.channel, err = s.conn.Channel()
	if err != nil {
		s.conn.Close()
		s.conn = nil
		return nil, err
	}

	if s.config.Qos() != nil {
		err = s.channel.Qos(
			s.config.Qos().PrefetchCount(),
			s.config.Qos().PrefetchSize(),
			s.config.Qos().Global(),
		)
		if err != nil {
			s.conn.Close()
			s.conn = nil
			return nil, err
		}
	}

	if s.config.OnConnect() != nil {
		s.config.OnConnect()(s)
	}

	// Setup reconnectPolicy.
	if s.config.ReconnectPolicy() != nil {
		s.setupReconnect(ctx)
	}

	return s.conn, nil
}

func (s *MyAMQP) setupReconnect(ctx context.Context) {
	errListener := s.config.ReconnectPolicy().ErrListener()
	if errListener == nil {
		errListener = func(err error) {}
	}

	s.rCtx, s.rCancel = context.WithCancel(ctx)

	go func() {
		for {
			select {
			// Check if the context is done after we try to connect.
			case <-s.rCtx.Done():
				if s.conn != nil && !s.conn.IsClosed() {
					err := s.conn.Close()
					if err != nil {
						errListener(err)
					}
				}
				errListener(s.rCtx.Err())
				return
			default:
				if s.conn != nil && !s.conn.IsClosed() {
					<-time.After(time.Millisecond * 100)
					continue
				}

				reconnectPolicy := s.config.ReconnectPolicy()
				if reconnectPolicy.Max() != MaxReconnectUnlimited && reconnectPolicy.Count() >= reconnectPolicy.Max() {
					errListener(ErrMaxReconnectsReached)
					return
				}

				reconnectPolicy.Inc()
				<-time.After(reconnectPolicy.Backoff())

				conn, err := s.config.DialFunc()()
				if err != nil {
					s.conn = nil
					errListener(err)
					continue
				}
				s.conn = conn
				s.channel, err = s.conn.Channel()
				if err != nil {
					errListener(err)
					err = s.conn.Close()
					if err != nil {
						errListener(err)
					}
					s.conn = nil
					continue
				}

				if s.config.Qos() != nil {
					err = s.channel.Qos(
						s.config.Qos().PrefetchCount(),
						s.config.Qos().PrefetchSize(),
						s.config.Qos().Global(),
					)
					if err != nil {
						errListener(err)
						err = s.conn.Close()
						if err != nil {
							errListener(err)
						}
						s.conn = nil
						continue
					}
				}

				if s.config.onConnect != nil {
					s.config.onConnect(s)
				}
			}
		}
	}()
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
