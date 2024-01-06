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

type MyAMQP struct {
	config     Config
	conn       *amqp091.Connection
	connErrors chan error
	connMu     sync.Mutex
	channel    *amqp091.Channel
}

type DialFunc func() (*amqp091.Connection, error)

func New(config Config) *MyAMQP {
	if config.Retry() == nil {
		config.retry = NoRetry
	}

	return &MyAMQP{
		config: config,
	}
}

func (s *MyAMQP) Connect(ctx context.Context) (*amqp091.Connection, chan error) {
	if s.conn != nil {
		return s.conn, s.connErrors
	}

	s.connErrors = s.connect(ctx)

	return s.conn, s.connErrors
}

func (s *MyAMQP) connect(ctx context.Context) chan error {
	errCh := make(chan error)
	go func() {
		for {
			// Check if the context is done before we try to connect.
			select {
			case <-ctx.Done():
				if s.conn == nil {
					errCh <- ctx.Err()
					return
				}
			default:
			}

			select {
			// Check if the context is done after we try to connect.
			case <-ctx.Done():
				s.connMu.Lock()
				if s.conn != nil && !s.conn.IsClosed() {
					err := s.conn.Close()
					if err != nil {
						errCh <- err
					}
				}
				errCh <- ctx.Err()
				s.connMu.Unlock()
				return
			default:
				s.connMu.Lock()
				if s.conn != nil && !s.conn.IsClosed() {
					<-time.After(time.Millisecond * 100)
					continue
				}

				retry := s.config.Retry()

				if retry.Max() != MaxRetryForever && retry.Count() >= retry.Max() {
					errCh <- ErrMaxRetries
					return
				}

				retry.Inc()
				<-time.After(retry.Wait())

				conn, err := s.config.DialFn()()
				if err != nil {
					errCh <- err
					continue
				}
				s.conn = conn
				if s.channel == nil {
					s.channel, err = s.conn.Channel()
					if err != nil {
						errCh <- err
						continue
					}
				}
				s.connMu.Unlock()
			}
		}
	}()

	return errCh
}
