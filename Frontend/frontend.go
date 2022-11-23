package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	proto "github.com/MadsRoager/AuctionReplication/proto"
	"google.golang.org/grpc"
)

type Frontend struct {
	proto.UnimplementedFrontendServer
	serverConnection []proto.ServerClient
}

var frontendPort = flag.Int("serverPort", 8000, "server port number")

func main() {
	flag.Parse()
	frontend := &Frontend{
		serverConnection: getServerConnection(),
	}
	go startFrontEnd(frontend)

	fmt.Println("frontend started")
	for {
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

func getServerConnection() []proto.ServerClient {
	conns := make([]proto.ServerClient, 3)
	for i := 0; i < 3; i++ {
		port := 8080 + i
		conn, err := grpc.Dial(":"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("Could not dial server")
		}
		log.Printf("dialed server %d\n", port)
		conns[i] = proto.NewServerClient(conn)
	}
	return conns
}

func (frontend *Frontend) Result(ctx context.Context, in *proto.Void) (*proto.BidResult, error) {
	counter := 0
	acks := make([]*proto.BidResult, 3)
	for i := 0; i < 3; i++ {
		ack, err := frontend.serverConnection[i].GetHighestBid(ctx, in)

		if err == nil {
			log.Printf("Frontend requested result in node %d\n", i)
			acks[counter] = ack
			counter++
		}
	}
	if len(acks) == 0 {
		return nil, status.Errorf(1, "error")
	}
	return acks[0], nil
}

func (frontend *Frontend) Bid(ctx context.Context, in *proto.BidRequest) (*proto.Ack, error) {
	counter := 0
	acks := make([]*proto.Ack, 3)
	log.Printf("Frontend received bid of %d by %s", in.Amount, in.Name)
	for i := 0; i < 3; i++ {
		conn := frontend.serverConnection[i]
		ack, err := conn.UpdateHighestBid(ctx, in)

		if err == nil {
			log.Printf("sends updatehighest bid to node %d\n", i)
			acks[counter] = ack
			counter++
		}
	}
	if len(acks) == 0 {
		return nil, status.Errorf(1, "error")
	}
	return acks[0], nil
}

func (frontend *Frontend) StartAuction(ctx context.Context, in *proto.Void) (*proto.Ack, error) {
	counter := 0
	acks := make([]*proto.Ack, 3)
	for i := 0; i < 3; i++ {

		ack, err := frontend.serverConnection[i].StartAuction(ctx, in)
		if err == nil {
			log.Printf("Frontend started auction in node %d\n", i)
			acks[counter] = ack
			counter++
		}
	}
	if len(acks) == 0 {
		return nil, status.Errorf(1, "error")
	}
	return acks[0], nil
}
