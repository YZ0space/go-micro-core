package db

import (
	"fmt"
	prom "github.com/prometheus/client_golang/prometheus"
)

// db 监控指标
// 慢查询 暂定200ms
// qps
// 错误数和耗时

// 有的信息，db name，sql 语句（可以进一步解析 表名）， 错误信息，耗时
var (
	DefBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 40, 60}

	dbClientRequestSend = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "db",
			Subsystem: "client",
			Name:      "request_send_total",
			Help:      "Total number of request send to the server.",
		}, []string{"db", "table", "operator"})

	dbClientResponse = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "db",
			Subsystem: "client",
			Name:      "response_err_total",
			Help:      "Total number of request received on the server.",
		}, []string{"db", "table", "operator"})

	serverHandledHistogramEnabled = false
	serverHandledHistogramOpts    = prom.HistogramOpts{
		Namespace: "db",
		Subsystem: "client",
		Name:      "handling_seconds",
		Help:      "Histogram of response latency (seconds) of http handled by the server.",
		Buckets:   DefBuckets,
	}
	serverHandledHistogram *prom.HistogramVec
)

func init() {
	fmt.Printf("prom db init...\n")
	prom.MustRegister(dbClientRequestSend)
	prom.MustRegister(dbClientResponse)
}

type HistogramOption func(*prom.HistogramOpts)

// WithHistogramBuckets allows you to specify custom bucket ranges for histograms if EnableHandlingTimeHistogram is on.
func WithHistogramBuckets(buckets []float64) HistogramOption {
	return func(o *prom.HistogramOpts) { o.Buckets = buckets }
}

// EnableHandlingTimeHistogram turns on recording of handling time of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableHandlingTimeHistogram(opts ...HistogramOption) {
	for _, o := range opts {
		o(&serverHandledHistogramOpts)
	}
	if !serverHandledHistogramEnabled {
		serverHandledHistogram = prom.NewHistogramVec(
			serverHandledHistogramOpts,
			[]string{"db", "table", "operator"},
		)
		prom.MustRegister(serverHandledHistogram)
	}
	serverHandledHistogramEnabled = true
}
