package event

import (
	"encoding/json"
	"fmt"
)

const (
	ORG_CODE_KEY = "orgcode"
	TRACE_ID     = "TRACE_ID"
)

const (
	//延迟消息级别
	//1=>1s 2=>5s 3=>10s 4=>30s 5=>1m 6=>2m 7=>3m 8=>4m 9=>5m 10=>6m 11=>7m 12=>8m 13=>9m 14=>10m 15=>20m 16=>30m 17=>1h 18=>2h
	DelayTimeLevel1Second   = 1
	DelayTimeLevel5Seconds  = 2
	DelayTimeLevel10Seconds = 3
	DelayTimeLevel30Seconds = 4
	DelayTimeLevel1Minute   = 5
	DelayTimeLevel2Minutes  = 6
	DelayTimeLevel3Minutes  = 7
	DelayTimeLevel4Minutes  = 8
	DelayTimeLevel5Minutes  = 9
	DelayTimeLevel6Minutes  = 10
	DelayTimeLevel7Minutes  = 11
	DelayTimeLevel8Minutes  = 12
	DelayTimeLevel9Minutes  = 13
	DelayTimeLevel10Minutes = 14
	DelayTimeLevel20Minutes = 15
	DelayTimeLevel30Minutes = 16
	DelayTimeLevel1Hour     = 17
	DelayTimeLevel2Hours    = 18
)

type Message struct {
	// topic
	Topic string `json:"topic"`
	// 消息体
	Body interface{} `json:"body"`
}

type DelayMessage struct {
	Message
	//延迟消息级别
	// WithDelayTimeLevel set message delay time to consume.
	// reference delay level definition:
	// 1=>1s 2=>5s 3=>10s 4=>30s 5=>1m 6=>2m 7=>3m 8=>4m 9=>5m 10=>6m 11=>7m 12=>8m 13=>9m 14=>10m 15=>20m 16=>30m 17=>1h 18=>2h
	// delay level starts from 1. for example, if we set param level=1, then the delay time is 1s.
	DelayTimeLevel int64 `json:"delay_time_level" `
}

func (msg Message) Unmarshal(body interface{}) error {
	str := fmt.Sprintf("%s", msg.Body)
	return json.Unmarshal([]byte(str), body)
}
