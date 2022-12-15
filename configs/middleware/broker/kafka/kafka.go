package kafka

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/aka-yz/go-micro-core/configs/log"
	"github.com/aka-yz/go-micro-core/configs/middleware/broker"
	"hash"
	"runtime"
	"time"

	"github.com/Shopify/sarama"
)

type Kafka struct {
	p    sarama.SyncProducer
	ap   sarama.AsyncProducer
	c    sarama.Consumer
	opts broker.Options
}

func (k *Kafka) Publish(topic string, message interface{}, opts ...broker.PublishOption) (err error) {
	var opt broker.PublishOptions
	for _, o := range opts {
		o(&opt)
	}

	msg, err := getMessage(message)
	if err != nil {
		return
	}

	partition, offset, err := k.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(opt.Key),
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		log.Infof(context.TODO(), "msg:%v partition:%v offset:%v err:%v", string(msg), partition, offset, err)
		return
	}

	log.Debugf(context.TODO(), "kafka publish success topic:%v msg:%v", topic, msg)
	return
}

// PublishWithCtx 发布msg
func (k *Kafka) PublishWithCtx(ctx context.Context, topic string, message interface{}, opts ...broker.PublishOption) (err error) {
	var opt broker.PublishOptions
	for _, o := range opts {
		o(&opt)
	}

	msg, err := getMessage(message)
	if err != nil {
		return
	}

	partition, offset, err := k.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(opt.Key),
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		log.Infof(ctx, "msg:%v partition:%v offset:%v err:%v", string(msg), partition, offset, err)
		return
	}
	log.Debugf(ctx, "kafka publish success topic:%v msg:%v", topic, msg)
	return
}

func getMessage(message interface{}) (msg []byte, err error) {
	if m, ok := message.(string); ok {
		return []byte(m), nil
	}

	return json.Marshal(message)
}

func (k *Kafka) BatchPublish(topic string, messages []interface{}, opts ...broker.PublishOption) (err error) {
	var opt broker.PublishOptions
	for _, o := range opts {
		o(&opt)
	}

	msgs, err := getMessageBatch(messages, topic, opt)
	if err != nil || len(msgs) == 0 {
		return
	}
	err = k.p.SendMessages(msgs)
	if err != nil {
		log.Infof(context.TODO(), "msgs len:%v partition:%v offset:%v err:%v", len(msgs), err)
		return
	}

	log.Debugf(context.TODO(), "kafka publish success topic:%v msg len:%v", topic, len(msgs))
	return
}

func (k *Kafka) BatchPublishWithPartition(topic string, messages map[string]interface{}, opts ...broker.PublishOption) (err error) {
	var opt broker.PublishOptions
	for _, o := range opts {
		o(&opt)
	}

	msgs, err := getMessageBatchWithPartition(messages, topic, opt)
	if err != nil || len(msgs) == 0 {
		return
	}
	err = k.p.SendMessages(msgs)
	if err != nil {
		log.Infof(context.TODO(), "msgs len:%v partition:%v offset:%v err:%v", len(msgs), err)
		return
	}

	log.Debugf(context.TODO(), "kafka publish success topic:%v msg len:%v", topic, len(msgs))
	return
}

// BatchPublishWithPartitionWithCtx 批量发布消息
func (k *Kafka) BatchPublishWithPartitionWithCtx(ctx context.Context, topic string, messages map[string]interface{}, opts ...broker.PublishOption) (err error) {
	var opt broker.PublishOptions
	for _, o := range opts {
		o(&opt)
	}

	msgs, err := getMessageBatchWithPartitionWithTrace(messages, topic, opt)
	if err != nil || len(msgs) == 0 {
		return
	}
	err = k.p.SendMessages(msgs)
	if err != nil {
		log.Infof(ctx, "msgs len:%v partition:%v offset:%v err:%v", len(msgs), err)
		return
	}

	log.Debugf(ctx, "kafka publish success topic:%v msg len:%v", topic, len(msgs))
	return
}

func getMessageBatchWithPartitionWithTrace(messages map[string]interface{}, topic string, opt broker.PublishOptions) (msgs []*sarama.ProducerMessage, err error) {
	for key, message := range messages {
		msg, err := getMessage(message)
		if err != nil {
			return nil, err
		}
		var msgHeader []sarama.RecordHeader
		if opt.Headers != nil {
			msgHeader = append(msgHeader, opt.Headers[key]...)
		}

		sendMsg := &sarama.ProducerMessage{
			Topic:   topic,
			Key:     sarama.StringEncoder(key),
			Value:   sarama.ByteEncoder(msg),
			Headers: msgHeader,
		}
		msgs = append(msgs, sendMsg)
	}
	return
}

func getMessageBatch(messages []interface{}, topic string, opt broker.PublishOptions) (msgs []*sarama.ProducerMessage, err error) {
	for _, message := range messages {
		msg, err := getMessage(message)
		if err != nil {
			return nil, err
		}
		sendMsg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(opt.Key),
			Value: sarama.ByteEncoder(msg),
		}
		msgs = append(msgs, sendMsg)
	}
	return
}

func getMessageBatchWithPartition(messages map[string]interface{}, topic string, opt broker.PublishOptions) (msgs []*sarama.ProducerMessage, err error) {
	for key, message := range messages {
		msg, err := getMessage(message)
		if err != nil {
			return nil, err
		}

		sendMsg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.ByteEncoder(msg),
		}

		if opt.Headers != nil {
			sendMsg.Headers = opt.Headers[key]
		}

		msgs = append(msgs, sendMsg)
	}
	return
}

