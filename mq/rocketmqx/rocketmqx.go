package rocketmqx

import (
	"context"
	"os"
	"strings"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	v2 "github.com/apache/rocketmq-clients/golang/v5/protocol/v2"
	config2 "github.com/og-sass/framework/mq/rocketmqx/config"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type RocketMqx struct {
	config config2.Config
}

func NewRocketMqx(config config2.Config) *RocketMqx {
	_ = os.Setenv("mq.consoleAppender.enabled", lo.Ternary(config.ConsoleAppenderEnabled, "true", "false"))
	rmq_client.ResetLogger()
	return &RocketMqx{config: config}
}

// 创建基础配置，减少重复代码
func (r *RocketMqx) createBaseConfig() *rmq_client.Config {
	return &rmq_client.Config{
		Endpoint:      r.config.Endpoint,
		NameSpace:     r.config.NameSpace,
		ConsumerGroup: r.config.ConsumerConfig.ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:     r.config.AccessKey,
			AccessSecret:  r.config.AccessSecret,
			SecurityToken: r.config.SecurityToken,
		},
	}
}

func (r *RocketMqx) NewProducer(options ...ProducerOption) (producer rmq_client.Producer, err error) {

	var rocketmqOpts []rmq_client.ProducerOption
	for _, opt := range options {
		rocketmqOpts = opt(rocketmqOpts)
	}

	producer, err = rmq_client.NewProducer(
		r.createBaseConfig(),
		rocketmqOpts...,
	)
	if err != nil {
		logx.Errorf("NewProducer err: %s", err.Error())
		return
	}
	// 启动生产者
	if err = producer.Start(); err != nil {
		logx.Errorf("Start producer err: %s", err.Error())
		return
	}
	return
}

func (r *RocketMqx) NewPullConsumer(handler config2.PullMessageHandler) (simpleConsumer rmq_client.SimpleConsumer, err error) {
	relations := map[string]*rmq_client.FilterExpression{
		r.config.ConsumerConfig.TopicRelations.Topic: rmq_client.NewFilterExpressionWithType(
			r.config.ConsumerConfig.TopicRelations.Expression,
			rmq_client.FilterExpressionType(r.config.ConsumerConfig.TopicRelations.ExpressionType),
		),
	}

	simpleConsumer, err = rmq_client.NewSimpleConsumer(
		r.createBaseConfig(),
		rmq_client.WithSimpleAwaitDuration(time.Duration(r.config.ConsumerConfig.AwaitDuration)*time.Second),
		rmq_client.WithSimpleSubscriptionExpressions(relations),
	)
	if err != nil {
		logx.Errorf("初始化消费者失败，原因为：%s", err.Error())
		return
	}

	if err = simpleConsumer.Start(); err != nil {
		logx.Errorf("启动消费者失败，原因为：%s", err.Error())
		return
	}

	// 将消息处理逻辑提取到单独的函数中
	go r.processMessages(simpleConsumer, handler, r.config.ConsumerConfig.TopicRelations.Topic)

	return
}

func (r *RocketMqx) NewPushConsumer(handler func(*rmq_client.MessageView) rmq_client.ConsumerResult) (pushConsumer rmq_client.PushConsumer, err error) {
	relations := map[string]*rmq_client.FilterExpression{
		r.config.ConsumerConfig.TopicRelations.Topic: rmq_client.NewFilterExpressionWithType(
			r.config.ConsumerConfig.TopicRelations.Expression,
			rmq_client.FilterExpressionType(r.config.ConsumerConfig.TopicRelations.ExpressionType),
		),
	}
	// In most case, you don't need to create many consumers, singleton pattern is more recommended.
	pushConsumer, err = rmq_client.NewPushConsumer(r.createBaseConfig(),
		rmq_client.WithPushAwaitDuration(time.Duration(r.config.ConsumerConfig.AwaitDuration)*time.Second),
		rmq_client.WithPushSubscriptionExpressions(relations),
		rmq_client.WithPushMessageListener(&rmq_client.FuncMessageListener{
			Consume: handler,
		}),
		rmq_client.WithPushConsumptionThreadCount(20),
		rmq_client.WithPushMaxCacheMessageCount(1024),
	)
	if err != nil {
		logx.Errorf("NewPushConsumer err: %s", err.Error())
		return
	}
	// start pushConsumer
	if err = pushConsumer.Start(); err != nil {
		logx.Errorf("Start pushConsumer err: %s", err.Error())
		return
	}
	return
}

// 处理消息的逻辑提取到单独的函数中
func (r *RocketMqx) processMessages(consumer rmq_client.SimpleConsumer, handler config2.PullMessageHandler, topic string) {
	for {
		// 1. 拉取消息 - Receive超时设置为 AwaitDuration + 5秒buffer
		receiveCtx, receiveCancel := context.WithTimeout(
			context.Background(),
			time.Duration(r.config.ConsumerConfig.AwaitDuration+5)*time.Second,
		)
		mvs, err := consumer.Receive(
			receiveCtx,
			int32(r.config.ConsumerConfig.PullBatchSize),
			time.Duration(r.config.ConsumerConfig.InvisibleDuration)*time.Second,
		)
		receiveCancel()

		// 2. 处理拉取错误
		if err != nil {
			if strings.Contains(err.Error(), v2.Code_name[int32(v2.Code_MESSAGE_NOT_FOUND)]) {
				// 无消息时短暂休眠
				time.Sleep(time.Millisecond * 200)
				continue
			}
			logx.Errorf("拉取消息失败，topic:%s,原因为:%s", topic, err.Error())
			continue
		}

		// 3. 处理消息 - 使用 InvisibleDuration 作为处理超时
		handlerCtx, handlerCancel := context.WithTimeout(
			context.Background(),
			time.Duration(r.config.ConsumerConfig.InvisibleDuration)*time.Second,
		)

		res, err := handler(handlerCtx, mvs...)
		handlerCancel()

		// 4. ACK确认
		if res && err == nil {
			// 确认ACK - 5秒超时
			ackCtx, ackCancel := context.WithTimeout(context.Background(), 5*time.Second)
			for _, mv := range mvs {
				if ackErr := consumer.Ack(ackCtx, mv); ackErr != nil {
					logx.Errorf("ack message failed, reason: %s, msgID:%s", ackErr.Error(), mv.GetMessageId())
				}
			}
			ackCancel()
		} else if err != nil {
			logx.Errorf("处理消息失败,topic:%s,原因为：%s", topic, err.Error())
		}
	}
}

// ProducerOption 定义生产者选项
type ProducerOption func([]rmq_client.ProducerOption) []rmq_client.ProducerOption

// WithMaxAttempts 设置最大重试次数
func WithMaxAttempts(attempts int32) ProducerOption {
	return func(opts []rmq_client.ProducerOption) []rmq_client.ProducerOption {
		return append(opts, rmq_client.WithMaxAttempts(attempts))
	}
}

// WithTopics 预声明主题
func WithTopics(topics ...string) ProducerOption {
	return func(opts []rmq_client.ProducerOption) []rmq_client.ProducerOption {
		return append(opts, rmq_client.WithTopics(topics...))
	}
}

// WithTransactionChecker 设置事务检查器
func WithTransactionChecker(checker *rmq_client.TransactionChecker) ProducerOption {
	return func(opts []rmq_client.ProducerOption) []rmq_client.ProducerOption {
		return append(opts, rmq_client.WithTransactionChecker(checker))
	}
}
