package event

import (
	"context"
)

type MessengerConsumer struct {
	topic       string
	consumer    string
	isBroadcast bool
	handler     ReceiveCallBack
	client      *MessengerClient
}

func NewMessengerConsumer(m *Messenger,
	topic, consumer string,
	handler ReceiveCallBack) (InnerConsumer, error) {
	return &MessengerConsumer{topic: topic, handler: handler, consumer: consumer, client: m.client}, nil
}

func (c *MessengerConsumer) Start() error {
	req := &SubscribeRequest{
		Topic:       c.topic,
		IsBroadcast: c.isBroadcast,
	}

	if err := c.client.Subscribe(context.Background(), req, c.handler); err != nil {
		return err
	}
	return nil
}

func (c *MessengerConsumer) Stop() error {
	if err := c.client.Unsubscribe(context.Background(), c.topic); err != nil {
		return err
	}
	return nil
}
