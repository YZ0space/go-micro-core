package http

import (
	"context"
	"github.com/aka-yz/go-micro-core"
	"github.com/aka-yz/go-micro-core/configs/log"
	"github.com/aka-yz/go-micro-core/providers/constants"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/config"
	"net/http"
	"os"
	"time"
)

type serverFactory struct{}

func (s *serverFactory) NewProvider(conf config.Provider) go_micro_core.Provider {
	if cfg := getServerConfig(conf); cfg != nil {
		return go_micro_core.NewProvider(newHTTPServer(cfg))
	}
	return nil
}

type Server struct {
	r      *gin.Engine
	Server *http.Server

	closeSyncJob  chan<- struct{}
	syncJobClosed <-chan struct{}
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
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// cors
	corsConfig := cors.New(cors.Config{
		AllowOrigins:     constants.AllowedOrigins,
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     constants.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "...."
		// },
		MaxAge: 12 * time.Hour,
	})
	r.Use(corsConfig)

	HTTPserver := &http.Server{
		Addr:              cfg.Addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second, // we should be safe behind istio
		ReadTimeout:       20 * time.Second, // setting them for go sec lint.
		WriteTimeout:      20 * time.Second,
	}
	//server.addHandlers()
	return &Server{
		r:      r,
		Server: HTTPserver,
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
