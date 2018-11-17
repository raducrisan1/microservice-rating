package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/raducrisan1/microservice-rating/stockinfo"
	"google.golang.org/grpc"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(msg)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:3001", grpc.WithInsecure())
	failOnError(err, "Could not connect to StockInfo gRPC server")
	defer conn.Close()

	stockInfoClient := stockinfo.NewStockInfoServiceClient(conn)

	rabbitQueue, rabbitChannel, err := setupRabbitMqWriter()
	failOnError(err, "Could not setup rabbitMq link to write data")
	defer rabbitChannel.Close()

	//every 5 seconds, a call to stockinfo microservice is made.

	impulse := make(chan int, 2)

	go func() {
		impulse <- 1
	}()

	ticker := time.Tick(time.Second * 3)
	osstop := setupsignal()
	stop := false
	for !stop {
		select {
		case <-ticker:
			impulse <- 1
		case <-impulse:
			startTime, _ := time.Parse(time.RFC3339, "2018-11-10 09:30Z")
			endTime, _ := time.Parse(time.RFC3339, "2018-11-10 10:00Z")
			req := new(stockinfo.StockInfoRequest)
			req.Stockname = "NVDA"
			req.Start = startTime.Unix()
			req.End = endTime.Unix()
			req.Resolution = 300
			res, err := stockInfoClient.StockInfo(context.Background(), req)
			if err != nil {
				fmt.Printf("An error occurred receiving data from gRPC stockinfo: %s\n", err)
				continue
			}

			//augment the stockinfo response with a random rating and send it to RabbitMQ for further processing
			rating := new(stockinfo.StockRating)
			rating.Rating = int32(rand.Intn(5) + 1)
			rating.Islongposition = rand.Intn(2) > 0
			rating.Stockname = res.Stockname
			rating.Timestamp = time.Now().Unix()

			if err != nil {
				fmt.Printf("Could not obtain data from the gRPC service StockInfo: %v", err)
				continue
			}

			err = sendMessage(rating, rabbitQueue, rabbitChannel)
		case <-osstop:
			stop = true
			fmt.Println("\nNode stop has been requested")
		}
	}

	fmt.Println("The node has been stopped")
}
