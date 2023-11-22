package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/" // Update with your RabbitMQ server URL
	queueName   = "example_queue"
	corrId      = "12345"
	replyQueue  = "reply_queue"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	// Connect to RabbitMQ server
	conn, err := amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {

		}
	}(ch)

	// Declare a queue declare
	q, err := ch.QueueDeclare(
		"",    // Queue name
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Declare a reply queue
	//replyQ, err := ch.QueueDeclare(
	//	"",    // Queue name
	//	false, // Durable
	//	false, // Delete when unused
	//	false, // Exclusive
	//	false, // No-wait
	//	nil,   // Arguments
	//)
	//failOnError(err, "Failed to declare a reply queue")

	// Publish a message to the queue
	//body := "Hello, test!"

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	body := bodyFrom(os.Args)

	err = ch.Publish(
		"",        // Exchange
		queueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")

	// Wait for the confirmation
	for d := range msgs {
		if d.CorrelationId == corrId {
			fmt.Printf(" [x] Message Confirmed: %s\n", d.Body)
			break
		}
	}

}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
