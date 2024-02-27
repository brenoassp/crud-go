package messageBroker

import (
	"context"
)

type Provider interface {
	Publish(ctx context.Context, exchange, routingKey string, mandatory, immediate bool, msg Message) error
	Consume(ctx context.Context, exchange, queueName, exchangeType string) (*Consumer, error)
}

type Message struct {
	ContentType string
	Body        []byte
	Ack         func(bool) error
	Nack        func(bool, bool) error
}

type Consumer struct {
	Done    chan error
	Message chan Message
	Close   func() error
}
