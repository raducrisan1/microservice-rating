package main

import (
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
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	failOnError(err, "Could not connect to StockInfo gRPC server")
	defer conn.Close()

	stockInfoClient := stockinfo.NewStockInfoServiceClient(conn)
	rabbitQueue, rabbitChannel, err := setupRabbitMqWriter()
	failOnError(err, "Could not setup rabbitMq link to write data")
	defer rabbitChannel.Close()

	ticker := time.Tick(time.Second * 5)
	pulse := setupsignal(ticker)
	for {
		startTime, _ := time.Parse(time.RFC3339, "2018-11-10 09:30Z")
		endTime, _ := time.Parse(time.RFC3339, "2018-11-10 10:00Z")
		req := &stockinfo.StockInfoRequest{
			Stockname:  "NVDA",
			Start:      startTime.Unix(),
			End:        endTime.Unix(),
			Resolution: 300}

		res, err := stockInfoClient.StockInfo(nil, req)
		rating := stockinfo.StockRating{
			Rating:         int32(rand.Intn(5) + 1),
			Islongposition: rand.Intn(2) > 0,
			Stockname:      res.Stockname,
			Timestamp:      time.Now().Unix()}

		if err != nil {
			fmt.Printf("Could not obtain data from the gRPC service StockInfo: %v", err)
			continue
		}
		exitflag := <-pulse
		if exitflag == 1 {
			break
		}
		err = sendMessage(&rating, rabbitQueue, rabbitChannel)
	}
}
