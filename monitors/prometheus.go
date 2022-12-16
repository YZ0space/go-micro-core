package monitors

import (
	"github.com/aka-yz/go-micro-core/extension"
	"github.com/aka-yz/go-micro-core/providers/constants"
	hp "github.com/aka-yz/go-micro-core/providers/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func (m *PrometheusMetrics) Init() {
	m.HTTPServer.AddHandlers([]*extension.GinHandlerRegister{
		{
			HttpMethod:   constants.GinMethodGet,
			RelativePath: "/metrics",
			Handlers: []gin.HandlerFunc{
				func(c *gin.Context) {
					promhttp.Handler().ServeHTTP(c.Writer, c.Request)
				},
			},
		},
	})
}
