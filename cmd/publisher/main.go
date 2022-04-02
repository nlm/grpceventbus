package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/nlm/grpceventbus/eventpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	connectPort  = flag.Int("port", 8080, "connection port")
	sendInterval = flag.Duration("interval", 3*time.Second, "event interval")
	topicName    = flag.String("topic", "default", "topic name")
)

func main() {
	flag.Parse()
	// dial gRPC Server
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", *connectPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := eventpb.NewApiClient(conn)
	// Event streaming
	stream, err := client.Publish(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(*sendInterval)
	defer ticker.Stop()
	for t := range ticker.C {
		// Demo Events
		events := []*eventpb.Event{
			{
				Kind: &eventpb.Event_FooEvent_{
					FooEvent: &eventpb.Event_FooEvent{
						Foo: fmt.Sprint(t.Clock()),
					},
				},
			},
			{
				Kind: &eventpb.Event_BarEvent_{
					BarEvent: &eventpb.Event_BarEvent{
						Bar: fmt.Sprint(t.Clock()),
					},
				},
			},
		}
		// Send Events
		for _, event := range events {
			err := stream.Send(&eventpb.PublishRequest{
				Topic: *topicName,
				Event: event,
			})
			if err != nil {
				log.Fatal(err)
			}
			log.Println("event sent:", event)
		}
	}
}
