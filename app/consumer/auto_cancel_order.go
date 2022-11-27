package consumer

import (
	"context"
	"fmt"
	"time"

	"alsritter.icu/stgo/domain/event"
	"alsritter.icu/stgo/infra/common"
	commonEvent "alsritter.icu/stgo/infra/common/event"
)

// AutoCancelOrderConsumer 自动取消订单消费者
type AutoCancelOrderConsumer struct{ common.BaseConsumer }

func NewAutoCancelOrderConsumer() common.IConsumer {
	return &AutoCancelOrderConsumer{common.BaseConsumer{Name: common.AutoCancelOrderConsumer}}
}

// RegisterNeededEvent 需要在这里调用 RegisterEvent 注册所有需要的事件，然后在构造方法中调用
func (c *AutoCancelOrderConsumer) RegisterNeededEvent() {
	// 注意，这里注册了两个消费处理器，它们绑定了不同的事件
	c.RegisterEvent(event.OrderSubmitEventTopic, c.handleOrderSubmitted)
	c.RegisterEvent(event.OrderDelayAutoCancelEventTopic, c.handleOrderAutoCancel)
}

// 处理收到订单消息，再发送延迟取消订单消息
func (c *AutoCancelOrderConsumer) handleOrderSubmitted(ctx context.Context, msg *commonEvent.Message) error {
	// 收到消息之后，先把消息序列化到 Event 里面
	e := event.OrderSubmittedEvent{}
	if err := msg.Unmarshal(&e); err != nil {
		return err
	}

	// 这里做一些处理
	fmt.Printf("取得提交的订单 %v \n", e)

	// 再发送延迟取消订单消息
	common.EventPublish(ctx, event.NewOrderDelayAutoCancelEvent(e.OrderId, commonEvent.DelayTimeLevel5Seconds))
	return nil
}

func (c *AutoCancelOrderConsumer) handleOrderAutoCancel(ctx context.Context, msg *commonEvent.Message) error {
	e := event.OrderDelayAutoCancelEvent{}
	if err := msg.Unmarshal(&e); err != nil {
		return err
	}

	fmt.Printf("收到需要延时取消的订单 %v \n", e)
	return nil
}

func TimeSetDefaultLoc(timeSuk *time.Time) time.Time {
	var loc, _ = time.LoadLocation("Asia/Shanghai")
	return timeSuk.In(loc)
}
