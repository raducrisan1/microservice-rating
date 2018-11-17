package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/raducrisan1/microservice-rating/stockinfo"
	"github.com/streadway/amqp"
)

func newRabbitMqWriter() (*amqp.Queue, *amqp.Channel, *amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "could not open the connection to rabbitmq exchange")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a rabbitmq channel", err)
		conn.Close()
		panic("Failed to open a rabbitmq channel")
	}

	q, err := ch.QueueDeclare(
		"StockRatingData",
		true,
		false,
		false,
		false,
		nil)

	return &q, ch, conn, err
}

func sendMessage(msg *stockinfo.StockRating, q *amqp.Queue, ch *amqp.Channel) error {
	content, err := proto.Marshal(msg)
	if err != nil {
		fmt.Printf("Could not serialize the message in order to be sent to rabbitmq: %v", err)
		return err
	}
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        content})

	return err
}
