package broker

import (
	"github.com/Shopify/sarama"
	"time"
)

type Options struct {
	Addrs             []string // producer addrs
	LookUpAddrs       []string // consumer addrs for nsq
	Suffix            string
	Timeout           int64
	MaxMessageBytes   int
	MaxRequestSize    int32
	MaxProcessingTime int64
	Username          string
	Password          string
}

type Option func(*Options)

// Addrs sets the host addresses to be used by the broker
func WithAddrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func WithTimeout(timeout int64) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

func WithMaxMessageBytes(maxMessageBytes int) Option {
	return func(o *Options) {
		o.MaxMessageBytes = maxMessageBytes
	}
}

func WithMaxRequestSize(maxRequestSize int32) Option {
	return func(o *Options) {
		o.MaxRequestSize = maxRequestSize
	}
}

func WithMaxProcessingTime(maxProcessingTime int64) Option {
	return func(o *Options) {
		o.MaxProcessingTime = maxProcessingTime
	}
}

func WithUsername(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

// Addrs sets the host addresses to be used by the broker
func WithLookUpAddrs(addrs ...string) Option {
	return func(o *Options) {
		o.LookUpAddrs = addrs
	}
}

func Suffix(suffix string) Option {
	return func(o *Options) {
		o.Suffix = suffix
	}
}

type PublishOptions struct {
	DelaySec time.Duration
	Key      string
	Headers  map[string][]sarama.RecordHeader
}

type PublishOption func(*PublishOptions)

func WithDelaySecs(sec time.Duration) PublishOption {
	return func(o *PublishOptions) {
		o.DelaySec = sec
	}
}

// for kafka
func WithKey(key string) PublishOption {
	return func(o *PublishOptions) {
		o.Key = key
	}
}

func WithHeaders(headers map[string][]sarama.RecordHeader) PublishOption {
	return func(o *PublishOptions) {
		o.Headers = headers
	}
}

type SubscribeOptions struct {
	Count int
	Addrs []string

	Partition bool
	Offset    int64
}

type SubscribeOption func(*SubscribeOptions)

func WithConcurrencyCount(cnt int) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Count = cnt
	}
}

// consumer find topic by addrs
func WithConsumerAddrs(addrs ...string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Addrs = addrs
	}
}

func WithPartition(partition bool) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Partition = partition
	}
}

func WithInitOffset(offset int64) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Offset = offset
	}
}

const (
	Partition = "partition"
	Topic     = "topic"
	GroupID   = "group-id"
)

type HandlerOptionKey struct{}

type HandlerOptions struct {
	ResetOption ResetOption
}

type ResetOption struct {
	NeedReset      bool  `json:"need_reset"`
	ResetTimestamp int64 `json:"reset_timestamp"`
}
