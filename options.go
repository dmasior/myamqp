package myamqp

import "github.com/rabbitmq/amqp091-go"

type ExchangeOptions struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	args       amqp091.Table
}

func NewExchangeOptions(name, kind string) *ExchangeOptions {
	return &ExchangeOptions{
		name: name,
		kind: kind,
	}
}

func (eo *ExchangeOptions) WithDurable(durable bool) *ExchangeOptions {
	eo.durable = durable
	return eo
}

func (eo *ExchangeOptions) WithAutoDelete(autoDelete bool) *ExchangeOptions {
	eo.autoDelete = autoDelete
	return eo
}

func (eo *ExchangeOptions) WithInternal(internal bool) *ExchangeOptions {
	eo.internal = internal
	return eo
}

func (eo *ExchangeOptions) WithNoWait(noWait bool) *ExchangeOptions {
	eo.noWait = noWait
	return eo
}

func (eo *ExchangeOptions) WithArgs(args amqp091.Table) *ExchangeOptions {
	eo.args = args
	return eo
}

type QueueOptions struct {
	name       string
	routingKey string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp091.Table
}

func NewQueueOptions(name string) *QueueOptions {
	return &QueueOptions{
		name: name,
	}
}

func (qo *QueueOptions) WithRoutingKey(routingKey string) *QueueOptions {
	qo.routingKey = routingKey
	return qo
}

func (qo *QueueOptions) WithDurable(durable bool) *QueueOptions {
	qo.durable = durable
	return qo
}

func (qo *QueueOptions) WithAutoDelete(autoDelete bool) *QueueOptions {
	qo.autoDelete = autoDelete
	return qo
}

func (qo *QueueOptions) WithExclusive(exclusive bool) *QueueOptions {
	qo.exclusive = exclusive
	return qo
}

func (qo *QueueOptions) WithNoWait(noWait bool) *QueueOptions {
	qo.noWait = noWait
	return qo
}

func (qo *QueueOptions) WithArgs(args amqp091.Table) *QueueOptions {
	qo.args = args
	return qo
}

type ConsumerOptions struct {
	name         string
	autoAck      bool
	exclusive    bool
	noLocal      bool
	noWait       bool
	args         amqp091.Table
	exchangeOpts ExchangeOptions
	queueOpts    QueueOptions
}

func NewConsumerOptions(name string) *ConsumerOptions {
	return &ConsumerOptions{
		name:         name,
		exchangeOpts: ExchangeOptions{},
		queueOpts:    QueueOptions{},
	}
}

func (co *ConsumerOptions) WithAutoAck(autoAck bool) *ConsumerOptions {
	co.autoAck = autoAck
	return co
}

func (co *ConsumerOptions) WithExchangeOptions(exchangeOpts ExchangeOptions) *ConsumerOptions {
	co.exchangeOpts = exchangeOpts
	return co
}

func (co *ConsumerOptions) WithQueueOptions(queueOpts QueueOptions) *ConsumerOptions {
	co.queueOpts = queueOpts
	return co
}

func (co *ConsumerOptions) WithExclusive(exclusive bool) *ConsumerOptions {
	co.exclusive = exclusive
	return co
}

func (co *ConsumerOptions) WithNoLocal(noLocal bool) *ConsumerOptions {
	co.noLocal = noLocal
	return co
}

func (co *ConsumerOptions) WithNoWait(noWait bool) *ConsumerOptions {
	co.noWait = noWait
	return co
}

func (co *ConsumerOptions) WithArgs(args amqp091.Table) *ConsumerOptions {
	co.args = args
	return co
}
