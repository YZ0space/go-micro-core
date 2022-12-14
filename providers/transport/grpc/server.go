package grpc

import (
	"fmt"
	go_micro_core "github.com/aka-yz/go-micro-core"
	grpc_interceptors "github.com/aka-yz/go-micro-core/providers/transport/grpc/interceptors"
	registry "github.com/aka-yz/go-micro-core/register"
	"github.com/aka-yz/go-micro-core/register/etcdv3"
	"github.com/aka-yz/go-micro-core/utils/json"
	netutils "github.com/aka-yz/go-micro-core/utils/net"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type serverFactory struct{}

func (s *serverFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	if cfg := getServerConfig(conf); cfg != nil {
		return go_micro_core.NewProvider(
			reflectRPCServer(newRPCServer(cfg)))
	}
	return nil
}

func reflectRPCServer(s *RPCServer) *RPCServer {
	reflection.Register(s.Server)
	return s
}

func newRPCServer(cfg *serverConfig) *RPCServer {
	interceptors := []grpc.UnaryServerInterceptor{
		grpc_interceptors.UnaryServerInterceptor(),
		recovery.UnaryServerInterceptor(),
	}

	if cfg.Metadata == "gin" {
		interceptors = append(interceptors, grpc_interceptors.GinUnaryServerInterceptor())
	}

	var register registry.Registry
	if cfg.Registry != nil {
		register = etcdv3.NewRegistry(
			registry.Addrs(cfg.Registry.Addrs...),
			registry.Timeout(time.Second*time.Duration(cfg.Registry.RegistryTTL)),
		)
	}

	server := NewServer(
		Addr(cfg.Addr),
		Service(cfg.Service),
		Registry(register),
		GRPCServerOption(
			grpc.UnaryInterceptor(middleware.ChainUnaryServer(interceptors...)),
		),
	)

	return server
}

type serverConfig struct {
	Addr     string
	Metadata string
	Registry *registryConfig
	Service  *registry.Service
}

func getServerConfig(conf config.Provider) *serverConfig {
	var cv config.Value
	if cv = conf.Get("rpcserver"); !cv.HasValue() {
		return nil
	}

	addrMap := make(map[string]string)
	if err := cv.Populate(&addrMap); err != nil {
		return nil
	}

	var cfg serverConfig
	cfg.Addr = port(addrMap["addr"])
	cfg.Metadata = addrMap["metadata"]
	cfg.Registry = getRegistryConfig(conf)
	cfg.Service = getRegistryService(conf, "-rpc")
	return &cfg
}

func port(addr string) string {
	port := os.Getenv("PORT_" + addr)
	if port != "" {
		return ":" + port
	}

	port = os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}

	return addr
}

// RPCServer 启动服务
// 注册registry
type RPCServer struct {
	*grpc.Server
	opts ServerOptions
}

func NewServer(opts ...ServerOption) *RPCServer {
	var opt ServerOptions
	for _, o := range opts {
		o(&opt)
	}

	rs := &RPCServer{
		opts:   opt,
		Server: grpc.NewServer(opt.serverOptions...),
	}

	//hs := health.NewServer()
	//proto.RegisterHealthServer(rs.Server, hs)
	return rs
}

// Start 启动服务
func (s *RPCServer) Start() {
	ip, ls, err := netutils.ListenAddr(s.opts.addr, func(addr string) (net.Listener, error) {
		return net.Listen("tcp", s.opts.addr)
	})
	if err != nil {
		panic(err)
	}

	s.opts.addr = ip + s.opts.addr
	log.Println("RPCServer listen on:", s.opts.addr)
	go s.Serve(ls)
	go s.register()
}

func (s *RPCServer) register() {
	if s.opts.registry == nil {
		return
	}
	addr := strings.Split(s.opts.addr, ":")
	if len(addr) != 2 {
		panic(fmt.Errorf("error register addr:%v", addr))
	}
	port, _ := strconv.Atoi(addr[1])
	s.opts.service.Nodes[0].Address = addr[0]
	s.opts.service.Nodes[0].Port = port

	for {
		if err := s.opts.registry.Register(s.opts.service); err == nil {
			log.Println("RPC Server register:", json.MustString(s.opts.service))
		}

		time.Sleep(time.Second * 15)
	}
}

func (s *RPCServer) Stop() {
	s.Server.GracefulStop()
	if s.opts.registry == nil {
		return
	}
	if err := s.opts.registry.Deregister(s.opts.service); err != nil {
		log.Printf("Deregister failed service:%v error:%v", json.MustString(s.opts.service), err)
	}
}
