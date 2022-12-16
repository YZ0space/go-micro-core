package http

import (
	"context"
	"github.com/aka-yz/go-micro-core"
	"github.com/aka-yz/go-micro-core/configs/log"
	"github.com/aka-yz/go-micro-core/extension"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"github.com/gin-gonic/gin"
	"go.uber.org/config"
	"net/http"
	"os"
	"time"
)

type Server struct {
	r      *gin.Engine
	Server *http.Server
	opts   ServerOptions

	closeSyncJob  chan<- struct{}
	syncJobClosed <-chan struct{}
}

func (s *Server) Init() {
	handler := go_micro_core.ScanGinHandler(constants.HandlerInjectName)
	if handler == nil {
		log.Warnf(context.Background(), "handler 异常")

		return
	}
	s.r.Use(handler.MiddlewareList()...)
	s.AddHandlers(handler.HandlerList())
}

func (s *Server) Start() {
	go func() {
		go_micro_core.HttpErrCh <- s.Server.ListenAndServe()
	}()
	//go func() {
	//	服务注册
	//s.register()
	//}()
}

func (s *Server) Stop() {
	//if s.opts.registry == nil {
	//	return
	//}
	//if err := s.opts.registry.Deregister(s.opts.service); err != nil {
	//	log.Infof(context.Background(), "Deregister failed service:%v error:%v", json.MustString(s.opts.service), err)
	//}
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

//func (s *Server) register() {
//	if s.opts.registry == nil {
//		return
//	}
//	addr := strings.Split(s.opts.addr, ":")
//	if len(addr) != 2 {
//		panic(fmt.Errorf("error register addr:%v", addr))
//	}
//	port, _ := strconv.Atoi(addr[1])
//	s.opts.service.Nodes[0].Address = addr[0]
//	s.opts.service.Nodes[0].Port = port
//
//	for {
//		if err := s.opts.registry.Register(s.opts.service); err == nil {
//			log.Infof(context.Background(), "HTTP Server register:", json.MustString(s.opts.service))
//		}
//
//		time.Sleep(time.Second * 15)
//	}
//}

func (s *Server) AddHandlers(HandlerList []*extension.GinHandlerRegister) {
	for _, l := range HandlerList {
		s.r.Handle(l.HttpMethod, l.RelativePath, l.Handlers...)
	}
}

func newHTTPServer(cfg *serverConfig) *Server {
	r := gin.New()
	r.Use(Logger())
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

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}
		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path
		log.Infof(c, "[GIN] %3d| %13v | %15s |%-7s %#v |%s",
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage)

	}
}

type serverConfig struct {
	Addr string
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
	port := os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}

	return addr
}
