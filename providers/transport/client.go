package transport

import (
	"github.com/aka-yz/go-micro-core"
	"github.com/aka-yz/go-micro-core/providers/option"
	http2 "github.com/aka-yz/go-micro-core/providers/transport/http"
	"go.uber.org/config"
	"strconv"
	"time"
)

type clientFactory struct{}

func (n *clientFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	return go_micro_core.NewProvider(newHTTPClient(getClientConfig(conf)))
}

func newHTTPClient(cfg *option.HttpClientConfig) *http2.HttpClient {
	var opt []http2.Option
	if cfg.MaxConnectionNum != 0 {
		opt = append(opt, http2.WithMaxConnectionNum(cfg.MaxConnectionNum))
	}

	if cfg.Timeout != 0 {
		opt = append(opt, http2.WithTimeout(cfg.Timeout))
	}

	if cfg.Name != "" {
		opt = append(opt, http2.WithServiceName(cfg.Name))
	}

	return http2.NewHttpClient(opt...)
}

func getClientConfig(conf config.Provider) (cfg *option.HttpClientConfig) {
	cfg = &option.HttpClientConfig{}
	var cv config.Value

	cfg.Name = conf.Get("name").String()

	if cv = conf.Get("http_client"); !cv.HasValue() {
		return
	}

	maxConnectionNum, err := strconv.Atoi(cv.Get("maxconnectionnum").String())
	if err != nil {
		panic(err)
	}

	timeout, err := strconv.Atoi(cv.Get("timeout").String())
	if err != nil {
		panic(err)
	}

	cfg.MaxConnectionNum = maxConnectionNum
	cfg.Timeout = time.Duration(timeout) * time.Second
	return
}
