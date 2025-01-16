package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	Environment string
	ProjectID   string
	Client      *pubsub.Client
}

func getTopicName(message proto.Message, env string) string {
	fullMessageName := string(proto.MessageName(message))
	return fmt.Sprintf("%s__event__%s", env, fullMessageName)
}

func (p Publisher) Initialize(ctx context.Context) {
	client, err := pubsub.NewClient(ctx, p.ProjectID)
	if err != nil {
		panic(err)
	}
	p.Client = client
}

// Publish publishes a message to the topic
func (p Publisher) Publish(ctx context.Context, message proto.Message) error {
	topicName := getTopicName(message, p.Environment)
	topic := p.Client.Topic(topicName)
	defer topic.Stop()
	msg, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	_, err = topic.Publish(ctx, &pubsub.Message{Data: msg}).Get(ctx)
	if err != nil {
		return err
	}
	return nil
}
