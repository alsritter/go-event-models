package event

import (
	"context"
)

type InnerProducer interface {
	publish(ctx context.Context, msg *Message) (string, error)
	publishDelay(ctx context.Context, msg *DelayMessage) (string, error)
}

type InnerConsumer interface {
	Start() error
	Stop() error
}

type EventServerIface interface {
	// 发布消息
	Publish(ctx context.Context, msg *Message) (string, error)
	// 发布延迟消息
	PublishDelay(ctx context.Context, msg *DelayMessage) (string, error)
	// 订阅消息
	Subscribe(topic, consumer string, handler ReceiveCallBack) error
	// 取消订阅
	Unsubscribe(topic, consumerName string) error
	// 获取消费者
	GetConsumer(topic, consumerName string) (InnerConsumer, error)
}
