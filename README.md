# AuctionSystem


To run the primary replica: go run node/node.go --id=1 --port=5454 

To run backup replicas: go run node/node.go --id=<some_id> --joinPort=5454 

To run client: go run client/client.go --id=<some_id> --port=5454

When bidding: type an integer amount in the client terminal

When asking for result: type a single string in the client terminal

When crashing a replica: Close its corresponding terminal


OBS. The timeframe for the nodes have been set to 5 minutes and the client for 4 minutes, so make sure to run the client after running the primary node.
