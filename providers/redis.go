package providers

import (
	"github.com/facebookgo/inject"
	"github.com/go-redis/redis/v8"
	go_micro_core "go-micro-core"
	"go.uber.org/config"
)

func init() {
	go_micro_core.RegisterProvider(&redisFactory{})
}

type redisFactory struct{}

func (n *redisFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	var opts map[string]*redis.Options

	var cv config.Value
	if cv = conf.Get("redis"); !cv.HasValue() {
		return nil
	}
	if err := cv.Populate(&opts); err != nil {
		panic(err)
	}

	return go_micro_core.ProvideFunc(func() []*inject.Object {
		var objects []*inject.Object
		for k, v := range opts {
			client := redis.NewClient(v)
			name := "redis." + k

			objects = append(objects, &inject.Object{Name: name, Value: client})
		}
		return objects
	})
}