func (k *Kafka) AsyncPublish(topic string, message interface{}, opts ...broker.PublishOption) (err error) {
	var opt broker.PublishOptions
	for _, o := range opts {
		o(&opt)
	}

	msg, err := getMessage(message)
	if err != nil {
		return
	}

	sendMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(opt.Key),
		Value: sarama.ByteEncoder(msg),
	}
	k.ap.Input() <- sendMsg

	log.Debugf(context.TODO(), "kafka asyncReport send message success topic:%v msg:%v", topic, msg)
	return
}

func (k *Kafka) KafkaAsyncProducer() {
	success := k.ap.Successes()
	errors := k.ap.Errors()
	for {
		select {
		case err := <-errors:
			if err != nil {
				log.Errorf(context.TODO(), "asyncProducer error: %+v", err)
			}
		case <-success:
			log.Debugf(context.TODO(), "asyncProducer success.")
		}
	}
}

func (k *Kafka) Subscribe(topic string, channel string, handler broker.HandlerFunc, opts ...broker.SubscribeOption) (err error) {
	var opt broker.SubscribeOptions
	for _, o := range opts {
		o(&opt)
	}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	client, err := sarama.NewClient(opt.Addrs, config)
	defer client.Close()
	if err != nil {
		panic(err)
	}
	k.c, err = k.getConsumer(client, channel)
	if err != nil {
		log.Errorf(context.TODO(), "kafka get consumer err:%v\n", err)
		panic(err)
	}
	partitions, err := k.c.Partitions(topic)
	if err != nil {
		log.Errorf(context.TODO(), "kafka partition get err:%v\n", err)
		panic(err)
	}

	for _, partitionId := range partitions {
		partitionConsumer, err := k.c.ConsumePartition(topic, partitionId, sarama.OffsetNewest)
		if err != nil {
			log.Errorf(context.TODO(), "kafka consumer partition partitionId:%v, err:%v\n", partitionId, err)
			panic(err)
		}

		go func(pc sarama.PartitionConsumer) {
			defer func() {
				if pErr := recover(); pErr != nil {
					//打印调用栈信息
					buf := make([]byte, 8192)
					n := runtime.Stack(buf, false)
					stackInfo := fmt.Sprintf("%s", buf[:n])
					log.Errorf(context.Background(), "err: %v, panic stack info %s", pErr, stackInfo)
				}
			}()
			for msg := range pc.Messages() {
				newCtx := context.TODO()
				log.Infof(newCtx, "partitionId: %d; offset:%d, value: %s\n", msg.Partition, msg.Offset, msg.Value)
				if err = handler(newCtx, msg.Value); err != nil {
					log.Errorf(newCtx, "topic:%v msg:%v handle error:%v\n", msg.Topic, string(msg.Value), err)
					continue
				}
			}
		}(partitionConsumer)
	}
	return
}

func (k *Kafka) getConsumer(client sarama.Client, channel string) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	defer consumer.Close()
	if err != nil {
		log.Errorf(context.TODO(), "kafka consumer close err:%v", err)
		panic(err)
	}
	if err := client.RefreshCoordinator(channel); err != nil {
		return nil, err
	}
	return consumer, nil
}

func (k *Kafka) UnSubscribe() (err error) {
	if errClose := k.c.Close(); errClose != nil {
		err = errClose
		topics, err := k.c.Topics()
		log.Errorf(context.TODO(), "kafka unsubscribe error, consumer: %+v, err: %v", topics, err)
	}
	return
}

func (k *Kafka) String() string {
	return "kafka"
}

func NewBroker(opts ...broker.Option) broker.Broker {
	var opt broker.Options
	for _, o := range opts {
		o(&opt)
	}
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Version = sarama.V1_0_0_0
	if opt.Timeout > 0 {
		kafkaConfig.Producer.Timeout = time.Duration(opt.Timeout) * time.Second
	}
	if opt.MaxMessageBytes > 0 {
		kafkaConfig.Producer.MaxMessageBytes = opt.MaxMessageBytes
	}
	if opt.MaxRequestSize > 0 {
		sarama.MaxRequestSize = opt.MaxRequestSize
	}
	if opt.MaxProcessingTime > 0 {
		kafkaConfig.Consumer.MaxProcessingTime = time.Duration(opt.MaxProcessingTime * int64(time.Second))
	}

	if opt.Username != "" && opt.Password != "" {
		withSASL(kafkaConfig, opt.Username, opt.Password)
	}

	p, err := sarama.NewSyncProducer(opt.Addrs, kafkaConfig)
	if err != nil {
		panic(err)
	}

	ap, err := sarama.NewAsyncProducer(opt.Addrs, kafkaConfig)
	if err != nil {
		panic(err)
	}
	return &Kafka{
		opts: opt,
		p:    p,
		ap:   ap,
	}
}

func withSASL(kafkaConfig *sarama.Config, username, password string) {
	kafkaConfig.Version = sarama.V2_3_0_0
	kafkaConfig.Net.SASL.Enable = true
	kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	kafkaConfig.Net.SASL.User = username
	kafkaConfig.Net.SASL.Password = password
	SHA512 := func() hash.Hash { return sha512.New() }
	kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
}
