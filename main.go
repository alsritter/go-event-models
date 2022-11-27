package main

import (
	"context"
	"time"

	"alsritter.icu/stgo/app/consumer"
	"alsritter.icu/stgo/domain"
	"alsritter.icu/stgo/domain/event"
	"alsritter.icu/stgo/infra/common"
	commonEvent "alsritter.icu/stgo/infra/common/event"
)

var consumers = []common.IConsumer{
	consumer.NewAutoCancelOrderConsumer(),
}

// RegisterEventHandler 注册事件处理器
func RegisterEventHandler(eventServer commonEvent.EventServerIface) error {
	// 注册所有消费者
	for _, c := range consumers {
		c.RegisterNeededEvent()
		if err := c.RegisterConsumer(eventServer); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := RegisterEventHandler(common.EventServer); err != nil {
		panic(err)
	}

	order := domain.Order{
		Id:                      "1234567",
		Name:                    "测试订单",
		Data:                    "这是一连串测试数据！！！！",
		AutoCancelCountdownTime: time.Now(),
	}

	time.Sleep(100 * time.Second)

	// 模拟发送事件
	for i := 0; i < 100000; i++ {
		common.EventPublish(context.TODO(), &event.OrderSubmittedEvent{
			OrderId: order.Id,
			Name:    order.Name,
			Data:    order.Data,
		})
	}

	time.Sleep(1000 * time.Second)
}

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	"github.com/apache/rocketmq-client-go/v2"
// 	"github.com/apache/rocketmq-client-go/v2/admin"
// 	"github.com/apache/rocketmq-client-go/v2/consumer"
// 	"github.com/apache/rocketmq-client-go/v2/primitive"
// 	"github.com/apache/rocketmq-client-go/v2/producer"
// )

// func main() {
// 	// 1. 创建主题，这一步可以省略，在send的时候如果没有topic，也会进行创建。
// 	CreateTopic("testTopic01")
// 	// 2.生产者向主题中发送消息
// 	// SendSyncMessage("hello world2022send test ，rocketmq go client!  too，是的")
// 	// 3.消费者订阅主题并消费
// 	SubcribeMessage()
// }

// func CreateTopic(topicName string) {
// 	endPoint := []string{"172.26.130.183:9876"}
// 	// 创建主题
// 	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(endPoint)))
// 	if err != nil {
// 		fmt.Printf("connection error: %s\n", err.Error())
// 	}
// 	err = testAdmin.CreateTopic(context.Background(), admin.WithTopicCreate(topicName), admin.WithBrokerAddrCreate("127.0.0.1:10911"))
// 	if err != nil {
// 		fmt.Printf("createTopic error: %s\n", err.Error())
// 	}
// }

// func SendSyncMessage(message string) {
// 	// 发送消息
// 	endPoint := []string{"172.26.130.183:9876"}
// 	// 创建一个producer实例
// 	p, _ := rocketmq.NewProducer(
// 		producer.WithNameServer(endPoint),
// 		producer.WithRetry(2),
// 		producer.WithGroupName("ProducerGroupName"),
// 	)
// 	// 启动
// 	err := p.Start()
// 	if err != nil {
// 		fmt.Printf("start producer error: %s", err.Error())
// 		os.Exit(1)
// 	}

// 	for i := 0; i < 100; i++ {
// 		// 发送消息
// 		result, err := p.SendSync(context.Background(), &primitive.Message{
// 			Topic: "testTopic01",
// 			Body:  []byte(message),
// 		})

// 		if err != nil {
// 			fmt.Printf("send message error: %s\n", err.Error())
// 		} else {
// 			fmt.Printf("send message seccess: result=%s\n", result.String())
// 		}
// 	}

// }

// func SubcribeMessage() {
// 	// 订阅主题、消费
// 	endPoint := []string{"172.26.130.183:9876"}
// 	// 创建一个consumer实例
// 	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer(endPoint),
// 		consumer.WithConsumerModel(consumer.Clustering),
// 		consumer.WithGroupName("ConsumerGroupName"),
// 	)
// 	if err != nil {
// 		fmt.Printf("start consumer error: %s", err.Error())
// 		os.Exit(1)
// 	}

// 	// 订阅topic
// 	err = c.Subscribe("testTopic01", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
// 		for i := range msgs {
// 			fmt.Printf("subscribe callback : %v \n", msgs[i])
// 		}
// 		return consumer.ConsumeSuccess, nil
// 	})

// 	if err != nil {
// 		fmt.Printf("subscribe message error: %s\n", err.Error())
// 	}

// 	// 启动consumer
// 	err = c.Start()

// 	if err != nil {
// 		fmt.Printf("consumer start error: %s\n", err.Error())
// 		os.Exit(-1)
// 	}

// 	err = c.Shutdown()
// 	if err != nil {
// 		fmt.Printf("shutdown Consumer error: %s\n", err.Error())
// 	}
// }
