package grpc

import (
	registry "github.com/aka-yz/go-micro-core/register"
	"google.golang.org/grpc"
)

type ServerOptions struct {
	registry      registry.Registry
	service       *registry.Service
	addr          string
	serverOptions []grpc.ServerOption
	interceptors  []grpc.UnaryServerInterceptor
}

type ServerOption func(*ServerOptions)

func UnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.interceptors = interceptors
	}
}

func Registry(registry registry.Registry) ServerOption {
	return func(o *ServerOptions) {
		o.registry = registry
	}
}

func Service(s *registry.Service) ServerOption {
	return func(o *ServerOptions) {
		o.service = s
	}
}

func Addr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.addr = addr
	}
}

func GRPCServerOption(opts ...grpc.ServerOption) ServerOption {
	return func(o *ServerOptions) {
		o.serverOptions = opts
	}
}
