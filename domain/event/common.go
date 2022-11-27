package event

type Event interface {
	// IsEmpty 判断是否是空事件
	IsEmpty() bool
	// GetTopic 获取事件发送对应的 Topic
	GetTopic() string
	// IsDelayEvent 是否是延迟事件
	IsDelayEvent() bool
	// GetDelayAt 获取延迟任务时间
	GetDelayAt() int64
}

type Base struct{}

func (b *Base) IsEmpty() bool {
	return false
}

func (b *Base) IsDelayEvent() bool {
	return false
}

func (b *Base) GetDelayAt() int64 {
	return 0
}

type DelayBase struct {
	Base
	delay int64 // 延迟任务执行时间点
}

func (b *DelayBase) IsDelayEvent() bool {
	return true
}

func (b *DelayBase) GetDelayAt() int64 {
	return b.delay
}

var EmptyEvent Event = new(emptyEvent)

type emptyEvent struct{}

func (e *emptyEvent) IsEmpty() bool      { return true }
func (e *emptyEvent) GetTopic() string   { return "" }
func (e *emptyEvent) IsDelayEvent() bool { return false }
func (e *emptyEvent) GetDelayAt() int64  { return 0 }
func (e *emptyEvent) SetOrgCode(string)  {}
