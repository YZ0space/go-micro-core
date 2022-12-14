package grpc

import (
	"github.com/aka-yz/go-micro-core/providers/transport/grpc/selector"
	"google.golang.org/grpc"
	"time"
)

type ClientOptions struct {
	suffix       string
	selector     selector.Selector
	interceptors []grpc.UnaryClientInterceptor
}

type ClientOption func(*ClientOptions)

func WithSuffix(suffix string) ClientOption {
	return func(o *ClientOptions) {
		o.suffix = suffix
	}
}

func WithInterceptor(interceptors ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *ClientOptions) {
		o.interceptors = interceptors
	}
}

func WithSelector(selector selector.Selector) ClientOption {
	return func(o *ClientOptions) {
		o.selector = selector
	}
}

type ConnOptions struct {
	block        bool
	connNum      int64
	balanceName  string
	maxSize      int
	timeout      time.Duration
	dials        []grpc.DialOption
	interceptors []grpc.UnaryClientInterceptor
}

type ConnOption func(*ConnOptions)

func WithBlock(block bool) ConnOption {
	return func(o *ConnOptions) {
		o.block = block
	}
}

func WithConnNum(count int64) ConnOption {
	return func(o *ConnOptions) {
		o.connNum = count
	}
}

func WithConnInterceptor(interceptors ...grpc.UnaryClientInterceptor) ConnOption {
	return func(o *ConnOptions) {
		o.interceptors = interceptors
	}
}

func WithBalanceName(name string) ConnOption {
	return func(o *ConnOptions) {
		o.balanceName = name
	}
}

func WithTimeout(d time.Duration) ConnOption {
	return func(o *ConnOptions) {
		o.timeout = d
	}
}

func WithMaxSize(s int) ConnOption {
	return func(o *ConnOptions) {
		o.maxSize = s
	}
}
