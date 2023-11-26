package main

import (
	proto "DS2023-Auction/grpc"
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"

	"slices"
	"strconv"
	"strings"

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
	port     = flag.Int("port", 0, "portnumber")
	joinport = flag.Int("joinPort", 0, "join port of primary replica")
)

var amountsFromBidders = make(map[string]int)
var biddersFromAmounts = make(map[int]string)
var max int

var idsSlice []string

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
		register(node)
		go waitForInfo(node)
	}

	timeLimit := time.After(5 * time.Minute)
	<-timeLimit
	log.Printf("AUCTION HAS ENDED\n")
	os.Exit(0)

	for {

	}
}

func register(n *Node) {
	serverConnection, _ := connect()

	// Ask the server for the time
	registerSuccess, err := serverConnection.RegisterBackup(context.Background(), &proto.RegisterRequest{
		BackupId: strconv.Itoa(n.id),
	})

	if err != nil {
		log.Fatal(err.Error()) //if error is not null, print the error message
	} else {
		log.Printf("Registered new replica id=%s", registerSuccess.BackupId)
	} //else if there are no errors, then print "Server <server_name> says the time is" + some_time_stamp

}

func election(idSlice []string) string {
	var maxId string
	var tempSlice []int
	for _, id := range idSlice {
		intId, _ := strconv.Atoi(id)
		tempSlice = append(tempSlice, intId)
	}
	maxId = strconv.Itoa(slices.Max(tempSlice))
	return maxId
}

func waitForInfo(n *Node) {
	connection, err := connect() //Auction Client
	waitchannel := make(chan struct{})

	if err != nil {
		log.Fatal("Could not connect inf waitForInfo")
	} else {
		for {
			stream, err := connection.AskForUpdate(context.Background())

			if err != nil {
				log.Fatal(err.Error())
			}

			updateResult, err := stream.Recv()
			if err == io.EOF {
				close(waitchannel)
				return
			}

			if err != nil {
				close(waitchannel)
				if len(idsSlice) > 0 {
					promotedReplicaId := election(idsSlice)

					log.Printf("Primary replica has crashed...initating new election")
					converted, _ := strconv.Atoi(promotedReplicaId)
					if n.id == converted {
						//print om back up har vundet
						log.Printf("Backup replica with id %d has won the election", n.id)
						var newIdsSlice []string

						for _, id := range idsSlice {
							if id != promotedReplicaId {
								newIdsSlice = append(newIdsSlice, id)
							}
						}

						idsSlice = newIdsSlice

						n.port = 5454
						*ownId = 1
						*joinport = 0
						nodeServer(n)
					} else {
						//print om at backup er backup
						log.Printf("Backup replica with id %d continues as backup", n.id)
						waitForInfo(n)
					}
				}

			} else {
				amountsFromBidders[updateResult.BidderId] = int(updateResult.Amount) //Gives amount value from bidder id
				biddersFromAmounts[int(updateResult.Amount)] = updateResult.BidderId //Gives bidder id from amount
				max = int(updateResult.Amount)
				idsSlice = strings.Split(updateResult.BackupIds, ",")
			}
		}
	}
	<-waitchannel

}

func (c *Node) AskForUpdate(str proto.Auction_AskForUpdateServer) error {
	message := "Primary replica with id " + strconv.Itoa(*ownId) + " has sent updated state to backup"
	//log.Print(message)

	state := &proto.SendState{
		BidderId:       biddersFromAmounts[max],
		Amount:         int64(max),
		SuccessMessage: message,
		BackupIds:      strings.Join(idsSlice[:], ","),
	}

	if err := str.Send(state); err != nil {
		log.Printf("send error %v", err)
	}
	return nil
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

func (backupNode *Node) RegisterBackup(ctx context.Context, in *proto.RegisterRequest) (*proto.Registered, error) {
	idsSlice = append(idsSlice, in.BackupId)
	log.Printf("[%d] Node %s registered as backup replica.", *ownId, in.BackupId)
	return &proto.Registered{BackupId: in.BackupId}, nil
}

func (c *Node) Bidding(ctx context.Context, in *proto.BidAmount) (*proto.BidSuccess, error) {
	log.Printf("Bidder with ID %s has bidded with amount %d\n", in.BidderId, int(in.Amount))
	//Saving bidded amount to map
	var maxAmount = 0
	//Initialize value variable and run through amountFromBidders map
	for _, input := range amountsFromBidders {
		//Compare and find highest amount
		if input > maxAmount {
			maxAmount = input
		}
	}

	if maxAmount < int(in.Amount) {
		amountsFromBidders[in.BidderId] = int(in.Amount) //gives an amount from bidder id
		biddersFromAmounts[int(in.Amount)] = in.BidderId //gives bidder from certain amount
		max = int(in.Amount)
	} else {
		log.Printf("Client with id %s has bid with invalid amount. Bid %d must be higher than %d\n", in.BidderId, in.Amount, maxAmount)
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
		BidderId: biddersFromAmounts[maxAmount],
		Amount:   int64(amountsFromBidders[winner]),
	}, nil
}
