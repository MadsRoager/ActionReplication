package main

import (
	"flag"
	"fmt"
	"github.com/MadsRoager/AuctionReplication/proto"
	"sync"
	"time"
)

var highestBid int32 = 0
var highestBidder string = "No bidder yet"
var highestBidderID int32 = 0

var mutex = sync.Mutex{}
var port = flag.Int("port", 8080, "server port number")

func main() {
	flag.Parse()
	defer runAuction()
}

var counter = 0

func runAuction() {
	for {
		time.Sleep(time.Second)
		fmt.Println(counter)
		counter++
	}

}

func evaluateNewBid(bid *proto.BidRequest) proto.Ack {
	if isWinningBet(bid) {
		updateHighestBid(bid)
		return success()
	}
	return fail()
}

func isWinningBet(bid *proto.BidRequest) bool {
	if bid.Amount > highestBid {
		return true
	}
	if bid.Amount == highestBid && bid.ProcessID < highestBidderID {
		return true
	}
	return false
}

func updateHighestBid(bid *proto.BidRequest) {
	mutex.Lock()
	highestBid = bid.Amount
	highestBidder = bid.Name
	mutex.Unlock()
}

func result() *proto.BidResult {
	result := &proto.BidResult{
		Amount:        highestBid,
		Name:          highestBidder,
		AuctionStatus: "Ongoing",
	}
	return result
}

func success() proto.Ack {
	return proto.Ack{
		Ack: "Success",
	}
}

func fail() proto.Ack {
	return proto.Ack{
		Ack: "Fail",
	}
}
