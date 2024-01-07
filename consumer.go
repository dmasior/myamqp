package myamqp

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

// Consumer represents a consumer.
type Consumer struct {
	channel *amqp091.Channel
	done    chan error
	options *ConsumerOptions
}

// HandleFunc is a function that handles incoming deliveries.
type HandleFunc func(deliveries <-chan amqp091.Delivery, done chan error)

// Consumer creates a new consumer with the given ConsumerOptions and HandleFunc.
func (s *MyAMQP) Consumer(options *ConsumerOptions, handler HandleFunc) (*Consumer, error) {
	if s.conn == nil || s.conn.IsClosed() {
		return nil, ErrNotConnected
	}

	if options == nil {
		return nil, ErrOptionsCannotBeNil
	}

	if options.exchangeOpts == nil {
		return nil, ErrExchangeOptionsCannotBeNil
	}

	if options.queueOpts == nil {
		return nil, ErrQueueOptionsCannotBeNil
	}

	if err := s.channel.ExchangeDeclare(
		options.exchangeOpts.name,
		options.exchangeOpts.kind,
		options.exchangeOpts.durable,
		options.exchangeOpts.autoDelete,
		options.exchangeOpts.internal,
		options.exchangeOpts.noWait,
		options.exchangeOpts.args,
	); err != nil {
		return nil, err
	}

	if _, err := s.channel.QueueDeclare(
		options.queueOpts.name,
		options.queueOpts.durable,
		options.queueOpts.autoDelete,
		options.queueOpts.exclusive,
		options.queueOpts.noWait,
		options.queueOpts.args,
	); err != nil {
		return nil, err
	}

	if err := s.channel.QueueBind(
		options.queueOpts.name,
		options.queueOpts.routingKey,
		options.exchangeOpts.name,
		options.queueOpts.noWait,
		options.queueOpts.args,
	); err != nil {
		return nil, err
	}

	deliveries, err := s.channel.Consume(
		options.queueOpts.name,
		options.name,
		options.autoAck,
		options.exclusive,
		options.noLocal,
		options.noWait,
		options.args,
	)
	if err != nil {
		return nil, err
	}

	consumer := &Consumer{
		options: options,
		channel: s.channel,
		done:    make(chan error),
	}

	go handler(deliveries, consumer.done)

	return consumer, nil
}

// SetCloseContext sets a context that closes the consumer.
func (c *Consumer) SetCloseContext(ctx context.Context) chan error {
	chErr := make(chan error)
	go func() {
		<-ctx.Done()
		if err := c.Cancel(); err != nil {
			chErr <- err
		}
	}()

	return chErr
}

// Cancel cancels the consumer.
func (c *Consumer) Cancel() error {
	if err := c.channel.Cancel(c.options.name, true); err != nil {
		return err
	}

	// Wait for the handler to finish. This is needed because the handler
	// is running in a goroutine.
	return <-c.done
}
