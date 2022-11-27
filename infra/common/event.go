package common

import (
	"context"
	"fmt"
	"sync"

	"alsritter.icu/stgo/domain/event"
	commonEvent "alsritter.icu/stgo/infra/common/event"
)

type EventWithCtx struct {
	Ctx context.Context
	E   event.Event
	WG  *sync.WaitGroup
}

// EventPublish 事件发布
func EventPublish(ctx context.Context, e event.Event) {
	if e.IsEmpty() {
		return
	}

	msg := commonEvent.Message{
		Topic: e.GetTopic(),
		Body:  e,
	}

	var (
		msgId string
		err   error
	)

	if !e.IsDelayEvent() {
		msgId, err = EventServer.Publish(ctx, &msg)
	} else {
		delayAt := e.GetDelayAt()
		msgId, err = EventServer.PublishDelay(ctx, &commonEvent.DelayMessage{
			Message: msg,
			// fmt.Printf("时间戳（秒）：%v\n", time.Now().Unix())
			// fmt.Printf("时间戳（纳秒）：%v\n", time.Now().UnixNano())
			// fmt.Printf("时间戳（毫秒）：%v\n", time.Now().UnixNano()/1e6)
			// fmt.Printf("时间戳（纳秒转换为秒）：%v\n", time.Now().UnixNano()/1e9)
			DelayTimeLevel: delayAt,
		})
	}

	if err != nil {
		fmt.Printf("事件队列发布事件失败, err:%v, topic:%s, body:%v \n", err, e.GetTopic(), e)
		return
	}
	fmt.Printf("事件队列发布事件成功, msgId:%s, topic:%s, body:%v \n", msgId, e.GetTopic(), e)
}
