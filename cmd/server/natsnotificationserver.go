package main

import (
	"io"
	"log"

	nats "github.com/nats-io/nats.go"
	"github.com/nlm/grpceventbus/eventpb"
	"google.golang.org/protobuf/proto"
)

type NatsNotificationServer struct {
	eventpb.UnimplementedApiServer
	logger *log.Logger
}

func NewNatsNotificationServer() *NatsNotificationServer {
	return &NatsNotificationServer{
		logger: log.New(
			log.Default().Writer(),
			"[NatsNotificationServer]",
			log.Default().Flags(),
		),
	}
}

func (ns *NatsNotificationServer) Close() {
}

func (ns *NatsNotificationServer) Publish(stream eventpb.Api_PublishServer) error {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()
	// Handle Publishing
	for {
		pr, err := stream.Recv()
		switch {
		case err == io.EOF:
			return stream.SendAndClose(&eventpb.Empty{})
		case err != nil:
			return err
		default:
			payload, err := proto.Marshal(pr.Event)
			if err != nil {
				return err
			}
			nc.Publish(pr.Topic, payload)
		}
	}
}

func (ns *NatsNotificationServer) Subscribe(sr *eventpb.SubscribeRequest, stream eventpb.Api_SubscribeServer) error {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()
	ch := make(chan *nats.Msg, 64)
	defer close(ch)
	sub, err := nc.ChanSubscribe(sr.Topic, ch)
	if err != nil {
		return err
	}
	ns.logger.Println("subscribed")

	ctx := stream.Context()
	event := &eventpb.Event{}
	for {
		select {
		case msg := <-ch:
			if err := proto.Unmarshal(msg.Data, event); err != nil {
				return err
			}
			stream.Send(event)
		case <-ctx.Done():
			//ns.pubsub.Unsubscribe(sub)
			if err := sub.Unsubscribe(); err != nil {
				return err
			}
			ns.logger.Println("unsubscribed")
			return ctx.Err()
		}
	}
}
