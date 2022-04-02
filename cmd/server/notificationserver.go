package main

import (
	"io"
	"log"

	"github.com/nlm/grpceventbus/eventpb"
	"github.com/nlm/grpceventbus/pkg/pubsub"
)

type NotificationServer struct {
	eventpb.UnimplementedApiServer
	pubsub *pubsub.PubSub[*eventpb.Event]
	logger *log.Logger
}

func NewNotificationServer() *NotificationServer {
	return &NotificationServer{
		pubsub: pubsub.NewPubSub[*eventpb.Event](),
		logger: log.New(
			log.Default().Writer(),
			"[NotificationServer]",
			log.Default().Flags(),
		),
	}
}

func (ns *NotificationServer) Close() {
	ns.pubsub.Close()
}

func (ns *NotificationServer) Publish(stream eventpb.Api_PublishServer) error {
	// Handle Publishing
	for {
		pr, err := stream.Recv()
		switch {
		case err == io.EOF:
			return stream.SendAndClose(&eventpb.Empty{})
		case err != nil:
			return err
		default:
			ns.pubsub.Publish(pr.Topic, pr.Event)
		}
	}
}

func (ns *NotificationServer) Subscribe(sr *eventpb.SubscribeRequest, stream eventpb.Api_SubscribeServer) error {
	// Connect to NATS
	sub := ns.pubsub.Subscribe(sr.Topic)
	ns.logger.Println("subscribed")
	ctx := stream.Context()
	for {
		select {
		case event := <-sub.C():
			stream.Send(event)
		case <-ctx.Done():
			ns.pubsub.Unsubscribe(sub)
			ns.logger.Println("unsubscribed")
			return ctx.Err()
		}
	}
}
