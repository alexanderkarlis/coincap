package rabbit

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const rabbitClientCreds = "guest:guest"
const rabbitHost = "localhost"
const rabbitPort = "5672"

var rabbitPath string

func init() {
	rabbitPath = fmt.Sprintf("amqp://%s@%s:%s", rabbitClientCreds, rabbitHost, rabbitPort)
}

// SendPrices some data
func SendPrices(body []byte) {
	rabbitPath := fmt.Sprintf("amqp://%s@%s:%s", rabbitClientCreds, rabbitHost, rabbitPort)
	conn, err := amqp.Dial(rabbitPath)
	FailOnError(err, "Failed to connect to RabbitMQ!")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"coin-prices", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	FailOnError(err, "Failed to declare a queue.")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	FailOnError(err, "Failed to publish a message")
}

// SendCoinNames sends the passed in Coin data to consumer
func SendCoinNames(body []byte) {
	conn, err := amqp.Dial(rabbitPath)
	FailOnError(err, "Failed to connect to RabbitMQ!")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"coin-names", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Failed to declare a queue.")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	FailOnError(err, "Failed to publish a message")
}

// FailOnError fails on any error captured
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
