package main

import (
	proto "DS2023-Auction/grpc"
	"context"
	"flag"
	"log"
	"net"
	"slices"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
	id   int
	port int
	proto.UnimplementedAuctionServer
}

var (
	ownId    = flag.Int("id", 0, "node id")
	port     = flag.Int("port", 5454, "portnumber")
	joinport = flag.Int("joinPort", 0, "join port of primary replica")
)

var amountsFromBidders = make(map[int]int)
var biddersFromAmounts = make(map[int]int)
var max int
var backupIDs []int64

func main() {

	flag.Parse()

	node := &Node{
		id:   *ownId,
		port: *port,
	}

	if int(*ownId) == 1 {
		// Primary replica manager code
		go nodeServer(node)
	} else {
		//Backup replica manager code
		go waitForInfo(node)
	}

	for {

	}
}

func election(idSlice []int64) (int64, error) {
	return slices.Max(idSlice), nil

}

func waitForInfo(n *Node) {
	connection, err := connect() //Auction Client

	if err != nil {
		log.Fatal("hej :(")

	} else {
		for {

			//Instantiate variable updateResult of type proto.SendState.
			updateResult, err := connection.AskForUpdate(context.Background(), &proto.RequestState{
				OwnId: int64(*ownId), //backup replica ID
			})

			if err != nil {
				//Before it was log.Fatalf which you print the error and afterwards quit the program.
				log.Print("Unable to connect to primary replica -> Initiating election")

				//var electionResult = election(updateResult.BackupsIds)
				log.Print(updateResult.BackupId)
				//BackupId represents the winner from election returned from server side.
				if updateResult.BackupId == int64(*ownId) {
					//If true, the backup node becomes the new primary node
					//*ownId = 1
					//*port = 5454
					log.Print("jeg er inde i if ")
				} else {
					//If false then the backup continues being a backup by requesting connection to primary
					//go waitForInfo(n)
					log.Print("This shouldn't occur")
				}

			} else {
				//Store values in respective variables

				amountsFromBidders[int(updateResult.BidderId)] = int(updateResult.Amount) //Gives amount value from bidder id
				biddersFromAmounts[int(updateResult.Amount)] = int(updateResult.BidderId) //Gives bidder id from amount
				max = int(updateResult.Amount)
				log.Printf("Backup replica with id %d has received updated state\n", *ownId)
				log.Printf("Backup received amount %d\n", strconv.Itoa(max))
				//log.Print(updateResult.BackupsIds)
			}
		}
	}
}

func (c *Node) AskForUpdate(ctx context.Context, in *proto.RequestState) (*proto.SendState, error) {
	message := "Primary replica with id " + strconv.Itoa(*ownId) + "has sent updated state to backup with id " + strconv.Itoa(int(in.OwnId))

	//Append unique ID's to backupIDs slice
	for _, id := range backupIDs {
		if id != int64(in.OwnId) {
			backupIDs = append(backupIDs, int64(in.OwnId))
		}
		break
	}
/*
	electionwinner := int64(1)

	if len(backupIDs) != 0 {
		winner, err := election(backupIDs)

		if err != nil {
			log.Print("Error in election call")
		} else {
			electionwinner = winner
		}
	}*/

	return &proto.SendState{
		BidderId:       int64(biddersFromAmounts[max]),
		Amount:         int64(max),
		SuccessMessage: message,
		BackupId:       int64(2),
	}, nil
}

func connect() (proto.AuctionClient, error) {
	conn, err := grpc.Dial(":"+strconv.Itoa(*joinport), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *joinport)
	} else {
		log.Printf("Backup with id %d has connected to primary at port %d", *ownId, *joinport)
	}
	return proto.NewAuctionClient(conn), nil
}

func nodeServer(n *Node) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(n.port))

	//if err (error) is not null, then print "Could not create the server" + err
	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", n.port)
	//if there are no errors, then print "Started server at port" + the designated serverport

	// Register the grpc server and serve its listener
	proto.RegisterAuctionServer(grpcServer, n)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (c *Node) Bidding(ctx context.Context, in *proto.BidAmount) (*proto.BidSuccess, error) {
	//Print for debugging
	log.Printf("Bidder with ID %d has bidded with amount %d\n", in.BidderId, int(in.Amount))
	//Saving bidded amount to map
	var maxAmount = 0
	//Initialize value variable and run through highestamount map
	for _, input := range amountsFromBidders {
		//Compare and find highest amount
		if input > maxAmount {
			maxAmount = input
		}
	}

	if maxAmount < int(in.Amount) {
		amountsFromBidders[int(in.BidderId)] = int(in.Amount) //gives an amount from bidder id
		biddersFromAmounts[int(in.Amount)] = int(in.BidderId) //gives bidder from certain amount
		max = int(in.Amount)
		//log.Printf("int amount: %d\n", int(in.Amount)) keeps tracks on the highest bid
		//log.Printf("int maxAmount: %d\n", maxAmount)
	} else {
		log.Printf("Client with id %d has bid with invalid amount. Bid %d must be higher than %d\n", in.BidderId, in.Amount, maxAmount)
		return &proto.BidSuccess{
			SuccessMessage: "Your bid has not been accepted, please bid a higher amount\n",
		}, nil
	}

	//return success message
	return &proto.BidSuccess{
		SuccessMessage: "Bid has been accepted",
	}, nil
}

func (c *Node) AskForResult(ctx context.Context, in *proto.ResultRequest) (*proto.Result, error) {
	var maxAmount = 0
	//Initialize value variable and run through highestamount map
	for _, input := range amountsFromBidders {
		//Compare and find highest amount
		if input > maxAmount {
			maxAmount = input
		}
	}

	var winner = biddersFromAmounts[maxAmount]

	return &proto.Result{
		BidderId: int64(biddersFromAmounts[maxAmount]),
		Amount:   int64(amountsFromBidders[winner]),
	}, nil
}
