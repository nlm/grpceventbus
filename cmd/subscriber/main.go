package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/nlm/grpceventbus/eventpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	connectPort = flag.Int("port", 8080, "connection port")
	topicName   = flag.String("topic", "default", "topic name")
)

func main() {
	// Dial Server
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", *connectPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := eventpb.NewApiClient(conn)
	// Event subscriber
	sc, err := client.Subscribe(context.Background(), &eventpb.SubscribeRequest{
		Topic: *topicName,
	})
	if err != nil {
		log.Fatal(err)
	}
	// Receive events
	for {
		event, err := sc.Recv()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("event received:", event.Event)
	}
}
