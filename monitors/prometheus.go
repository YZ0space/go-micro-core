package monitors

import (
	"github.com/aka-yz/go-micro-core/extension"
	"github.com/aka-yz/go-micro-core/providers/constants"
	hp "github.com/aka-yz/go-micro-core/providers/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type PrometheusMetrics struct {
	HTTPServer *hp.Server
}

const (
	SUCCESS string = "success"
	Fail    string = "fail"
)

func NewPrometheusMetrics(hs *hp.Server) *PrometheusMetrics {
	return &PrometheusMetrics{
		HTTPServer: hs,
	}
}

func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func (m *PrometheusMetrics) Init() {
	m.HTTPServer.AddHandlers([]*extension.GinHandlerRegister{
		{
			HttpMethod:   constants.GinMethodGet,
			RelativePath: "/metrics",
			Handlers: []gin.HandlerFunc{
				PromHandler(promhttp.Handler()),
			},
		},
	})
}
