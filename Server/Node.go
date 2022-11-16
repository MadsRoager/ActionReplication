package Server

import (
	"flag"
	"github.com/MadsRoager/AuctionReplication/proto"
	"sync"
)

var highestBid int32 = 0
var highestBidder string = ""
var highestBidderID int32 = 0

var mutex = sync.Mutex{}
var port = flag.Int("port", 8080, "server port number")

func main() {
	flag.Parse()
}

func auctionIsOver() {

}

func evaluateNewBid(bid *proto.BidRequest) proto.Ack {
	if isWinningBet(bid) {
		updateHighestBid(bid)
		return proto.Ack{}
	}
	return proto.Ack{}
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
