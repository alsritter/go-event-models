package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/sirupsen/logrus"
)

type ReceiveCallBack func(ctx context.Context, msg *Message) error

type SubscribeRequest struct {
	Topic       string
	IsBroadcast bool
}

type MessengerClient struct {
	producer     rocketmq.Producer
	consumer     rocketmq.PushConsumer
	consumerOnce sync.Once

	opts *mqOptions
}

// NewMessengerClient 创建一个MQ客户端
func NewMessengerClient(opts ...Option) *MessengerClient {
	rlog.SetLogLevel("warn") // 避免产生太多日志
	var m = new(MessengerClient)
	defaultOpts := defaultOptions()
	for _, apply := range opts {
		apply(&defaultOpts)
	}
	m.opts = &defaultOpts
	return m
}

func (m *MessengerClient) createConsumer() error {
	//消息主动推送给消费者
	c, err := rocketmq.NewPushConsumer(
		consumer.WithInstance(m.opts.consumerInstance), //必须设置，否则广播模式会重复消费
		consumer.WithGroupName(m.opts.groupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{m.opts.addr})),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset), //选择消费时间（首次/当前/根据时间）
		consumer.WithConsumerModel(consumer.Clustering),                //消费模式(集群消费:消费完,同组的其他人不能再读取/广播消费：所有人都能读)
	)

	if err != nil {
		logrus.Error("createConsumer NewPushConsumer failed : ", err.Error())
		return err
	}
	m.consumer = c
	return nil
}

// 创建生产者
func (m *MessengerClient) createProducer() error {
	addr, err := primitive.NewNamesrvAddr(m.opts.addr)
	if err != nil {
		logrus.Error("createProducer NewNamesrvAddr failed : ", err.Error())
		return err
	}
	p, err := rocketmq.NewProducer(
		producer.WithGroupName(m.opts.groupName),
		producer.WithNameServer(addr),
		producer.WithRetry(m.opts.retry),
	)
	if err != nil {
		logrus.Error("createProducer NewProducer failed : ", err.Error())
		return err
	}
	if err = p.Start(); err != nil {
		logrus.Error("createProducer Start failed : ", err.Error())
		return err
	}
	m.producer = p
	return nil
}

// 创建生产者和消费者
func (m *MessengerClient) Connect() error {
	err := m.createProducer()
	if err != nil {
		return err
	}

	err = m.createConsumer()
	if err != nil {
		return err
	}
	return nil
}

// // shutdown生产者和消费者
// func (m *MessengerClient) DisConnect() error {
// 	if m.producer != nil {
// 		m.producer.Shutdown()
// 	}
// 	if m.consumer != nil {
// 		m.consumer.Shutdown()
// 	}
// 	return nil
// }

// 订阅
func (c *MessengerClient) Subscribe(ctx context.Context, req *SubscribeRequest, f ReceiveCallBack) error {
	// 创建主题
	tmp, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})))
	if err != nil {
		fmt.Printf("connection error: %s\n", err.Error())
	}
	err = tmp.CreateTopic(context.Background(), admin.WithTopicCreate(req.Topic), admin.WithBrokerAddrCreate("127.0.0.1:10911"))
	if err != nil {
		fmt.Printf("createTopic error: %s\n", err.Error())
	}

	if err := c.consumer.Subscribe(req.Topic, consumer.MessageSelector{},
		// 定义的消费函数
		func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, m := range me {
				err := f(ctx, &Message{Topic: m.Topic, Body: m.Body})
				if err != nil {
					return consumer.ConsumeRetryLater, err
				}
				// fmt.Printf("subscribe callback : %v \n", m)
			}
			return consumer.ConsumeSuccess, nil
		}); err != nil {
		return err
	}

	// 只有第一次注册的时候才需要启动
	c.consumerOnce.Do(func() {
		err = c.consumer.Start()
		if err != nil {
			logrus.Error("createConsumer Start failed : ", err.Error())
		}
	})
	if err != nil {
		return err
	}

	return nil
}

// 退订
func (c *MessengerClient) Unsubscribe(ctx context.Context, topic string) error {
	if err := c.consumer.Unsubscribe(topic); err != nil {
		return err
	}
	return nil
}

type PublishRequest struct {
	Topic string `json:"topic"`
	Body  string `json:"body"`
}

type PublishDelayRequest struct {
	PublishRequest
	//延迟消息级别
	// WithDelayTimeLevel set message delay time to consume.
	// reference delay level definition:
	// 1=>1s 2=>5s 3=>10s 4=>30s 5=>1m 6=>2m 7=>3m 8=>4m 9=>5m 10=>6m 11=>7m 12=>8m 13=>9m 14=>10m 15=>20m 16=>30m 17=>1h 18=>2h
	// delay level starts from 1. for example, if we set param level=1, then the delay time is 1s.
	DelayTimeLevel int64 `json:"delay_time_level" binding:"required"`
}

// 发布消息
func (c *MessengerClient) Publish(ctx context.Context, req *PublishRequest) (string, error) {
	res, err := c.producer.SendSync(ctx, primitive.NewMessage(req.Topic, []byte(req.Body)))
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

// 发布延迟消息
func (c *MessengerClient) PublishDelay(ctx context.Context, req *PublishDelayRequest) (string, error) {
	msg := primitive.NewMessage(req.Topic, []byte(req.Body))
	msg.WithDelayTimeLevel(int(req.DelayTimeLevel))
	res, err := c.producer.SendSync(ctx, msg)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
