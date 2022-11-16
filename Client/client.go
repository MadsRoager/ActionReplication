package main

import (
	"AuctionReplication/proto"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MadsRoager/AuctionReplication/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	name string
	id   int
}

var name = flag.String("clientName", "DefaultName", "client name")
var id = flag.Int("id", 0, "id")
var frontendPort = flag.Int("frontendPort", 8000, "frontend port")

func main() {

	client := &Client{
		name: 	*name,
		id: 	*id,
	}

	go startClient(client)

	for {
		time.Sleep(100 * time.Second)
	}
}

func startClient(client *Client) {

	serverConnection := getServerConnection(client)

	sendMessage(serverConnection)	

}

func getServerConnection(client *Client) proto.MessagingServiceClient {
	conn, err := grpc.Dial(":"+strconv.Itoa(*frontendPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Could not dial server")
	}
	return proto.NewMessagingServiceClient(conn)
}

func sendMessage(serverConnection proto.FrontendClient) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.HasPrefix(input, "bid") {
			split := strings.Split(input, " ")
			ans, err := serverConnection.Bid(&proto.BidRequest{
				Amount: split[1],
				Name: name,
				ProcessID: id,
			})

		} 
		if input == "result" {
			ans, err := serverConnection.Result(&proto.Void{})
			log.Println(ans.Amount)
		} 
	}
}