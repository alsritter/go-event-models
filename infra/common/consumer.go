package common

import (
	"context"
	"fmt"

	"alsritter.icu/stgo/infra/common/crypt"
	"alsritter.icu/stgo/infra/common/event"
)

type IConsumer interface {
	// RegisterEvent 注册 Consumer 需要观察的 EventHandler
	RegisterEvent(topic string, handler EventHandlerFunc)
	// RegisterConsumer 注册当前消费者（订阅 Consumer 需要的 Event）
	RegisterConsumer(eventServer event.EventServerIface) error
	// RegisterNeededEvent 子类需要在这里调用 RegisterEvent 注册所有需要的事件,然后在构造方法中调用
	RegisterNeededEvent()
	// GetHandleMap 获取注册表
	GetHandleMap() map[string]EventHandlerFunc
}

type BaseConsumer struct {
	Name       string                      // 消费者名字，唯一
	handlerMap map[string]EventHandlerFunc // 事件处理器集合 topic => handler
}

type EventHandlerFunc func(ctx context.Context, message *event.Message) error

// RegisterEvent 注册 Consumer 需要观察的 EventHandler
func (b *BaseConsumer) RegisterEvent(topic string, handler EventHandlerFunc) {
	if nil == b.handlerMap {
		b.handlerMap = make(map[string]EventHandlerFunc)
	}
	b.handlerMap[topic] = b.wrapWithRetry(handler)
}

// RegisterConsumer 注册当前消费者（订阅 Consumer 需要的 Event）
func (b *BaseConsumer) RegisterConsumer(eventServer event.EventServerIface) error {
	for topic := range b.handlerMap {
		internalConsumerName := b.getInternalConsumerName(topic)
		if err := eventServer.Subscribe(topic, internalConsumerName, b.handle); err != nil {
			return err
		}
		fmt.Printf("消费者 %s 注册 topic %s 内部消费者名称 %s \n", b.Name, topic, internalConsumerName)
	}
	return nil
}

func (b BaseConsumer) getInternalConsumerName(topic string) string {
	// 这里用 MD5 是因为要防止 Consume Name 超长（最长32位）
	return crypt.MD5([]byte(fmt.Sprintf("%s-%s", b.Name, topic)))
}

// WrapWithRetry 重试目前只处理 IO 错误（DB、Redis、Conn），其余错误不抛回给消息队列也不会记录日志
func (b *BaseConsumer) wrapWithRetry(handler func(ctx context.Context, msg *event.Message) error) func(ctx context.Context, msg *event.Message) error {
	return func(ctx context.Context, msg *event.Message) error {
		retry, err := event.RecoverToRetry(ctx, msg, handler)
		if err != nil {
			fmt.Printf("[%s] 消费者 %s 处理事件, body [%s] retry:[%t] err:[%s] \n", b.getInternalConsumerName(msg.Topic), b.Name, msg.Body, retry, err)
		}
		if retry {
			return err
		}
		// 不需要重试只记录日志
		return nil
	}
}

// 事件处理分发
func (b *BaseConsumer) handle(ctx context.Context, msg *event.Message) error {
	// 这里一定要过 Handler，然后再分发事件
	// 不能直接注册对应的 EventHandlerFunc，因为会违反《订阅关系一致》
	handler, ok := b.handlerMap[msg.Topic]

	if !ok {
		return fmt.Errorf("[%s] 消费者处理事件时发生错误：事件 Topic %s 未注册对应 Handler", b.Name, msg.Topic)
	}

	fmt.Printf("[%s] 消费者处理事件 topic [%s], body [%s] \n", b.Name, msg.Topic, msg.Body)
	return handler(ctx, msg)
}

func (b *BaseConsumer) GetHandleMap() map[string]EventHandlerFunc {
	return b.handlerMap
}

func (b *BaseConsumer) PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
