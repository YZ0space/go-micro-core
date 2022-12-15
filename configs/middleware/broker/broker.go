package broker

import (
	"context"
)

type Broker interface {
	BatchPublish(topic string, messages []interface{}, opts ...PublishOption) (err error)
	BatchPublishWithPartition(topic string, messages map[string]interface{}, opts ...PublishOption) (err error)
	BatchPublishWithPartitionWithCtx(ctx context.Context, topic string, messages map[string]interface{}, opts ...PublishOption) (err error)
	Publish(topic string, message interface{}, opts ...PublishOption) (err error)
	PublishWithCtx(ctx context.Context, topic string, message interface{}, opts ...PublishOption) (err error)
	Subscribe(topic string, channel string, handler HandlerFunc, opts ...SubscribeOption) (err error)
	UnSubscribe()(err error)
	AsyncPublish(topic string, message interface{}, opts ...PublishOption) (err error)
	KafkaAsyncProducer()
	String() string
}

type HandlerFunc func(ctx context.Context, data []byte) error
