package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strconv"

	proto "github.com/MadsRoager/AuctionReplication/proto"
	"google.golang.org/grpc"
)


type Frontend struct {
	proto.UnimplementedFrontendServer
	serverNodePorts		[]int32
}

var frontendPort = flag.Int("serverPort", 8080, "server port number")

func main() {

	frontend := &Frontend{
		serverNodePorts:    make([]int32, 4),
	}

	grpcServer := grpc.NewServer()

	lister, err := net.Listen("tcp", ":"+strconv.Itoa(*frontendPort))

	if err != nil {
		log.Fatalln("could not start listener")
	}

	proto.RegisterFrontendServer(grpcServer, frontend)
	serverError := grpcServer.Serve(lister)

	if serverError != nil {
		log.Fatalln("could not start server")
	}

}

func (f *Frontend) Result(ctx context.Context, in *proto.Void) (*proto.BidResult, error) {
	// send GetHighestBid to all servernodes

	// when receiving a response, return it and wait for the other serverNodes to reply
	
}

func (f *Frontend) Bid(ctx context.Context, in *proto.BidRequest) (*proto.Ack, error) {
	// send UpdateHighestBid to all servernodes

	// when receiving a response, return it and wait for the other serverNodes to reply
	
}
