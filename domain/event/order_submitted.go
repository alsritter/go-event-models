package event

const OrderSubmitEventTopic = "orderSubmitted"

type OrderSubmittedEvent struct {
	Base
	OrderId string
	Name    string
	Data    string
}

func (a *OrderSubmittedEvent) GetTopic() string {
	return OrderSubmitEventTopic
}
