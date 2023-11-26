package main

import (
	proto "DS2023-Auction/grpc"
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

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

	go waitForBid(client)

	for {

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
	var amount int
	for {
		fmt.Scan(&amount)
		log.Printf("Client with id %d has bid with the amount: %d\n", c.id, amount)

		// Ask the server for the time
		bidSuccess, err := serverConnection.Bidding(context.Background(), &proto.BidAmount{
			BidderId: strconv.Itoa(c.id),
			Amount:   int64(amount),
		})

		if err != nil {
			log.Printf(err.Error()) //if error is not null, print the error message
		} else {
			log.Printf(bidSuccess.SuccessMessage) //bid has been accepted.
		} //else if there are no errors, then print "Server <server_name> says the time is" + some_time_stamp
	}

}
