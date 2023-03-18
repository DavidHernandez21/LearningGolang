package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Note about queue size
	// If all the workers are busy, your queue can fill up.
	// You will want to keep an eye on that, and maybe add more workers,
	//  or have some other strategy.
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//This tells RabbitMQ not to give more than one message to a worker at a time.
	//Or, in other words, don't dispatch a new message to a worker until it has processed and acknowledged the previous one.
	//Instead, it will dispatch it to the next worker that is not still busy.
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// By default, RabbitMQ will send each message to the next consumer,
	//  in sequence. On average every consumer will get the same number of messages.
	//  This way of distributing messages is called round-robin
	msgs, err := ch.Consume(
		q.Name,        // queue
		"my-consumer", // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		var err error
		var dotCount int
		for d := range msgs {
			log.Printf("Received a message: %s - tag: %d - id: %s - resend: %t", d.Body, d.DeliveryTag, d.MessageId, d.Redelivered)
			dotCount = bytes.Count(d.Body, []byte("."))
			// t := time.Duration(dotCount)
			time.Sleep(time.Duration(dotCount) * time.Second)
			log.Printf("Done")
			err = d.Ack(false)
			failOnError(err, "Failed to acknowledge a message")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	log.Printf("received signal: %v...cancelling the channel", s)

	ch.Cancel("my-consumer", false)

}

func consumeMessages(ctx context.Context, message amqp.Delivery) error {
	var err error

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

		log.Printf("Received a message: %s - tag: %d - id: %s - resend: %t", message.Body, message.DeliveryTag, message.MessageId, message.Redelivered)
		dotCount := bytes.Count(message.Body, []byte("."))
		// t := time.Duration(dotCount)
		time.Sleep(time.Duration(dotCount) * time.Second)
		log.Printf("Done")
		err = message.Ack(false)
		// if err != nil {
		// 	return err
		// }

		return err
	}
}
