package redis

import (
	"github.com/aka-yz/go-micro-core"
	cf "github.com/aka-yz/go-micro-core/configs/middleware"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"github.com/aka-yz/go-micro-core/providers/option"
	"github.com/facebookgo/inject"
	"go.uber.org/config"
)

func init() {
	go_micro_core.RegisterProvider(&redisFactory{})
}

type redisFactory struct{}

func (n *redisFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	var opts map[string]*option.Redis

	var cv config.Value
	if cv = conf.Get(constants.ConfigKeyRedis); !cv.HasValue() {
		return nil
	}
	if err := cv.Populate(&opts); err != nil {
		panic(err)
	}

	return go_micro_core.ProvideFunc(func() []*inject.Object {
		var objects []*inject.Object
		for k, v := range opts {
			client := cf.NewClient(v)
			name := constants.ConfigKeyRedis + "." + k

			objects = append(objects, &inject.Object{Name: name, Value: client})
		}
		return objects
	})
}
