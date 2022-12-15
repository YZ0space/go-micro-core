package cron

import (
	go_micro_core "github.com/aka-yz/go-micro-core"
	"github.com/aka-yz/go-micro-core/configs/middleware"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"github.com/facebookgo/inject"
	"go.uber.org/config"
)

func init() {
	go_micro_core.RegisterProvider(&cronFactory{})
}

type cronFactory struct{}

func (n *cronFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	// TODO: 读取 yaml 配置自定义 cron
	// ....
	return go_micro_core.ProvideFunc(func() []*inject.Object {
		cronPool := middleware.NewCron()
		return []*inject.Object{
			{Name: constants.ConfigCron, Value: cronPool},
		}
	})
}
