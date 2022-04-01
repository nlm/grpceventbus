package main

import (
	"flag"
	"fmt"
	"github.com/nlm/grpceventbus/eventpb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type NotificationServer struct {
	eventpb.UnimplementedApiServer
}

var (
	listenPort     = flag.Int("port", 8080, "listen port")
	tickerInterval = flag.Duration("interval", 3*time.Second, "emission interval")
)

func (ns *NotificationServer) Subscribe(in *eventpb.SubscribeRequest, srv eventpb.Api_SubscribeServer) error {
	log.Println("subscribed")
	ticker := time.NewTicker(*tickerInterval)
	defer ticker.Stop()
	ctx := srv.Context()

	for {
		select {
		case <-ticker.C:
			log.Println("emitted event")
			err := srv.Send(&eventpb.Event{
				Event: &eventpb.Event_FooEvent{
					FooEvent: &eventpb.FooEvent{
						Foo: "foo",
					},
				},
			})
			if err != nil {
				log.Println("error:", err)
			}
		case <-ctx.Done():
			log.Println("unsubscribed")
			return ctx.Err()
		}
	}
}

func main() {
	// Parse Flags
	flag.Parse()
	// Build NotificationServer
	notifServer := &NotificationServer{}
	// Listen
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *listenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Build gRPC Server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	// Register Services
	eventpb.RegisterApiServer(grpcServer, notifServer)
	// Seerver
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
