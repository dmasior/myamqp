package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmasior/myamqp"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Create a new context with a done channel.
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	// Create a new Config. It needs a dialer function which returns an amqp091.Connection.
	config, err := myamqp.NewConfig(func() (*amqp091.Connection, error) {
		return amqp091.Dial("amqp://guest:guest@localhost:5672/")
	})
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	// Setup on connect callback.
	config = config.
		WithOnConnect(func(myAMQP *myamqp.MyAMQP) {
			slog.InfoContext(ctx, "connected")
			setupConsumer(ctx, myAMQP)
		}).
		// Setup reconnect policy.
		WithReconnectPolicy(
			myamqp.NewReconnectPolicy(myamqp.MaxReconnectUnlimited, 1*time.Second).
				// Setup error listener.
				WithErrorListener(func(err error) {
					slog.ErrorContext(ctx, err.Error())
				}),
		)

	// Create and run a new MyAMQP.
	amqp, err := myamqp.New(config)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	err = amqp.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	slog.InfoContext(ctx, "done, bye")
}

func setupConsumer(ctx context.Context, amqp *myamqp.MyAMQP) {
	// Create a new ConsumerOptions.
	consumerOptions := myamqp.NewConsumerOptions(
		"consumer-tag",
		myamqp.NewExchangeOptions("/", myamqp.ExchangeTypeDirect),
		myamqp.NewQueueOptions("queue-name"),
	)

	// Deliveries handler.
	handler := func(deliveries <-chan amqp091.Delivery, done chan error) {
		for d := range deliveries {
			slog.InfoContext(ctx, string(d.Body))
			d.Ack(false)
		}

		// Signal that the handler is done. It allows consumer.Cancel() to return.
		done <- nil
	}

	// Attach a consumer to the MyAMQP.
	consumer, err := amqp.Consumer(consumerOptions, handler)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	// Set a close context. It allows to close the consumer when the context is done.
	consumer.SetCloseContext(ctx)
}
