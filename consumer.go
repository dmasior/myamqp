package myamqp

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel *amqp091.Channel
	done    chan error
	options ConsumerOptions
}

type HandleFunc func(deliveries <-chan amqp091.Delivery, done chan error)

func (s *MyAMQP) Consumer(options ConsumerOptions, handler HandleFunc) (*Consumer, error) {
	if s.conn == nil {
		return nil, ErrNotConnected
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

func (c *Consumer) Cancel() error {
	if err := c.channel.Cancel(c.options.name, true); err != nil {
		return err
	}

	return <-c.done
}
