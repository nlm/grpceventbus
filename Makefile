.PHONY: proto

all: proto server publisher subscriber

proto: eventpb/event.pb.go eventpb/event_grpc.pb.go

eventpb/event_grpc.pb.go: eventpb/event.proto
	protoc --go-grpc_out=eventpb eventpb/event.proto

eventpb/event.pb.go: eventpb/event.proto
	protoc --go_out=eventpb eventpb/event.proto

server: cmd/server/*.go eventpb/*.go pkg/pubsub/*.go
	go build -o $@ ./cmd/server

publisher: cmd/publisher/*.go eventpb/*.go
	go build -o $@ ./cmd/publisher

subscriber: cmd/subscriber/*.go eventpb/*.go
	go build -o $@ ./cmd/subscriber

.PHONY: clean

clean:
	rm -f server publisher subscriber
