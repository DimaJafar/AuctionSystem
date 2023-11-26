package main

import (
	proto "DS2023-Auction/grpc"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id   int
	port int
}

var (
	ownId = flag.Int("id", 0, "client id")
	port  = flag.Int("port", 5454, "portnumber")
)

func main() {
	flag.Parse()

	client := &Client{
		id: *ownId,
	}

	timeLimit := time.After(time.Minute)

	go waitForBid(client)

	<-timeLimit
	log.Printf("AUCTION HAS ENDED:\n")
	winner()
	os.Exit(0)

	for {

	}
}

func winner() {
	serverConnection, _ := connectToServer()
	result, err := serverConnection.AskForResult(context.Background(), &proto.ResultRequest{})

	if err != nil {
		log.Printf("Could not retrieve auction result") //if error is not null, print the error message
	} else {
		log.Printf("Winner is client with id %s and amount %d !!! \n", result.BidderId, result.Amount)
	}
}

func connectToServer() (proto.AuctionClient, error) {

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *port)
	} else {
		log.Printf("Connected to the server at port %d\n", *port)
	}
	return proto.NewAuctionClient(conn), nil
}

func waitForBid(c *Client) {
	serverConnection, _ := connectToServer()
	var amount string

	for {
		fmt.Scan(&amount)

		if _, err := strconv.Atoi(amount); err == nil {
			log.Printf("Client with id %d has bid with the amount: %s\n", c.id, amount)

			temp, _ := strconv.Atoi(amount)
			bidSuccess, err := serverConnection.Bidding(context.Background(), &proto.BidAmount{
				BidderId: strconv.Itoa(c.id),
				Amount:   int64(temp),
			})

			if err != nil {
				//log.Print(err.Error()) //if error is not null, print the error message
				log.Print("Auction has ended.")
			} else {
				log.Printf(bidSuccess.SuccessMessage) //bid has been accepted.
			}
		} else {
			log.Printf("Client with id %d has asked for the result of the auction\n", c.id)

			result, err := serverConnection.AskForResult(context.Background(), &proto.ResultRequest{})

			if err != nil {
				log.Printf("Could not retrieve auction result") //if error is not null, print the error message
			} else {
				log.Printf("Result is: client with ID %s is in the lead with amount %d\n", result.BidderId, result.Amount)
			}
		}
	}
}
