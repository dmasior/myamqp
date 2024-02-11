package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dmasior/myamqp"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Create a new context with a done channel.
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	// Create a new Config with a dial function that returns a new AMQP connection.
	config, err := myamqp.NewConfig(func() (*amqp091.Connection, error) {
		return amqp091.Dial("amqp://guest:guest@localhost:5672/")
	})
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
	// Setup on connect hook.
	config = config.
		WithOnConnect(func(myAMQP *myamqp.MyAMQP) {
			slog.InfoContext(ctx, "connected")
			setupProducer(ctx, myAMQP)
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

	slog.InfoContext(ctx, "context done, bye")
}

func setupProducer(ctx context.Context, amqp *myamqp.MyAMQP) {
	// Create a new ProducerOptions.
	producerOptions := myamqp.NewProducerOptions(
		myamqp.NewExchangeOptions("/", myamqp.ExchangeTypeDirect),
	)

	// Create a new Producer.
	producer, err := amqp.Producer(producerOptions)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	// Start a goroutine to publish messages.
	go func() {
		for i := 0; ; i++ {
			// Check if the context is done before we try to publish a message.
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Publish a message.
			err = producer.Publish(ctx, "", false, false, amqp091.Publishing{
				Body: []byte("hello world, " + strconv.Itoa(i)),
			})
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				return
			}

			slog.InfoContext(ctx, fmt.Sprintf("published `%d` message", i))

			<-time.After(1 * time.Second)
		}
	}()
}
