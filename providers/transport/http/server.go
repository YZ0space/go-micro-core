package http

import (
	"context"
	"github.com/aka-yz/go-micro-core"
	"github.com/aka-yz/go-micro-core/configs/log"
	"github.com/aka-yz/go-micro-core/extension"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"go.uber.org/config"
	"net/http"
	"os"
	"time"
)

type serverFactory struct{}

func (s *serverFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	if cfg := getServerConfig(conf); cfg != nil {
		srv := newHTTPServer(cfg)
		return go_micro_core.ProvideFunc(func() []*inject.Object {
			name := constants.ConfigSrvKey
			return []*inject.Object{
				&inject.Object{Name: name, Value: srv},
			}
		})
	}
	return nil
}

type Server struct {
	r      *gin.Engine
	Server *http.Server

	closeSyncJob  chan<- struct{}
	syncJobClosed <-chan struct{}
}

func (s *Server) Init() {
	handler := go_micro_core.ScanGinHandler(constants.HandlerInjectName)
	if handler == nil {
		log.Errorf(context.Background(), "handler 异常")
		return
	}
	s.r.Use(handler.MiddlewareList()...)
	s.addHandlers(handler.HandlerList())
}

func (s *Server) addHandlers(HandlerList []*extension.GinHandlerRegister) {
	for _, l := range HandlerList {
		s.r.Handle(l.HttpMethod, l.RelativePath, l.Handlers...)
	}
}

func (s *Server) Start() {
	go func() {
		go_micro_core.HttpErrCh <- s.Server.ListenAndServe()
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Errorf(context.TODO(), "Failed to gracefully shutdown server: %s", err)
	}
	select {
	case <-ctx.Done():
		log.Info(context.TODO(), "gracefully shutdown gin server and attached go routines: timeout")
	}
}

func newHTTPServer(cfg *serverConfig) *Server {
	r := gin.New()
	return &Server{
		r: r,
		Server: &http.Server{
			Addr:              cfg.Addr,
			Handler:           r,
			ReadHeaderTimeout: 10 * time.Second, // we should be safe behind istio
			ReadTimeout:       20 * time.Second, // setting them for go sec lint.
			WriteTimeout:      20 * time.Second,
		},
	}
}

type serverConfig struct {
	Addr  string
	PProf string
}

func getServerConfig(conf config.Provider) *serverConfig {
	var cv config.Value

	if cv = conf.Get("httpserver"); !cv.HasValue() {
		return nil
	}

	addrMap := make(map[string]string)
	if err := cv.Populate(&addrMap); err != nil {
		return nil
	}

	var cfg serverConfig
	cfg.Addr = port(addrMap["addr"])
	return &cfg
}

func port(addr string) string {
	port := os.Getenv("PORT_" + addr)
	if port != "" {
		return ":" + port
	}

	port = os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}

	return addr
}
