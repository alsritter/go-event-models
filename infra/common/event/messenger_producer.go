package event

import (
	"context"
	"encoding/json"
	"time"
)

type MessengerProducer struct {
	topic  string
	client *MessengerClient
}

func NewMessengerProducer(m *Messenger, topic string) (InnerProducer, error) {
	return &MessengerProducer{
		topic:  topic,
		client: m.client,
	}, nil
}

func (p *MessengerProducer) publish(ctx context.Context, msg *Message) (string, error) {
	msgBody, err := json.Marshal(msg.Body)
	if err != nil {
		return "", err
	}

	req := &PublishRequest{
		Topic: p.topic,
		Body:  string(msgBody),
	}

	retryCount := 4
	var msgId string
	for i := 1; i <= retryCount; i++ {
		msgId, err = p.client.Publish(ctx, req)
		if err == nil {
			break
		}

		//防止 broker 连续限流，一次限流3秒
		time.Sleep(time.Second * time.Duration(i*3))

		if i == retryCount-1 {
			return "", err
		}
	}

	return msgId, nil
}

func (p *MessengerProducer) publishDelay(ctx context.Context, msg *DelayMessage) (string, error) {
	msgBody, err := json.Marshal(msg.Body)
	if err != nil {
		return "", err
	}

	req := &PublishDelayRequest{
		PublishRequest: PublishRequest{
			Topic: p.topic,
			Body:  string(msgBody),
		},
		DelayTimeLevel: msg.DelayTimeLevel,
	}

	retryCount := 4
	var msgId string
	for i := 1; i <= retryCount; i++ {
		msgId, err = p.client.PublishDelay(ctx, req)
		if err == nil {
			break
		}

		//防止 broker 连续限流，一次限流3秒
		time.Sleep(time.Second * time.Duration(i*3))

		if i == retryCount-1 {
			return "", err
		}
	}

	return msgId, nil
}
