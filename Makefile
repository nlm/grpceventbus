.PHONY: proto

all: proto server publisher subscriber

proto: eventpb/event.pb.go eventpb/event_grpc.pb.go

eventpb/event_grpc.pb.go: eventpb/event.proto
	protoc --go-grpc_out=eventpb eventpb/event.proto

eventpb/event.pb.go: eventpb/event.proto
	protoc --go_out=eventpb eventpb/event.proto

server: cmd/server/*.go
	go build -o $@ $<

publisher: cmd/publisher/*.go
	go build -o $@ $<

subscriber: cmd/subscriber/*.go
	go build -o $@ $<

.PHONY: clean

clean:
	rm -f server publisher subscriber
