package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/nlm/grpceventbus/eventpb"

	"google.golang.org/grpc"
)

var (
	listenPort = flag.Int("port", 8080, "listen port")
	useNats    = flag.Bool("nats", false, "use nats backend")
)

func main() {
	// Parse Flags
	flag.Parse()
	// Listen
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *listenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Build gRPC Server
	grpcServer := grpc.NewServer()
	// Register Services
	if *useNats {
		eventpb.RegisterApiServer(grpcServer, NewNatsNotificationServer())
	} else {
		eventpb.RegisterApiServer(grpcServer, NewNotificationServer())
	}
	// Server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
