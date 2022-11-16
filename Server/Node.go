package Server

import (
	"flag"
	"github.com/MadsRoager/AuctionReplication/proto"
	"sync"
)

var highestBid = 0
var mutex = sync.Mutex{}
var port = flag.Int("port", 8080, "server port number")

func main() {
	flag.Parse()
}

func auctionIsOver() {

}

func updateHighestBid(bid proto.BidRequest) proto.Ack {
	if highestBid < bid.Amount {

	}
}

func result() proto.BidRequest {

}
