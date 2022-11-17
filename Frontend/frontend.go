package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"strconv"
	"time"

	proto "github.com/MadsRoager/AuctionReplication/proto"
	"google.golang.org/grpc"
)

type Frontend struct {
	proto.UnimplementedFrontendServer
	serverNodePorts  []int32
	serverConnection proto.ServerClient
}

var frontendPort = flag.Int("serverPort", 8080, "server port number")

func main() {
	frontend := &Frontend{
		serverNodePorts:  make([]int32, 4),
		serverConnection: getServerConnection(),
	}
	go startFrontEnd(frontend)

	for {
		fmt.Println("frontend started")
		time.Sleep(100 * time.Second)
	}
}

func startFrontEnd(frontend *Frontend) {
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

func getServerConnection() proto.ServerClient { // Hard coded port 8081
	conn, err := grpc.Dial(":"+strconv.Itoa(8081), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Could not dial server")
	}
	return proto.NewServerClient(conn)
}

func (frontend *Frontend) Result(ctx context.Context, in *proto.Void) (*proto.BidResult, error) {
	// send GetHighestBid to all servernodes
	return frontend.serverConnection.GetHighestBid(ctx, in)

	// when receiving a response, return it and wait for the other serverNodes to reply

}

func (frontend *Frontend) Bid(ctx context.Context, in *proto.BidRequest) (*proto.Ack, error) {
	// send UpdateHighestBid to all servernodes
	return frontend.serverConnection.UpdateHighestBid(ctx, in)

	// when receiving a response, return it and wait for the other serverNodes to reply

}
