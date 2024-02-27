package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/brenoassp/crud-go/adapters/messageBroker"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	connectionString string
}

func New(
	connectionString string,
) Client {
	if connectionString == "" {
		connectionString = "amqp://guest:guest@localhost:5672/"
	}

	return Client{
		connectionString: connectionString,
	}
}

func (c Client) Publish(
	ctx context.Context,
	exchange, routingKey string,
	mandatory, immediate bool,
	msg messageBroker.Message,
) error {
	conn, err := amqp.Dial(c.connectionString)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		routingKey, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = ch.PublishWithContext(
		ctx,
		exchange,  // exchange
		q.Name,    // routing key
		mandatory, // mandatory
		immediate, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  msg.ContentType,
			Body:         msg.Body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Println(ctx, "Published message", map[string]interface{}{
		"exchange":   exchange,
		"routingKey": routingKey,
		"msg":        msg,
	})

	return nil
}

func (c Client) Consume(
	ctx context.Context,
	exchange, queueName, exchangeType string,
) (*messageBroker.Consumer, error) {
	consumer := &messageBroker.Consumer{
		Done:    make(chan error),
		Message: make(chan messageBroker.Message),
	}

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("create-client-consumer")
	conn, err := amqp.DialConfig(c.connectionString, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	consumer.Close = conn.Close

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %s", err)
	}

	queue, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	deliveries, err := ch.Consume(
		queue.Name,                    // name
		"create-clients-consumer-tag", // consumerTag,
		false,                         // autoAck
		false,                         // exclusive
		false,                         // noLocal
		false,                         // noWait
		nil,                           // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %s", err)
	}

	go handle(deliveries, consumer)

	return consumer, nil
}

func handle(deliveries <-chan amqp.Delivery, consumer *messageBroker.Consumer) {
	deliveryCount := 0

	for d := range deliveries {
		deliveryCount++
		consumer.Message <- messageBroker.Message{
			ContentType: d.ContentType,
			Body:        d.Body,
			Ack:         d.Ack,
			Nack:        d.Nack,
		}
	}
}
