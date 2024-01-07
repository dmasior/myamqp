# MyAMQP
The AMQP tools for Golang with modern and developer friendly API.
Out of the box: reconnections, context support (graceful shutdown) and more.
It wraps [rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go) library.

## Docs
[![GoDoc](https://godoc.org/github.com/dmasior/myamqp?status.svg)](http://godoc.org/github.com/dmasior/myamqp)

## Installation
```bash
go get github.com/dmasior/myamqp
```

## Usage
### Configuration and connection
```go
// Create a new Config. It needs a dialer function which returns an amqp091.Connection.
config, err := myamqp.NewConfig(func() (*amqp091.Connection, error) {
    return amqp091.Dial("amqp://guest:guest@localhost:5672/")
})
if err != nil {
    // handle error
}
// Setup on connect callback.
config = config.
    WithOnConnect(func(myAMQP *myamqp.MyAMQP) {
        slog.InfoContext(ctx, "connected")
        // You can also start a consumer here.
        // See example consumer in examples/consumer/main.go
    }).
    // Setup reconnect policy.
    WithReconnectPolicy(
        myamqp.NewReconnectPolicy(myamqp.MaxReconnectUnlimited, 1*time.Second).
            // Setup error listener.
            WithErrorListener(func(err error) {
                slog.ErrorContext(ctx, err.Error())
            }),
    )

// Create a new MyAMQP instance and connect to the server.
instance, err := myamqp.New(config)
_, err = instance.Connect(ctx)
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

// Start a consumer.
consumer, err := myAMQP.Consumer(consumerOptions, handler)
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

// Create a new Producer.
producer, err := myAMQP.Producer(producerOptions)
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
