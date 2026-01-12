package config

import (
	"context"

	"github.com/apache/rocketmq-clients/golang/v5"
)

type Config struct {
	ConsoleAppenderEnabled bool           `json:"console_appender_enabled,default=true"`
	LogLevel               string         `json:"log_level,default=info"`
	Endpoint               string         `json:"endpoint"`
	AccessKey              string         `json:"access_key,optional"`
	AccessSecret           string         `json:"access_secret,optional"`
	NameSpace              string         `json:"namespace,optional"`
	SecurityToken          string         `json:"security_token,optional"`
	ConsumerConfig         ConsumerConfig `json:"consumer_config,optional"`
}

type ConsumerConfig struct {
	ConsumerGroup              string        `json:"consumer_group"`
	AwaitDuration              int64         `json:"await_duration,default=5"`
	PullBatchSize              int           `json:"pull_batch_size,default=32"`
	InvisibleDuration          int64         `json:"invisible_duration,default=60"`
	PushConsumptionThreadCount int32         `json:"push_consumption_thread_count,default=20"`
	PushMaxCacheMessageCount   int32         `json:"push_max_cache_message_count,default=1024"`
	TopicRelations             TopicRelation `json:"topic_relations,optional"`
}

type TopicRelation struct {
	Topic          string `json:"topic"`
	Expression     string `json:"expression,default=*"`
	ExpressionType int    `json:"expression_type,default=1"`
}

type PullMessageHandler func(ctx context.Context, messages ...*golang.MessageView) (bool, error)
type PushMessageHandler func(message *golang.MessageView) golang.ConsumerResult
