package providers

import (
	"github.com/facebookgo/inject"
	go_micro_core "go-micro-core"
	cf "go-micro-core/config"
	"go-micro-core/option"
	"go.uber.org/config"
)

func init() {
	go_micro_core.RegisterProvider(&pgFactory{})
}

type pgFactory struct{}

func (n *pgFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	// 读取 yaml 配置并初始化 connection
	var opts map[string]*option.Postgresql
	var cv config.Value
	if cv = conf.Get("pg"); !cv.HasValue() {
		return nil
	}
	if err := cv.Populate(&opts); err != nil {
		panic(err)
	}

	return go_micro_core.ProvideFunc(func() []*inject.Object {
		var objects []*inject.Object
		for k, v := range opts {
			conn := cf.OpenPG(v)
			name := "db." + k

			objects = append(objects, &inject.Object{Name: name, Value: conn})
		}
		return objects
	})
}
