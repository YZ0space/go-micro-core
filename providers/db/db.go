package db

import (
	"github.com/aka-yz/go-micro-core"
	cf "github.com/aka-yz/go-micro-core/configs/middleware/db"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"github.com/aka-yz/go-micro-core/providers/option"
	"github.com/facebookgo/inject"
	"go.uber.org/config"
)

func init() {
	go_micro_core.RegisterProvider(&mysqlFactory{})
}

type mysqlFactory struct{}

func (n *mysqlFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	// 读取 yaml 配置并初始化 connection
	var opts map[string]*option.DB
	var cv config.Value
	if cv = conf.Get(constants.ConfigKeyMysql); !cv.HasValue() {
		return nil
	}
	if err := cv.Populate(&opts); err != nil {
		panic(err)
	}

	return go_micro_core.ProvideFunc(func() []*inject.Object {
		var objects []*inject.Object
		for k, v := range opts {
			conn := cf.OpenDB(v)
			name := constants.ConfigKeyMysql + "." + k

			objects = append(objects, &inject.Object{Name: name, Value: conn})
		}
		return objects
	})
}
