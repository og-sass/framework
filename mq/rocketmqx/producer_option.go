package rocketmqx

import (
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
)

// ProducerOption 定义生产者选项
type ProducerOption func([]rmqClient.ProducerOption) []rmqClient.ProducerOption

// WithMaxAttempts 设置最大重试次数
func WithMaxAttempts(attempts int32) ProducerOption {
	return func(opts []rmqClient.ProducerOption) []rmqClient.ProducerOption {
		return append(opts, rmqClient.WithMaxAttempts(attempts))
	}
}

// WithTopics 预声明主题
func WithTopics(topics ...string) ProducerOption {
	return func(opts []rmqClient.ProducerOption) []rmqClient.ProducerOption {
		return append(opts, rmqClient.WithTopics(topics...))
	}
}

// WithTransactionChecker 设置事务检查器
func WithTransactionChecker(checker *rmqClient.TransactionChecker) ProducerOption {
	return func(opts []rmqClient.ProducerOption) []rmqClient.ProducerOption {
		return append(opts, rmqClient.WithTransactionChecker(checker))
	}
}
