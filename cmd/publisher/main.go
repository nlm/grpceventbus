package main

import (
	"context"
	"flag"
	"fmt"
	"grpctest/eventpb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	connectPort = flag.Int("port", 8080, "connection port")
)

func main() {
	var dialOptions []grpc.DialOption
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", *connectPort), dialOptions...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := eventpb.NewApiClient(conn)
	sc, err := client.Subscribe(context.Background(), &eventpb.SubscribeRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		event, err := sc.Recv()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("event received:", event.Event)
	}
}
