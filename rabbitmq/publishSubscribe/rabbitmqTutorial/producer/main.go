package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/DavidHernandez21/rabbitmq/publishmq"
	"github.com/streadway/amqp"
	"golang.org/x/sync/errgroup"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	mex := flag.String("mex", "daje..", "The message body to publish")
	reps := flag.Int("reps", 3, "The number of times to publish the message")
	ctxTimeout := flag.Duration("timeout", time.Second*5, "The timeout for the context")
	exchangeName := flag.String("exchange", "logs", "The exchange name")

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

	ctx, cancel := context.WithTimeout(context.Background(), *ctxTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	body := *mex
	for i := 0; i < *reps; i++ {
		g.Go(func() error {
			return publishmq.PublishExchange(ctx, ch, *mex, *exchangeName)
		})
	}

	err = g.Wait()
	failOnError(err, "Failed to publish a message")

	// os.Args[2] | 4
	log.Printf(" [x] Sent %s - %d times", body, *reps)
}
