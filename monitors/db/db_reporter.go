package db

import (
	"time"
)

// 监控 耗时
type DBReporter struct {
	dbName        string
	costThreshold int64
}

func NewDBReporter(dbName string, costThreshold int64) *DBReporter {
	EnableHandlingTimeHistogram()
	return &DBReporter{
		dbName:        dbName,
		costThreshold: costThreshold,
	}
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data
func (s *DBReporter) EventErr(table string, operator string, err error) {
	if err != nil {
		dbClientResponse.WithLabelValues(s.dbName, table, operator).Inc()
	}
	return
}

// TimingKv receives the time an event took to happen along with optional key/value data
func (s *DBReporter) Timing(table string, operator string, nanoseconds int64) {
	dbClientRequestSend.WithLabelValues(s.dbName, table, operator).Inc()
	sec := time.Duration(nanoseconds) / time.Second
	nsec := time.Duration(nanoseconds) % time.Second
	serverHandledHistogram.WithLabelValues(s.dbName, table, operator).Observe(float64(sec) + float64(nsec)/1e9)
}
