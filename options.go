package myamqp

import (
	"errors"

	"github.com/rabbitmq/amqp091-go"
)

var (
	ErrOptionsCannotBeNil         = errors.New("options cannot be nil")
	ErrExchangeOptionsCannotBeNil = errors.New("exchange options cannot be nil")
	ErrQueueOptionsCannotBeNil    = errors.New("queue options cannot be nil")
)

// ExchangeOptions represents options for configuring an exchange.
type ExchangeOptions struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	args       amqp091.Table
}

// NewExchangeOptions creates a new ExchangeOptions with the given name and kind.
func NewExchangeOptions(name, kind string) *ExchangeOptions {
	return &ExchangeOptions{
		name: name,
		kind: kind,
	}
}

// WithDurable sets the durable on the ExchangeOptions.
func (eo *ExchangeOptions) WithDurable(durable bool) *ExchangeOptions {
	eo.durable = durable
	return eo
}

// WithAutoDelete sets the autoDelete on the ExchangeOptions.
func (eo *ExchangeOptions) WithAutoDelete(autoDelete bool) *ExchangeOptions {
	eo.autoDelete = autoDelete
	return eo
}

// WithInternal sets the internal on the ExchangeOptions.
func (eo *ExchangeOptions) WithInternal(internal bool) *ExchangeOptions {
	eo.internal = internal
	return eo
}

// WithNoWait sets the noWait on the ExchangeOptions.
func (eo *ExchangeOptions) WithNoWait(noWait bool) *ExchangeOptions {
	eo.noWait = noWait
	return eo
}

// WithArgs sets the args on the ExchangeOptions.
func (eo *ExchangeOptions) WithArgs(args amqp091.Table) *ExchangeOptions {
	eo.args = args
	return eo
}

// QueueOptions represents options for configuring a queue.
type QueueOptions struct {
	name       string
	routingKey string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp091.Table
}

// NewQueueOptions creates a new QueueOptions with the given name.
func NewQueueOptions(name string) *QueueOptions {
	return &QueueOptions{
		name: name,
	}
}

// WithRoutingKey sets the routingKey on the QueueOptions.
func (qo *QueueOptions) WithRoutingKey(routingKey string) *QueueOptions {
	qo.routingKey = routingKey
	return qo
}

// WithDurable sets the durable on the QueueOptions.
func (qo *QueueOptions) WithDurable(durable bool) *QueueOptions {
	qo.durable = durable
	return qo
}

// WithAutoDelete sets the autoDelete on the QueueOptions.
func (qo *QueueOptions) WithAutoDelete(autoDelete bool) *QueueOptions {
	qo.autoDelete = autoDelete
	return qo
}

// WithExclusive sets the exclusive on the QueueOptions.
func (qo *QueueOptions) WithExclusive(exclusive bool) *QueueOptions {
	qo.exclusive = exclusive
	return qo
}

// WithNoWait sets the noWait on the QueueOptions.
func (qo *QueueOptions) WithNoWait(noWait bool) *QueueOptions {
	qo.noWait = noWait
	return qo
}

// WithArgs sets the args on the QueueOptions.
func (qo *QueueOptions) WithArgs(args amqp091.Table) *QueueOptions {
	qo.args = args
	return qo
}

// ConsumerOptions represents options for configuring a consumer.
type ConsumerOptions struct {
	name         string
	autoAck      bool
	exclusive    bool
	noLocal      bool
	noWait       bool
	args         amqp091.Table
	exchangeOpts *ExchangeOptions
	queueOpts    *QueueOptions
}

// NewConsumerOptions creates a new ConsumerOptions with the given name, ExchangeOptions, and QueueOptions.
func NewConsumerOptions(name string, exchangeOptions *ExchangeOptions, queueOptions *QueueOptions) *ConsumerOptions {
	return &ConsumerOptions{
		name:         name,
		exchangeOpts: exchangeOptions,
		queueOpts:    queueOptions,
	}
}

// WithName sets the name on the ConsumerOptions.
func (co *ConsumerOptions) WithName(name string) *ConsumerOptions {
	co.name = name
	return co
}

// WithAutoAck sets the autoAck on the ConsumerOptions.
func (co *ConsumerOptions) WithAutoAck(autoAck bool) *ConsumerOptions {
	co.autoAck = autoAck
	return co
}

// WithExchangeOptions sets the ExchangeOptions on the ConsumerOptions.
func (co *ConsumerOptions) WithExchangeOptions(exchangeOpts *ExchangeOptions) *ConsumerOptions {
	co.exchangeOpts = exchangeOpts
	return co
}

// WithQueueOptions sets the QueueOptions on the ConsumerOptions.
func (co *ConsumerOptions) WithQueueOptions(queueOpts *QueueOptions) *ConsumerOptions {
	co.queueOpts = queueOpts
	return co
}

// WithExclusive sets the exclusive on the ConsumerOptions.
func (co *ConsumerOptions) WithExclusive(exclusive bool) *ConsumerOptions {
	co.exclusive = exclusive
	return co
}

// WithNoLocal sets the noLocal on the ConsumerOptions.
func (co *ConsumerOptions) WithNoLocal(noLocal bool) *ConsumerOptions {
	co.noLocal = noLocal
	return co
}

// WithNoWait sets the noWait on the ConsumerOptions.
func (co *ConsumerOptions) WithNoWait(noWait bool) *ConsumerOptions {
	co.noWait = noWait
	return co
}

// WithArgs sets the args on the ConsumerOptions.
func (co *ConsumerOptions) WithArgs(args amqp091.Table) *ConsumerOptions {
	co.args = args
	return co
}

// ProducerOptions represents options for configuring a producer.
type ProducerOptions struct {
	exchangeOpts *ExchangeOptions
	queueOpts    *QueueOptions
}

// NewProducerOptions creates a new ProducerOptions with the given ExchangeOptions.
func NewProducerOptions(exchangeOptions *ExchangeOptions) *ProducerOptions {
	return &ProducerOptions{
		exchangeOpts: exchangeOptions,
	}
}

// WithQueueOptions sets the QueueOptions on the ProducerOptions.
func (po *ProducerOptions) WithQueueOptions(queueOpts *QueueOptions) *ProducerOptions {
	po.queueOpts = queueOpts
	return po
}

// Qos represents options for configuring Qos.
type Qos struct {
	prefetchCount int
	prefetchSize  int
	global        bool
}

// NewQos creates a new Qos with the given prefetchCount, prefetchSize, and global.
func NewQos(prefetchCount, prefetchSize int, global bool) *Qos {
	return &Qos{
		prefetchCount: prefetchCount,
		prefetchSize:  prefetchSize,
		global:        global,
	}
}

// PrefetchCount returns the prefetchCount on the Qos.
func (q *Qos) PrefetchCount() int {
	return q.prefetchCount
}

// PrefetchSize returns the prefetchSize on the Qos.
func (q *Qos) PrefetchSize() int {
	return q.prefetchSize
}

// Global returns the global on the Qos.
func (q *Qos) Global() bool {
	return q.global
}
