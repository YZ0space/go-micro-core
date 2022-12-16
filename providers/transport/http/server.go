package http

import (
	"github.com/aka-yz/go-micro-core"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"github.com/facebookgo/inject"
	"go.uber.org/config"
)

type serverFactory struct{}

func (s *serverFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	if cfg := getServerConfig(conf); cfg != nil {
		srv := newHTTPServer(cfg)
		return go_micro_core.ProvideFunc(func() []*inject.Object {
			name := constants.ConfigSrvKey
			return []*inject.Object{
				{Name: name, Value: srv},
				//{Value: monitor.NewPrometheusMetrics(srv)},
			}
		})
	}
	return nil
}
