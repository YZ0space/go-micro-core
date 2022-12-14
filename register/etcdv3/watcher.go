package etcdv3

import (
	"context"
	"errors"
	"github.com/aka-yz/go-micro-core/configs/log"
	registry "github.com/aka-yz/go-micro-core/register"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type watcher struct {
	watchChan clientv3.WatchChan
}

func newEtcdV3Watcher(e *etcdv3Registry, path string) registry.Watcher {
	return &watcher{
		watchChan: e.Client.Watch(e.Ctx, path, clientv3.WithPrefix(), clientv3.WithPrevKV()),
	}
}

// Next 监听服务变化
func (w *watcher) Next() (result *registry.Result, err error) {
	for resp := range w.watchChan {
		err = resp.Err()
		if err != nil {
			return
		}
		for _, event := range resp.Events {
			service := decode(event.Kv.Value)
			var action string
			switch event.Type {
			case clientv3.EventTypePut:
				action = "create"
			case clientv3.EventTypeDelete:
				action = "delete"
				service = decode(event.PrevKv.Value)
			}

			if service == nil {
				continue
			}

			result = &registry.Result{
				Action:  action,
				Service: service,
			}
			return
		}
	}
	return nil, errors.New("not found next")
}

func (w *watcher) Stop() {
	log.Info(context.TODO(), "stop the etcd watcher")
}
