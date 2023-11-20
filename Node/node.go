package main;

import {
proto "mod/grpc"




}


type Node struct {
id int
port int
proto.UnimplementedAuctionServer
}


var (
ownId = flag.Int("id", 0, "node id")
port = flag.Int("port", 5454, "portnumber")
)

highestAmount = make(map[int]int)

func main() {

	flag.Parse()


node := &Node {
id: *ownId,
port: *port


go nodeServer(n *Node)


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
	log.Printf("Bidder with ID %d has bidded\n", in.BidderId)
	return &proto.BidSuccess{
	SuccessMessage: "Bid has been accepted",
	}, nil	
}

func (c *Node) AskForResult(ctx context.Context, in *proto.ResultRequest) (*proto.Result, error) {
	log.Print("Bidder with ID %d has asked for Auction result\n", in.BidderId)
	return &proto.Result{
	bidderId: in.BidderId,
	amount: in.Amount,
	}
}


} 