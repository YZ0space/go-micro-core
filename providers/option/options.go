package option

import (
	"crypto/tls"
	"time"
)

// 以下数据都是需要在yaml中配置的

type DB struct {
	Driver     string `default:"mysql"`
	DataSource string
	DBName     string `json:"dbname"`
	UserName   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	ReadHost   string `json:"readhost"`

	MaxIdleConns    int `json:"maxidleconns"`
	MaxOpenConns    int `json:"maxopenconns"`
	ConnMaxLifetime int `json:"connMaxLifetime"`

	ReadTimeout  int `json:"readtimeout"`
	WriteTimeout int `json:"writetimeout"`
	Timeout      int `json:"timeout"`
}

type Postgresql struct {
	// Network type, either tcp or unix.
	// Default is tcp.
	Network string
	// TCP host:port or Unix socket depending on Network.
	Addr     string
	User     string
	Password string
	Database string

	// ApplicationName is the application name. Used in logs on Pg side.
	// Only available from pg-9.0.
	ApplicationName string

	// TLS configs for secure connections.
	TLSConfig *tls.Config

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration

	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	WriteTimeout time.Duration

	// Maximum number of retries before giving up.
	// Default is to not retry failed queries.
	MaxRetries int
	// Whether to retry queries cancelled because of statement_timeout.
	RetryStatementTimeout bool
	// Minimum backoff between each retry.
	// Default is 250 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 4 seconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
	// Connection age at which client retires (closes) the connection.
	// It is useful with proxies like PgBouncer and HAProxy.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
	// Time for which client waits for free connection if all
	// connections are busy before returning an error.
	// Default is 30 seconds if ReadTimeOut is not defined, otherwise,
	// ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration
}

type Redis struct {
	ClusterAddr []string
	// host:port address.
	Addr string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string
	// Database to be selected after connecting to the server.
	DB int
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout int
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout int
	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout int64
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout int
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout int
	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int
	//
	IsClusterMode bool

	IsElastiCache bool
}

type HttpClientConfig struct {
	MaxConnectionNum int
	Timeout          time.Duration
	Name             string
}
