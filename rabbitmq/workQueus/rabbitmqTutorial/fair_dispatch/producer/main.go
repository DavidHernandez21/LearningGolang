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

	flag.Parse()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// err = ch.Publish(
	// 	"",     // exchange
	// 	q.Name, // routing key
	// 	false,  // mandatory
	// 	false,
	// 	amqp.Publishing{
	// 		// Marking messages as persistent doesn't fully guarantee that a message won't be lost.
	// 		//Although it tells RabbitMQ to save the message to disk,
	// 		//there is still a short time window when RabbitMQ has accepted a message and hasn't saved it yet
	// 		DeliveryMode: amqp.Persistent,
	// 		ContentType:  "text/plain",
	// 		Body:         []byte(body),
	// 	})
	ctx, cancel := context.WithTimeout(context.Background(), *ctxTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	body := *mex
	repsValue := *reps
	for i := 0; i < repsValue; i++ {
		g.Go(func() error {
			return publishmq.PublishQueue(ctx, ch, body, q.Name)
		})
	}

	err = g.Wait()
	failOnError(err, "Failed to publish a message")

	// os.Args[2] | 4
	log.Printf(" [x] Sent %s - %d times", body, repsValue)
}

// func bodyFrom(args []string) string {
// 	var s string
// 	if (len(args) < 2) || os.Args[1] == "" {
// 		s = "hello"
// 	} else {
// 		s = strings.Join(args[1:], " ")
// 	}
// 	return s
// }
