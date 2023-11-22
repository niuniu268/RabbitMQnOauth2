package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/" // Update with your RabbitMQ server URL
	queueName   = "example_queue"
	corrId      = "12345"
	replyQueue  = "reply_queue"
)

var oauthConfig = &oauth2.Config{
	ClientID:     "000000",
	ClientSecret: "999999",
	Endpoint: oauth2.Endpoint{
		TokenURL: "http://localhost:9096/token",
	},
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	// Obtain OAuth2 token
	token, err := getOAuthToken()
	failOnError(err, "Failed to obtain OAuth2 token")
	log.Printf("token %v \n", token.AccessToken)

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
			Headers: amqp.Table{
				"Authorization": token.AccessToken,
			},
			Body: []byte(body),
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

// Function to obtain OAuth2 token
func getOAuthToken() (*oauth2.Token, error) {
	// Build the request URL
	url := "http://localhost:9096/token?grant_type=client_credentials&client_id=000000&client_secret=999999&scope=read"

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
