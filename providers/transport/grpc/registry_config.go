package grpc

import (
	registry "github.com/aka-yz/go-micro-core/register"
	"github.com/aka-yz/go-micro-core/utils/uuid"
	"go.uber.org/config"
)

type registryConfig struct {
	Addrs       []string
	RegistryTTL int
	Name        string
}

func getRegistryConfig(conf config.Provider) *registryConfig {
	var cv config.Value
	if cv = conf.Get("registry"); !cv.HasValue() {
		return nil
	}

	var cfg registryConfig
	if err := cv.Populate(&cfg); err != nil {
		return nil
	}

	if cfg.RegistryTTL == 0 {
		cfg.RegistryTTL = 30
	}

	return &cfg
}

func getRegistryService(conf config.Provider, suffix string) *registry.Service {
	var service registry.Service
	service.Name = conf.Get("name").String() + suffix
	service.Version = conf.Get("version").String()
	service.Nodes = make([]*registry.Node, 1)
	service.Nodes[0] = &registry.Node{
		Id: uuid.New(),
	}
	return &service
}
