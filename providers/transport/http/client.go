package http

import (
	"github.com/aka-yz/go-micro-core/providers/constants"
	"github.com/aka-yz/go-micro-core/providers/option"
	"go.uber.org/config"
	"time"
)

func newHTTPClient(cfg *option.HttpClientConfig) *HttpClient {
	var opt []Option
	if cfg.MaxConnectionNum != 0 {
		opt = append(opt, WithMaxConnectionNum(cfg.MaxConnectionNum))
	}

	if cfg.Timeout != 0 {
		opt = append(opt, WithTimeout(cfg.Timeout))
	}

	if cfg.Name != "" {
		opt = append(opt, WithServiceName(cfg.Name))
	}

	return NewHttpClient(opt...)
}

func getClientConfig(conf config.Provider) (cfg *option.HttpClientConfig) {
	cfg = &option.HttpClientConfig{}

	var cv config.Value
	if cv = conf.Get(constants.ConfigHttpClient); !cv.HasValue() {
		return
	}
	var hcc option.HttpClientConfig
	if err := cv.Populate(&cfg); err != nil {
		panic(err)
	}

	cfg.MaxConnectionNum = hcc.MaxConnectionNum
	cfg.Timeout = cfg.Timeout * time.Second
	return
}
