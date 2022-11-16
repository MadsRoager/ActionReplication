package Client

import (
	"flag"
)

type Client struct {
	name string
	id   int
}

var name = flag.String("clientName", "DefaultName", "client name")
var id = flag.Int()

func main() {

	client := &Client{
		name: *name,
	}
}