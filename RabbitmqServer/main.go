package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	// Declare a queue
	q, err := ch.QueueDeclare(
		"example_queue", // Queue name
		false,           // Durable
		false,           // Delete when unused
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Consume messages from the queue
	msg, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer
		false,  // Auto-acknowledge
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	failOnError(err, "Failed to register a consumer")

	// Receive and print messages
	forever := make(chan bool)

	go func() {
		for d := range msg {

			token, _ := d.Headers["Authorization"].(string)
			log.Printf("Received token %v\n", token)

			if token != "ZJE3MZRLNTKTNWRIZS0ZYZQ1LTLHYZCTNDDKZTI0ZTQ1ODLL" {
				log.Println("Invalid Oauth2 token")
				continue
			} else {
				log.Println("valid Oauth2 token")

			}
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte(".")) // Count the dots in the message
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second) // Simulate time-consuming task

			// Send a confirmation message to the client
			err = ch.Publish(
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte("Message processed successfully"),
				})
			failOnError(err, "Failed to publish confirmation to client")

			err := d.Ack(false)
			if err != nil {
				return
			}
			log.Printf("Done")
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
