package main

import (
	"flag"
	"log"

	"github.com/DavidHernandez21/rabbitmq/keyboardinterrup"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	exchangeName := flag.String("exchange", "logs", "The exchange name")
	consumerName := flag.String("consumer", "my-consumer", "The consumer name")

	flag.Parse()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		*exchangeName, // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,        // queue name
		"",            // routing key
		*exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name,        // queue
		*consumerName, // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s - tag: %d - id: %s - resend: %t", d.Body, d.DeliveryTag, d.MessageId, d.Redelivered)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")

	keyboardinterrup.Listening()

	ch.Cancel(*consumerName, false)

}
