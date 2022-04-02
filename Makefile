.PHONY: proto

PBGOFILES=eventpb/event.pb.go eventpb/event_grpc.pb.go

all: proto server publisher subscriber

proto: $(PBGOFILES)

eventpb/event_grpc.pb.go: eventpb/event.proto
	protoc --go-grpc_out=eventpb eventpb/event.proto

eventpb/event.pb.go: eventpb/event.proto
	protoc --go_out=eventpb eventpb/event.proto

server: cmd/server/*.go pkg/pubsub/*.go $(PBGOFILES)
	go build -o $@ ./cmd/server

publisher: cmd/publisher/*.go $(PBGOFILES)
	go build -o $@ ./cmd/publisher

subscriber: cmd/subscriber/*.go $(PBGOFILES)
	go build -o $@ ./cmd/subscriber

.PHONY: clean distclean

clean:
	rm -f server publisher subscriber

distclean: clean
	rm -f $(PBGOFILES)
