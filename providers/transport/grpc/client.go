package grpc

import (
	"fmt"
	"github.com/aka-yz/go-micro-core"
	"github.com/aka-yz/go-micro-core/providers/transport/grpc/interceptors"
	"github.com/aka-yz/go-micro-core/providers/transport/grpc/selector"
	registry "github.com/aka-yz/go-micro-core/register"
	"github.com/aka-yz/go-micro-core/register/etcdv3"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
	"sync"
	"time"

	"go.uber.org/config"
)

type clientFactory struct{}

func (n *clientFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	if cfg := getRegistryConfig(conf); cfg != nil {
		return go_micro_core.NewProvider(newRPCClient(cfg))
	}
	return nil
}

func newRPCClient(options *registryConfig) *RPCClient {
	if options == nil {
		return nil
	}

	register := etcdv3.NewRegistry(
		registry.Addrs(options.Addrs...),
		registry.Timeout(time.Second*time.Duration(options.RegistryTTL)),
	)

	return NewClient(
		WithSuffix("-rpc"),
		WithSelector(selector.NewSelector(
			selector.Registry(register),
			selector.SetStrategy(selector.Random),
		)),
		WithInterceptor(
			interceptors.UnaryClientInterceptor(),
		),
	)
}

type RPCClient struct {
	connMap map[string]*serviceConn
	opts    ClientOptions
	sync.Mutex
}

func NewClient(opts ...ClientOption) *RPCClient {
	var opt ClientOptions
	for _, o := range opts {
		o(&opt)
	}

	return &RPCClient{
		opts:    opt,
		connMap: make(map[string]*serviceConn),
	}
}

func (c *RPCClient) GetConn(service string, opts ...ConnOption) (conn *grpc.ClientConn, err error) {
	target := service + c.opts.suffix
	c.Lock()
	defer c.Unlock()
	if sc, ok := c.connMap[target]; ok {
		return sc.getConn(), nil
	}

	// 不存在连接
	sc := NewServiceConn(target, opts...)
	if conn, err = sc.CreateConn(c.opts.selector.Options().Registry, c.opts.interceptors); err == nil {
		c.connMap[target] = sc
	}
	return
}

func (c *RPCClient) AddInterceptorsTail(interceptors ...grpc.UnaryClientInterceptor) {
	c.opts.interceptors = append(c.opts.interceptors, interceptors...)
}

func (c *RPCClient) AddInterceptorsHead(interceptors ...grpc.UnaryClientInterceptor) {
	c.opts.interceptors = append(interceptors, c.opts.interceptors...)
}

type serviceConn struct {
	conns  []*grpc.ClientConn
	target string
	opts   ConnOptions
}

func NewServiceConn(target string, opts ...ConnOption) *serviceConn {
	var opt ConnOptions
	for _, o := range opts {
		o(&opt)
	}
	return &serviceConn{
		target: target,
		opts:   opt,
	}
}

func (c *serviceConn) CreateConn(registry registry.Registry, inters []grpc.UnaryClientInterceptor) (conn *grpc.ClientConn, err error) {
	c.opts.interceptors = append(c.opts.interceptors, inters...)

	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(middleware.ChainUnaryClient(c.opts.interceptors...)),
	}
	if c.opts.block {
		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if c.opts.timeout != 0 {
		dialOptions = append(dialOptions, grpc.WithTimeout(c.opts.timeout))
	}

	if c.opts.maxSize != 0 {
		dialOptions = append(dialOptions, grpc.WithMaxMsgSize(c.opts.maxSize))
	}

	if c.opts.balanceName == "" {
		dialOptions = append(dialOptions, grpc.WithBalancer(grpc.RoundRobin(&resolver{registry: registry})))
	} else {
		dialOptions = append(dialOptions, grpc.WithBalancerName(c.opts.balanceName))
	}

	conn, err = grpc.Dial(c.target, dialOptions...)
	if err != nil {
		return
	}

	c.conns = append(c.conns, conn)
	return
}

func (c *serviceConn) getConn() *grpc.ClientConn {
	return c.conns[0]
}

type resolver struct {
	registry registry.Registry
}

func (r *resolver) Resolve(target string) (naming.Watcher, error) {
	w, err := r.registry.Watch(target)
	if err != nil {
		return nil, err
	}
	return &watcher{
		w:      w,
		r:      r.registry,
		target: target,
	}, nil
}

type watcher struct {
	isinitialized bool
	target        string
	w             registry.Watcher
	r             registry.Registry
}

func (w *watcher) Next() (updates []*naming.Update, err error) {
	if !w.isinitialized {
		ss, err := w.r.GetService(w.target)
		if err != nil {
			return nil, err
		}

		for _, s := range ss {
			updates = append(updates, w.update("create", s)...)
		}
		w.isinitialized = true
		return updates, nil
	}

	ret, err := w.w.Next()
	if err != nil {
		return
	}

	return w.update(ret.Action, ret.Service), nil
}

// Close closes the Watcher.
func (w *watcher) Close() {
	w.w.Stop()
}

func (w *watcher) update(action string, service *registry.Service) (updates []*naming.Update) {
	var op naming.Operation
	switch action {
	case "create":
		op = naming.Add
	case "delete":
		op = naming.Delete
	}

	for _, node := range service.Nodes {
		updates = append(updates, &naming.Update{Op: op, Addr: fmt.Sprintf("%s:%d", node.Address, node.Port)})
	}
	return
}
