package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

type mexCount struct {
	count int
}

func joinBuilder(grow int, words ...string) (string, error) {

	// sb.Reset()
	var sb strings.Builder
	sb.Grow(90)

	for _, word := range words {
		_, err := sb.WriteString(word)

		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func loopThroughNotifyPublish(ch *amqp.Channel, mexCount *mexCount, wg *sync.WaitGroup) {

	defer wg.Done()

	log.Println("collecting messages from notifyPublish channel")

	notifyConfirmation := make(chan amqp.Confirmation, mexCount.count)

	notifyPublish := ch.NotifyPublish(notifyConfirmation)

	_, ok := <-notifyPublish

	if ok {
		close(notifyPublish)
	}

	for msg := range notifyPublish {
		if !msg.Ack {
			log.Printf("Failed to confirm delivery of message %d\n", msg.DeliveryTag)
		}
		log.Println(msg)
	}

	log.Println("done with NotifyPublish")

}

func loopThroughNotifyReturn(ch *amqp.Channel, mexCount *mexCount, wg *sync.WaitGroup) {

	defer wg.Done()

	log.Println("collecting messages from notifyReturn channel")

	amqpReturn := make(chan amqp.Return, mexCount.count)

	notifyReturn := ch.NotifyReturn(amqpReturn)

	_, ok := <-notifyReturn

	if ok {
		close(notifyReturn)
	}

	for ret := range notifyReturn {
		log.Printf("notifyReturn: %v\n", ret)
	}

	log.Println("done with NotifyReturn")

}

func listenForKeyboardInterrupt(done chan struct{}) {
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	log.Printf("received signal: %v...closing publishing channel", s)
	done <- struct{}{}
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

	if err := ch.Confirm(false); err != nil {
		log.Fatalf("Channel could not be put into confirm mode: %s", err)
	}

	q, err := ch.QueueDeclare("gophers", false, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	//debug only
	fmt.Println(q)

	done := make(chan struct{})

	go listenForKeyboardInterrupt(done)

	payload, err := joinBuilder(90, "Forza Roma ", "Daje!!")

	if err != nil {
		log.Fatalf("error strings builder: %s", err.Error())
	}

	mexCount := mexCount{}

	// go runGoFunctions(done, ch, &mexCount, &wg)

loop:
	for {

		if err != nil {
			log.Fatalf("error strings builder: %v", err)
		}
		log.Println("publishing message")
		err = ch.Publish("", q.Name, true, false,
			amqp.Publishing{
				Headers:     nil,
				ContentType: "text/plain",
				Body:        []byte(payload),
				Timestamp:   time.Now(),
			})

		if err != nil {
			log.Println(err)
			break loop
		}

		mexCount.count += 1

		//wait 2 seconds until send another message
		time.Sleep(2 * time.Second)

		select {
		case <-done:
			log.Println("received from done...exiting the loop")
			break loop
		default:
		}
	}

	log.Printf("total messages published: %d", mexCount.count)

	var wg sync.WaitGroup

	wg.Add(2)

	go loopThroughNotifyPublish(ch, &mexCount, &wg)
	go loopThroughNotifyReturn(ch, &mexCount, &wg)

	go func() {
		wg.Wait()
		log.Println("exiting program")
	}()
}
