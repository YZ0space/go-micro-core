package etcdv3

import (
	"context"
	"encoding/json"
	"errors"
	registry "github.com/aka-yz/go-micro-core/register"
	clientv3 "go.etcd.io/etcd/client/v3"
	"path"
	"strings"
	"time"
)

// register  etcd 服务注册 （预留）
type etcdv3Registry struct {
	Client *clientv3.Client
	opt    registry.Options
	Ctx    context.Context
	stop   chan bool
	leaser clientv3.Lease
}

const (
	prefix = "/services/"
)

func nodePath(name string, id string) string {
	return path.Join(servicePath(name), id)
}

func servicePath(name string) string {
	return path.Join(prefix, strings.Replace(name, "-", "/", -1))
}

func encode(service *registry.Service) string {
	val, _ := json.Marshal(service)
	return string(val)
}

func decode(val []byte) (ret *registry.Service) {
	json.Unmarshal(val, &ret)
	return
}

// Register 一个client 对应一个服务
func (e *etcdv3Registry) Register(service *registry.Service, opt ...registry.RegisterOption) (err error) {
	if len(service.Nodes) == 0 {
		return errors.New("service nodes empty")
	}

	registerOptions := registry.RegisterOptions{
		TTL: 30 * time.Second,
	}
	for _, o := range opt {
		o(&registerOptions)
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.opt.Timeout)
	defer cancel()
	grantResp, err := e.leaser.Grant(ctx, int64(registerOptions.TTL.Seconds()))
	if err != nil {
		return err
	}

	if _, err = e.leaser.KeepAliveOnce(e.Ctx, grantResp.ID); err != nil {
		return
	}

	_, err = e.Client.Put(ctx, nodePath(service.Name, service.Nodes[0].Id), encode(service), clientv3.WithLease(grantResp.ID))
	if err != nil {
		return
	}

	return
}

func (e *etcdv3Registry) Deregister(service *registry.Service) (err error) {
	if len(service.Nodes) == 0 {
		return errors.New("service nodes empty")
	}

	if _, err = e.Client.Delete(e.Ctx, nodePath(service.Name, service.Nodes[0].Id), clientv3.WithPrefix()); err != nil {
		return
	}

	if err = e.leaser.Close(); err != nil {
		return
	}

	close(e.stop)
	return
}

// GetService 通过name获取所有注册的服务
func (e *etcdv3Registry) GetService(name string) (services []*registry.Service, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.opt.Timeout)
	defer cancel()
	getResp, err := e.Client.Get(ctx, servicePath(name)+"/", clientv3.WithPrefix())
	if err != nil {
		return
	}

	// 获取到注册val 反序列化成service
	for _, kv := range getResp.Kvs {
		if sv := decode(kv.Value); sv != nil {
			services = append(services, sv)
		}
	}
	return
}

// ListServices 获取client所有服务
func (e *etcdv3Registry) ListServices() (services []*registry.Service, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.opt.Timeout)
	defer cancel()
	getResp, err := e.Client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return
	}

	// 获取到注册val 反序列化成service
	for _, kv := range getResp.Kvs {
		if sv := decode(kv.Value); sv != nil {
			services = append(services, sv)
		}
	}
	return
}

// Watch 监听服务变化
func (e *etcdv3Registry) Watch(service string) (w registry.Watcher, err error) {
	return newEtcdV3Watcher(e, servicePath(service)), nil
}

func (e *etcdv3Registry) String() string {
	return "etcdv3"
}

func (e *etcdv3Registry) Options() registry.Options {
	return e.opt
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	var opt registry.Options
	for _, o := range opts {
		o(&opt)
	}

	if len(opt.Addrs) == 0 {
		return nil
	}

	cfg := clientv3.Config{
		Endpoints: opt.Addrs,
	}

	ctx, cancel := context.WithCancel(context.TODO())
	// .... 不支持TLS
	client, _ := clientv3.New(cfg)
	regist := &etcdv3Registry{
		Client: client,
		Ctx:    ctx,
		opt:    opt,
		leaser: clientv3.NewLease(client),
		stop:   make(chan bool),
	}

	go func() {
		<-regist.stop
		cancel()
	}()

	return regist
}
