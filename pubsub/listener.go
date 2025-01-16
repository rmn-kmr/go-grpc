package pubsub

import (
	"context"
	"fmt"

	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/ptypes/empty"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rmnkmr/go-common/errors"
	"github.com/rmnkmr/lsp/log"
	"github.com/rmnkmr/protonium"
	"google.golang.org/protobuf/proto"
)

type Listener struct {
	Environment    string
	Client         *pubsub.Client
	SubscriptionID string
	Message        proto.Message
}

func (s Subscription[T]) Decode(data []byte) (T, error) {
	var x struct {
		Name    string `json:"name"`
		Payload T      `json:"payload"`
	}
	err := json.Unmarshal(data, &x)
	return x.Payload, err
}

type Subscription[T any] struct {
	Sub     *pubsub.Subscription
	Handler func(context.Context, T) (*empty.Empty, error)
	Decoder func([]byte) (T, error)
}

func (s Subscription[T]) Initialize(*nr.Application) {}
func (s Subscription[T]) Start(ctx context.Context) error {
	return s.Sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		// TODO: fix the logger to support Derive and WithFields
		//ctx = log.Derive(ctx, log.WithFields("event", s.Sub.String()))
		var msg T
		if s.Decoder != nil {
			// decode the message
			message, err := s.Decoder(m.Data)
			if err != nil {
				log.Error(ctx, err, "malformed-event", "evt", string(m.Data))
				m.Ack()
				return
			}
			msg = message
		} else {
			// decode the message
			message, err := s.Decode(m.Data)
			if err != nil {
				log.Error(ctx, err, "malformed-event", "evt", string(m.Data))
				m.Ack()
				return
			}
			msg = message
		}
		// call the handler
		if _, err := s.Handler(ctx, msg); err != nil {
			log.Error(ctx, err, "failed to process event")
			if errors.IsRetryable(err) {
				m.Nack()
				return
			}
		}
		m.Ack()
	})
}

func getSubscriptionName(message proto.Message, env string) string {
	fullMessageName := string(proto.MessageName(message))
	return fmt.Sprintf("%s__%s__listener__%s", env, "lsp", fullMessageName)
}

func ListenerOption[T any](l *Listener, handler func(context.Context, T) (*empty.Empty, error), decoder func([]byte) (T, error)) protonium.Option {
	if l.Message == nil && l.SubscriptionID == "" {
		return nil
	}
	subscriptionID := l.SubscriptionID
	var name string
	if subscriptionID == "" {
		subscriptionID = getSubscriptionName(l.Message, l.Environment)
		name = string(l.Message.ProtoReflect().Descriptor().Name())
	} else {
		name = subscriptionID
	}
	s := Subscription[T]{
		Sub:     l.Client.Subscription(subscriptionID),
		Handler: handler,
		Decoder: decoder,
	}
	return protonium.Component(name, s)
}
