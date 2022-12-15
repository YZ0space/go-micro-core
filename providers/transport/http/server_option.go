package http

type ServerOptions struct {
	// todo: 预留的服务注册
	//registry registry.Registry
	//service  *registry.Service
	addr string
}

type ServerOption func(*ServerOptions)
