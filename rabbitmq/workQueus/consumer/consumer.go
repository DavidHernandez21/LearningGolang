package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

// consumes messages from a channel if they arrived with intervals of t. Otherwise returns. the channel must be manually closed afterwards
func consumeMexwithtimeout(chDelivery <-chan amqp.Delivery, t time.Duration) {

outerloop:
	for {
		select {
		case delivery := <-chDelivery:
			log.Printf("msg: %s - %d - %d - %s - %s - %s", string(delivery.Body), delivery.DeliveryMode,
				delivery.DeliveryTag, delivery.Exchange, delivery.Expiration, delivery.RoutingKey)
		case <-time.After(t):
			log.Printf("No Messages recieved after %v....Exiting\n", t)

			break outerloop

		}
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}

	defer ch.Close()

	chDelivery, err := ch.Consume(
		"gophers",
		"my-consumer",
		true,
		false,
		false,
		false, nil)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for delivery := range chDelivery {
			log.Printf("msg: %s - %d - %d - %s - %s - %s - %t - %v", string(delivery.Body), delivery.DeliveryMode,
				delivery.DeliveryTag, delivery.Exchange, delivery.Expiration, delivery.RoutingKey, delivery.Redelivered, delivery.Timestamp)

			// delivery.Ack(false)
		}
	}()

	// go consumeMexwithtimeout(chDelivery, 10*time.Second)

	// ch.Cancel("my-consumer", false)
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	log.Printf("received signal: %v...cancelling the channel", s)

	ch.Cancel("my-consumer", false)
}
