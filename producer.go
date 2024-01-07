package myamqp

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

// Producer represents an AMQP producer.
type Producer struct {
	channel *amqp091.Channel
	options *ProducerOptions
}

// Producer creates a new producer with the given ProducerOptions.
func (s *MyAMQP) Producer(options *ProducerOptions) (*Producer, error) {
	if s.conn == nil || s.conn.IsClosed() {
		return nil, ErrNotConnected
	}

	if options == nil {
		return nil, ErrOptionsCannotBeNil
	}

	if options.exchangeOpts == nil {
		return nil, ErrExchangeOptionsCannotBeNil
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

	if options.queueOpts != nil {
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
	}

	producer := &Producer{
		options: options,
		channel: s.channel,
	}

	return producer, nil
}

// Publish publishes a message to the AMQP server.
func (s *Producer) Publish(ctx context.Context, routingKey string, mandatory, immediate bool, msg amqp091.Publishing) error {
	return s.channel.PublishWithContext(
		ctx,
		s.options.exchangeOpts.name,
		routingKey,
		mandatory,
		immediate,
		msg,
	)
}

// PublishWithDeferredConfirm publishes a message to the AMQP server and returns a DeferredConfirmation.
func (s *Producer) PublishWithDeferredConfirm(ctx context.Context, routingKey string, mandatory, immediate bool, msg amqp091.Publishing) (*amqp091.DeferredConfirmation, error) {
	return s.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		s.options.exchangeOpts.name,
		routingKey,
		mandatory,
		immediate,
		msg,
	)
}
