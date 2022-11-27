package event

import (
	"context"
	"fmt"
	"sync"
)

type Messenger struct {
	client      *MessengerClient
	consumerMap sync.Map
	producerMap sync.Map
}

func NewMessenger(opts ...Option) (EventServerIface, error) {
	client := NewMessengerClient(opts...)
	if err := client.Connect(); err != nil {
		return nil, err
	}

	return &Messenger{
		client: client,
	}, nil
}

// 发布消息
func (m *Messenger) Publish(ctx context.Context, msg *Message) (string, error) {
	var (
		producer InnerProducer
		err      error
	)

	if p, ok := m.producerMap.Load(msg.Topic); ok {
		producer = p.(InnerProducer)
	} else {
		producer, err = NewMessengerProducer(m, msg.Topic)
		if err != nil {
			return "", fmt.Errorf("NewMessengerProducer err: %v", err)
		}
		m.producerMap.Store(msg.Topic, producer)
	}

	return producer.publish(ctx, msg)
}

// 发布延迟消息
func (m *Messenger) PublishDelay(ctx context.Context, msg *DelayMessage) (string, error) {
	var (
		producer InnerProducer
		err      error
	)

	if p, ok := m.producerMap.Load(msg.Topic); ok {
		producer = p.(InnerProducer)
	} else {
		producer, err = NewMessengerProducer(m, msg.Topic)
		if err != nil {
			return "", fmt.Errorf("NewMessengerProducer err: %v", err)
		}
		m.producerMap.Store(msg.Topic, producer)
	}
	return producer.publishDelay(ctx, msg)
}

// 订阅消息
func (m *Messenger) Subscribe(topic, consumer string, handler ReceiveCallBack) error {
	c, err := m.GetConsumer(topic, consumer)
	if err != nil {
		return err
	}
	if c == nil {
		mc, err := NewMessengerConsumer(m, topic, consumer, handler)
		if err != nil {
			return fmt.Errorf("NewMessengerConsumer err: %v", err)
		}
		m.consumerMap.Store(fmt.Sprintf("%s_%s", topic, consumer), mc)
		c = mc
	}

	if err := c.Start(); err != nil {
		return fmt.Errorf("consumer Start err: %v", err)
	}
	return nil
}

// 取消订阅
func (m *Messenger) Unsubscribe(topic, consumer string) error {
	c, err := m.GetConsumer(topic, consumer)
	if err != nil {
		return err
	}
	if c == nil {
		return nil
	}

	mc := c.(*MessengerConsumer)
	if err := mc.Stop(); err != nil {
		fmt.Printf("Consumer Stop err: %v", err)
		return err
	}

	m.consumerMap.Delete(fmt.Sprintf("%s_%s", mc.topic, mc.consumer))
	return nil
}

// 获取消费者
func (m *Messenger) GetConsumer(topic, consumer string) (InnerConsumer, error) {
	c, ok := m.consumerMap.Load(fmt.Sprintf("%s_%s", topic, consumer))
	if ok {
		return c.(InnerConsumer), nil
	}
	return nil, nil
}
