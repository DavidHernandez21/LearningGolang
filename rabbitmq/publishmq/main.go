package publishmq

import (
	"context"

	"github.com/streadway/amqp"
)

func PublishQueue(ctx context.Context, ch *amqp.Channel, body, routingKey string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := ch.Publish(
			"",         // exchange
			routingKey, // routing key
			false,      // mandatory
			false,
			amqp.Publishing{
				// Marking messages as persistent doesn't fully guarantee that a message won't be lost.
				//Although it tells RabbitMQ to save the message to disk,
				//there is still a short time window when RabbitMQ has accepted a message and hasn't saved it yet
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		return err
	}

}

func PublishExchange(ctx context.Context, ch *amqp.Channel, body, exchange string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := ch.Publish(
			exchange, // exchange
			"",       // routing key
			false,    // mandatory
			false,
			amqp.Publishing{
				// Marking messages as persistent doesn't fully guarantee that a message won't be lost.
				//Although it tells RabbitMQ to save the message to disk,
				//there is still a short time window when RabbitMQ has accepted a message and hasn't saved it yet
				// DeliveryMode: amqp.Persistent,
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		return err
	}

}
