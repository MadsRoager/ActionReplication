# AuctionReplication
## Commands for renewing proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/proto.proto

go mod tidy