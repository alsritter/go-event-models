package event

const OrderDelayAutoCancelEventTopic = "orderDelayAutoCancel"

type OrderDelayAutoCancelEvent struct {
	DelayBase
	OrderId string
}

func NewOrderDelayAutoCancelEvent(orderId string, delay int64) *OrderDelayAutoCancelEvent {
	return &OrderDelayAutoCancelEvent{OrderId: orderId, DelayBase: DelayBase{delay: delay}}
}

func (a *OrderDelayAutoCancelEvent) GetTopic() string {
	return OrderDelayAutoCancelEventTopic
}
