# MyAMQP
The AMQP tools for Golang with developer-friendly API. It's a wrapper around the [rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go) library, providing helper functions for faster and easier development.
Reconnections, context support (graceful shutdown), and more.

## Docs
[![GoDoc](https://godoc.org/github.com/dmasior/myamqp?status.svg)](http://godoc.org/github.com/dmasior/myamqp)

## Installation
```bash
go get github.com/dmasior/myamqp
```

## Usage
### New MyAMQP instance
```go
// Create a new Config. It needs a dialer function which returns an amqp091.Connection.
config, err := myamqp.NewConfig(func() (*amqp091.Connection, error) {
    return amqp091.Dial("amqp://guest:guest@localhost:5672/")
})
if err != nil {
    // handle error
}
config = config.
    // WithOnConnect allows to set a callback which is called after a successful connection.
    WithOnConnect(func(myAMQP *myamqp.MyAMQP) {
        slog.InfoContext(ctx, "connected")
        // You can also start a consumer here.
        // See example consumer in examples/consumer/main.go
    }).
	// WithReconnectPolicy allows to set a reconnect policy.
    WithReconnectPolicy(
        myamqp.NewReconnectPolicy(myamqp.MaxReconnectUnlimited, 1*time.Second).
            // ReconnectPolicy can be extended with error listener.
            WithErrorListener(func(err error) {
                slog.ErrorContext(ctx, err.Error())
            }),
    )

// Create and run a new MyAMQP.
amqp, err := myamqp.New(config)
if err != nil {
    // handle error
}

// Use amqp.Run after setting up all consumers and producers in `WithOnConnect` callback.
err = amqp.Run(ctx)
if err != nil {
    // handle error
}
```

### Consumer
```go
// Create a new ConsumerOptions
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

// Attach a new consumer to the MyAMQP.
consumer, err := amqp.Consumer(consumerOptions, handler)
if err != nil {
    // handle error
}

// Set a close context. It allows to close the consumer when the context is done.
consumer.SetCloseContext(ctx)
```

See full [consumer example](./examples/consumer/main.go)

### Producer
```go
// Create a new ProducerOptions.
producerOptions := myamqp.NewProducerOptions(
    myamqp.NewExchangeOptions("/", myamqp.ExchangeTypeDirect),
)

// Attach a new producer to the MyAMQP.
producer, err := amqp.Producer(producerOptions)
if err != nil {
    // handle error
}

// Publish a message.
err = producer.Publish(ctx, "", false, false, amqp091.Publishing{
    Body: []byte("hello world"),
})
if err != nil {
	// handle error
}
```
See full [producer example](./examples/producer/main.go)
