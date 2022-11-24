package main

import (
	"log"

	"github.com/go-rabbitmq-sample/shared"
		amqp "github.com/rabbitmq/amqp091-go"

)

var (
	err error
)

func main() {
	// setup rabbitmq
	connection, err := amqp.Dial(shared.RABBITMQ_SERVER_URL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	// open a channel to the instance over the created connection
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// subscribe to service
	messages, err := channel.Consume(
		shared.SERVICE_ONE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("error subscribing to message - %v", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// make a channel to receive messgaes into infinite loop
	forever := make(chan bool)

	go func() {
		for message := range messages {
			log.Printf(" > Received message: %s\n", message.Body)
		}
	}()

	<-forever
}
