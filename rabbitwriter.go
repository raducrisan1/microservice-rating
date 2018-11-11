package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/raducrisan1/microservice-rating/stockinfo"
	"github.com/streadway/amqp"
)

func setupRabbitMqWriter() (*amqp.Queue, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "could not open the connection to rabbitmq exchange")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"StockRatingData",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to declare a queue")
	return &q, ch, err
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
